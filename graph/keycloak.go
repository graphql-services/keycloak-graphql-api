package graph

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"

	"golang.org/x/oauth2/clientcredentials"
)

type KeycloakUser struct {
	ID        string
	Email     string
	FirstName string
	LastName  string
}

type KeycloakAPI struct {
	client *http.Client
}

func NewKeycloakAPI() *KeycloakAPI {
	return &KeycloakAPI{}
}

func (api *KeycloakAPI) getClient(ctx context.Context) *http.Client {
	c := clientcredentials.Config{
		ClientID:     os.Getenv("KEYCLOAK_CLIENT_ID"),
		ClientSecret: os.Getenv("KEYCLOAK_CLIENT_SECRET"),
		TokenURL:     os.Getenv("KEYCLOAK_TOKEN_URL"),
	}
	return c.Client(ctx)
}
func (api *KeycloakAPI) isActive() bool {
	return os.Getenv("KEYCLOAK_REST_API_URL") != ""
}

func (api *KeycloakAPI) getURL(p string, qs map[string]string) string {
	u, _ := url.Parse(os.Getenv("KEYCLOAK_REST_API_URL"))
	u.Path = path.Join(u.Path, p)
	q := u.Query()
	for key, val := range qs {
		q.Add(key, val)
	}
	u.RawQuery = q.Encode()
	return u.String()
}

func (api *KeycloakAPI) InviteMember(ctx context.Context, email, firstName, lastName string) (user KeycloakUser, err error) {
	if !api.isActive() {
		user = KeycloakUser{
			ID: email,
		}
		return
	}
	users, err := api.GetUsersByUsername(ctx, email)
	if err != nil {
		return
	}
	if len(users) == 0 {
		err = api.CreateUser(ctx, email, firstName, lastName)
		if err != nil {
			return
		}
		users, err = api.GetUsersByUsername(ctx, email)
		if err != nil {
			return
		}
		if len(users) == 0 {
			err = fmt.Errorf("unable to create user")
			return
		}
		err = api.ExecuteActionsEmail(ctx, users[0].ID, []string{"VERIFY_EMAIL", "UPDATE_PROFILE", "UPDATE_PASSWORD"})
		if err != nil {
			return
		}
	}

	user = users[0]

	return
}

func (api *KeycloakAPI) GetUser(ctx context.Context, id string) (user *KeycloakUser, err error) {
	c := api.getClient(ctx)
	res, err := c.Get(api.getURL("/users/"+id, map[string]string{}))
	if err != nil {
		return
	}
	fmt.Println("??", res.StatusCode)
	if res.StatusCode == 404 {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&user)
	return
}

func (api *KeycloakAPI) GetUsersByUsername(ctx context.Context, email string) (users []KeycloakUser, err error) {
	c := api.getClient(ctx)
	res, err := c.Get(api.getURL("/users", map[string]string{"username": email}))
	if err != nil {
		return
	}

	err = json.NewDecoder(res.Body).Decode(&users)
	return
}

func (api *KeycloakAPI) CreateUser(ctx context.Context, email, firstname, lastname string) (err error) {
	c := api.getClient(ctx)

	type User struct {
		Username  string `json:"username"`
		Email     string `json:"email"`
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
		Enabled   bool   `json:"enabled"`
	}

	u := User{Username: email, Email: email, FirstName: firstname, LastName: lastname, Enabled: true}
	reqBody, _ := json.Marshal(u)
	res, err := c.Post(api.getURL("/users", map[string]string{}), "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}

	if res.StatusCode != 201 {
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	return
}

func (api *KeycloakAPI) ExecuteActionsEmail(ctx context.Context, userId string, actions []string) (err error) {
	c := api.getClient(ctx)

	reqBody, _ := json.Marshal(actions)
	req, err := http.NewRequest("PUT", api.getURL("/users/"+userId+"/execute-actions-email", map[string]string{}), bytes.NewBuffer(reqBody))
	if err != nil {
		return
	}
	req.Header.Set("content-type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return
	}

	if res.StatusCode != 204 {
		err = fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

	return
}

package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/graphql-services/id/graph/generated"
	"github.com/graphql-services/id/graph/model"
)

func (r *entityResolver) FindUserByID(ctx context.Context, id string) (user *model.User, err error) {
	fmt.Println("start resolver", id)
	kc := NewKeycloakAPI()
	keycloakUser, err := kc.GetUser(ctx, id)
	if err != nil {
		return
	}
	if keycloakUser == nil {
		return
	}

	user = &model.User{
		ID:         keycloakUser.ID,
		Email:      keycloakUser.Email,
		FamilyName: &keycloakUser.FirstName,
		GivenName:  &keycloakUser.LastName,
	}
	fmt.Println("end resolver", id)
	return
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }

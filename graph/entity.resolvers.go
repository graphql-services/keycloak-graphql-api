package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/graphql-services/id/graph/generated"
	"github.com/graphql-services/id/graph/model"
)

func (r *entityResolver) FindUserByID(ctx context.Context, id string) (*model.User, error) {
	kc := NewKeycloakAPI()
	keycloakUser, err := kc.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	if keycloakUser == nil {
		return nil, nil
	}

	user := &model.User{
		ID:         keycloakUser.ID,
		Email:      keycloakUser.Email,
		FamilyName: &keycloakUser.FirstName,
		GivenName:  &keycloakUser.LastName,
	}
	return user, nil
}

// Entity returns generated.EntityResolver implementation.
func (r *Resolver) Entity() generated.EntityResolver { return &entityResolver{r} }

type entityResolver struct{ *Resolver }

package graphql

import (
	model1 "api-shapes/transport/graphql/graph/model"
	"api-shapes/store"
	"context"

	"github.com/google/uuid"
)

type Resolver interface {
	List(ctx context.Context, limit *int) (*model1.ListRes, error)
	Retrieve(ctx context.Context, userID uuid.UUID) (*model1.UserRes, error)
	Create(ctx context.Context, input model1.CreateUser) (*model1.UserRes, error)
	Update(ctx context.Context, input model1.UpdateUser) (*model1.UserRes, error)
	Delete(ctx context.Context, userID uuid.UUID) (*model1.GeneralMutationRes, error)
}

type resolver struct{}

func NewResolver() Resolver {
	return &resolver{}
}

func (r *resolver) List(ctx context.Context, limit *int) (*model1.ListRes, error) {
	l := store.List()

	var res model1.ListRes
	for _, e := range l {
		var ur model1.UserRes
		ur.ID = e.ID
		ur.Name = e.Name
		res.Users = append(res.Users, &ur)
	}

	return &res, nil
}

func (r *resolver) Retrieve(ctx context.Context, userID uuid.UUID) (*model1.UserRes, error) {
	e := store.Retrieve(userID.String())

	var res model1.UserRes
	res.ID = e.ID
	res.Name = e.Name
	return &res, nil
}

func (r *resolver) Create(ctx context.Context, input model1.CreateUser) (*model1.UserRes, error) {
	e := store.Create(store.User{ID: uuid.NewString(), Name: input.Name})

	var res model1.UserRes
	res.ID = e.ID
	res.Name = e.Name
	return &res, nil
}

func (r *resolver) Update(ctx context.Context, input model1.UpdateUser) (*model1.UserRes, error) {
	e := store.Retrieve(input.ID)
	e.Name = input.Name
	e = store.Update(e)

	var res model1.UserRes
	res.ID = e.ID
	res.Name = e.Name
	return &res, nil
}

func (r *resolver) Delete(ctx context.Context, userID uuid.UUID) (*model1.GeneralMutationRes, error) {
	store.Delete(userID.String())

	var res model1.GeneralMutationRes
	res.Success = true
	return &res, nil
}

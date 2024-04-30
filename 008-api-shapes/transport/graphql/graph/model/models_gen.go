// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type CreateUser struct {
	Name string `json:"name"`
}

type GeneralMutationRes struct {
	Success bool `json:"success"`
}

type ListRes struct {
	Users []*UserRes `json:"users,omitempty"`
}

type Mutation struct {
}

type Query struct {
}

type UpdateUser struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type UserRes struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

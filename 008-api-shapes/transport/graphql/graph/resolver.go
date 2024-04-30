package graph

import "api-shapes/transport/graphql"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	userResolver graphql.Resolver
}

func NewResolver(userResolver graphql.Resolver) ResolverRoot {
	return &Resolver{userResolver}
}

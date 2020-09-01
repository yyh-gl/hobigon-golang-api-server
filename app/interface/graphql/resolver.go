package graphql

import "github.com/yyh-gl/hobigon-golang-api-server/app/domain/repository"

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver : GraphQL標準Resolver
type Resolver struct {
	BlogRepository repository.Blog
}

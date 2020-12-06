package graphql

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"fmt"

	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/graphql/generated"
	"github.com/yyh-gl/hobigon-golang-api-server/app/presentation/graphql/model"
)

func (r *mutationResolver) CreateBlog(ctx context.Context, input model.NewBlog) (*model.Blog, error) {
	panic(fmt.Errorf("not implemented"))
}

func (r *queryResolver) Blog(ctx context.Context, title string) (*model.Blog, error) {
	blog, err := r.BlogRepository.FindByTitle(ctx, title)
	if err != nil {
		// TODO: NotFound時のエラーハンドリング
		return nil, fmt.Errorf("blogRepository.FindByTitle()内でのエラー: %w", err)
	}
	return &model.Blog{
		Title: blog.Title().String(),
		Count: blog.Count().Int(),
	}, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

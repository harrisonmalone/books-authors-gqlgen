package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"

	"github.com/harrisonmalone/authors-books-app/graph/database"
	"github.com/harrisonmalone/authors-books-app/graph/generated"
	"github.com/harrisonmalone/authors-books-app/graph/model"
	"github.com/harrisonmalone/authors-books-app/graph/utils"
)

func (r *mutationResolver) CreateBook(ctx context.Context, input model.BookInput) (*model.Book, error) {
	book := &database.Book{
		Title:    input.Title,
		AuthorId: input.AuthorID,
	}

	createdBook, err := book.Create()

	if err != nil {
		return nil, err
	}

	return &model.Book{
		ID:    createdBook.ID,
		Title: createdBook.Title,
		Author: &model.Author{
			ID:   createdBook.Author.ID,
			Name: createdBook.Author.Name,
		},
	}, nil
}

func (r *mutationResolver) CreateAuthor(ctx context.Context, input model.AuthorInput) (*model.Author, error) {
	author := &database.Author{
		Name: input.Name,
	}

	createdAuthor, err := author.Create()

	if err != nil {
		return nil, err
	}

	return &model.Author{
		ID:    createdAuthor.ID,
		Name:  createdAuthor.Name,
		Books: []*model.Book{},
	}, nil
}

func (r *queryResolver) Book(ctx context.Context, id string) (*model.Book, error) {
	book, err := database.Book{}.Find(id)

	if err != nil {
		return nil, err
	}

	return &model.Book{
		ID:    book.ID,
		Title: book.Title,
		Author: &model.Author{
			ID:   book.Author.ID,
			Name: book.Author.Name,
		},
	}, nil
}

func (r *queryResolver) Author(ctx context.Context, id string) (*model.Author, error) {
	author, err := database.Author{}.Find(id)

	if err != nil {
		return nil, err
	}

	return &model.Author{
		ID:    author.ID,
		Name:  author.Name,
		Books: utils.MapGraphBooks(author.Books),
	}, nil
}

func (r *queryResolver) Authors(ctx context.Context) ([]*model.Author, error) {
	authors, err := database.Author{}.All()

	if err != nil {
		return nil, err
	}

	return utils.MapGraphAuthors(authors), nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }

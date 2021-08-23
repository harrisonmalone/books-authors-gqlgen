package utils

import (
	"github.com/harrisonmalone/authors-books-app/graph/database"
	"github.com/harrisonmalone/authors-books-app/graph/model"
)

func MapGraphBooks(books []database.Book) []*model.Book {
	var graphBooks []*model.Book
	for _, book := range books {
		graphBooks = append(graphBooks, &model.Book{
			ID:    book.ID,
			Title: book.Title,
		})
	}
	return graphBooks
}


func MapGraphAuthors(authors []database.Author) []*model.Author {
	var graphAuthors []*model.Author
	for _, author := range authors {
		graphAuthors = append(graphAuthors, &model.Author{
			ID: author.ID,
			Name: author.Name,
			Books: MapGraphBooks(author.Books),
		})
	}
	return graphAuthors
}
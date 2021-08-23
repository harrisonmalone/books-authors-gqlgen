package database

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

var connection *sql.DB

type Author struct {
	ID    string
	Name  string
	Books []Book
}

type Book struct {
	ID       string
	Title    string
	AuthorId string
	Author   Author
}

func (a Author) Find(id string) (Author, error) {
	var author Author

	authorQuery := `
	SELECT
		authors.id,
		authors. "name"
	FROM
		authors
	WHERE
		authors.id = ?;
	`

	stmt, err := connection.Prepare(authorQuery)

	if err != nil {
		return author, err
	}

	rows, _ := stmt.Query(id)

	var authorId, name string

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&authorId, &name)

		if err != nil {
			return author, err
		}
	}

	author = Author{
		ID:   authorId,
		Name: name,
	}

	authorBooksQuery := `
	SELECT
		books.id,
		books.title
	FROM
		books
	WHERE
		books.author_id = ?;
	`

	stmt, err = connection.Prepare(authorBooksQuery)

	if err != nil {
		return author, err
	}

	rows, _ = stmt.Query(id)

	var bookId, title string

	for rows.Next() {
		err := rows.Scan(&bookId, &title)

		if err != nil {
			return author, err
		}

		author.Books = append(author.Books, Book{
			ID:    bookId,
			Title: title,
		})
	}

	rows.Close()

	return author, nil
}

func (a *Author) Create() (Author, error) {
	var author Author

	stmt, err := connection.Prepare("INSERT INTO authors(name) values(?)")

	if err != nil {
		return author, err
	}

	res, err := stmt.Exec(a.Name)

	if err != nil {
		return author, err
	}

	id, err := res.LastInsertId()

	strInt := strconv.Itoa(int(id))

	author, err = a.Find(strInt)

	return author, err
}

func (a Author) All() ([]Author, error) {
	var authors []Author

	authorQuery := `
	SELECT
		authors.id,
		authors. "name"
	FROM
		authors;
	`

	stmt, err := connection.Prepare(authorQuery)

	if err != nil {
		return authors, err
	}

	authorRows, _ := stmt.Query()

	var authorId, name string

	for authorRows.Next() {
		err := authorRows.Scan(&authorId, &name)

		if err != nil {
			return authors, err
		}

		booksQuery := `
		SELECT
			books.id,
			books.title
		FROM
			books
		WHERE
			author_id = ?;	
		`

		stmt, err := connection.Prepare(booksQuery)

		if err != nil {
			return authors, err
		}

		bookRows, _ := stmt.Query(authorId)

		var (
			bookId, title string
			books         []Book
		)

		for bookRows.Next() {
			err := bookRows.Scan(&bookId, &title)

			if err != nil {
				return authors, err
			}

			books = append(books, Book{
				ID:    bookId,
				Title: title,
			})
		}

		bookRows.Close()

		authors = append(authors, Author{
			ID:    authorId,
			Name:  name,
			Books: books,
		})
	}

	authorRows.Close()

	return authors, nil
}

func (b *Book) Create() (Book, error) {
	var book Book

	stmt, err := connection.Prepare("INSERT INTO books(title, author_id) values(?,?)")

	if err != nil {
		return book, err
	}

	res, err := stmt.Exec(b.Title, b.AuthorId)

	if err != nil {
		return book, err
	}

	id, err := res.LastInsertId()

	strInt := strconv.Itoa(int(id))

	book, err = b.Find(strInt)

	return book, err
}

func (b Book) Find(id string) (Book, error) {
	var book Book

	query := `
	SELECT
		books.id,
		books.title, 
		authors.id,
		authors."name"
	FROM
		books
		INNER JOIN authors ON books.author_id
	WHERE 
		books.author_id = authors.id
		AND books.id = ?;
	`

	stmt, err := connection.Prepare(query)

	if err != nil {
		return book, err
	}

	rows, _ := stmt.Query(id)

	var bookID, title, authorId, authorName string

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&bookID, &title, &authorId, &authorName)

		if err != nil {
			return book, err
		}
	}

	return Book{
		ID:    bookID,
		Title: title,
		Author: Author{
			ID:   authorId,
			Name: authorName,
		},
	}, err
}

func Initialize() {
	var err error
	db, err := sql.Open("sqlite3", "./books.db")

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}
	// ensure all migrations have been run
	connection = db
}

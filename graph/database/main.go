package database

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/lib/pq"
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
		authors.id = $1;
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
		books.author_id = $1;
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

	lastInsertId := 0
	err := connection.QueryRow(`INSERT INTO "authors" (name) values($1) RETURNING id`, a.Name).Scan(&lastInsertId)

	if err != nil {
		return author, err
	}

	strInt := strconv.Itoa(int(lastInsertId))

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
			author_id = $1;	
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

	lastInsertId := 0
	err := connection.QueryRow(`INSERT INTO "books" (title, author_id) values($1, $2) RETURNING id`, b.Title, b.AuthorId).Scan(&lastInsertId)

	if err != nil {
		return book, err
	}

	strInt := strconv.Itoa(int(lastInsertId))

	book, err = b.Find(strInt)

	return book, err
}

func (b Book) Find(id string) (Book, error) {
	var book Book

	log.Println(id)
	query := `
	SELECT
		books.id,
		books.title,
		authors.id,
		authors. "name"
	FROM
		books
		INNER JOIN authors ON books.author_id = authors.id
	WHERE
		books.id = $1;
	`

	stmt, err := connection.Prepare(query)

	if err != nil {
		log.Println(err)
		return book, err
	}

	rows, _ := stmt.Query(id)

	var bookID, title, authorId, authorName string

	defer rows.Close()

	if rows.Next() {
		err := rows.Scan(&bookID, &title, &authorId, &authorName)
		log.Println(bookID)
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
	config := "user=harrisonmalone password= host=127.0.0.1 port=5432 dbname=books connect_timeout=20 sslmode=disable"
	db, err := sql.Open("postgres", config)

	if err != nil {
		fmt.Println(err)
		panic("failed to connect database")
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	// ensure all migrations have been run
	connection = db
}

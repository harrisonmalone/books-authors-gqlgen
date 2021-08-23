CREATE TABLE authors (
	id SERIAL PRIMARY KEY,
	name VARCHAR
);

CREATE TABLE books (
	id SERIAL PRIMARY KEY,
	title text,
	author_id bigint,
	CONSTRAINT fk_authors_books FOREIGN KEY (author_id) REFERENCES authors (id)
);
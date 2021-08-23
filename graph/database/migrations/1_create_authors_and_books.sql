CREATE TABLE `authors` (
	`id` integer,
	`name` text,
	PRIMARY KEY (`id`)
);

CREATE TABLE `books` (
	`id` integer,
	`title` text,
	`author_id` integer,
	PRIMARY KEY (`id`),
	CONSTRAINT `fk_authors_books` FOREIGN KEY (`author_id`) REFERENCES `authors` (`id`)
);
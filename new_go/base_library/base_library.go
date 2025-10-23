package base_library

import (
	"book-lib/book"
)

type Library interface {
	AddBook(book book.Book)
	ChangeIdentificator(idf func(book book.Book) int)
	ChangeBooks(books []book.Book)
	GetBook(name string) book.Book
}

type BaseLibrary struct {
	BookNameToId  map[string]int
	Identificator func(book book.Book) int
}

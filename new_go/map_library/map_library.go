package map_library

import (
	"book-lib/base_library"
	"book-lib/book"
	"fmt"
)

type MapLibrary struct {
	BaseLibrary base_library.BaseLibrary
	Database    map[int]book.Book
}

func NewMapLibrary() *MapLibrary {
	lib := &MapLibrary{}
	lib.Database = make(map[int]book.Book)
	lib.BaseLibrary.Identificator = func(book book.Book) int {
		return len(lib.Database)
	}
	lib.BaseLibrary.BookNameToId = make(map[string]int)
	return lib
}

func (lib *MapLibrary) AddBook(book book.Book) {
	var book_id = lib.BaseLibrary.Identificator(book)
	lib.BaseLibrary.BookNameToId[book.Name] = book_id
	lib.Database[book_id] = book
}

func (lib *MapLibrary) ChangeIdentificator(idf func(book book.Book) int) {
	lib.BaseLibrary.Identificator = idf
	OldToNewId := make(map[int]int)
	for id, book := range lib.Database {
		OldToNewId[id] = idf(book)
	}
	NewBookNameToId := make(map[string]int)
	NewDatabase := make(map[int]book.Book)
	for name, id := range lib.BaseLibrary.BookNameToId {
		NewBookNameToId[name] = OldToNewId[id]
		NewDatabase[OldToNewId[id]] = lib.Database[id]
	}
	lib.BaseLibrary.BookNameToId = NewBookNameToId
	lib.Database = NewDatabase
}

func (lib *MapLibrary) ChangeBooks(books []book.Book) {
	for k := range lib.BaseLibrary.BookNameToId {
		delete(lib.BaseLibrary.BookNameToId, k)
	}
	for k := range lib.Database {
		delete(lib.Database, k)
	}
	for _, book := range books {
		lib.AddBook(book)
	}
}

func (lib MapLibrary) GetBook(name string) book.Book {
	if id, exists := lib.BaseLibrary.BookNameToId[name]; exists {
		return lib.Database[id]
	}
	fmt.Println("Такой книги нет!")
	return book.Book{}
}

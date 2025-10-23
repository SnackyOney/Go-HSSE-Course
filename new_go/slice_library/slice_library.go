package slice_library

import (
	"book-lib/base_library"
	"book-lib/book"
	"fmt"
)

type IdAndBook struct {
	Id   int
	Book book.Book
}

type SliceLibrary struct {
	BaseLibrary base_library.BaseLibrary
	Database    []IdAndBook
}

func NewSliceLibrary() *SliceLibrary {
	lib := &SliceLibrary{}
	lib.Database = make([]IdAndBook, 0)
	lib.BaseLibrary.Identificator = func(book book.Book) int {
		return len(lib.Database)
	}
	lib.BaseLibrary.BookNameToId = make(map[string]int)
	return lib
}

func (lib *SliceLibrary) AddBook(book book.Book) {
	book_id := lib.BaseLibrary.Identificator(book)
	lib.BaseLibrary.BookNameToId[book.Name] = book_id
	lib.Database = append(lib.Database, IdAndBook{book_id, book})
}

func (lib *SliceLibrary) ChangeIdentificator(idf func(book book.Book) int) {
	lib.BaseLibrary.Identificator = idf
	OldToNewId := make(map[int]int)
	for id, id_and_book := range lib.Database {
		OldToNewId[id] = idf(id_and_book.Book)
	}
	NewBookNameToId := make(map[string]int)
	NewDatabase := make([]IdAndBook, 0)
	for name, id := range lib.BaseLibrary.BookNameToId {
		NewBookNameToId[name] = OldToNewId[id]
		NewDatabase = append(NewDatabase, IdAndBook{OldToNewId[id], lib.Database[id].Book})
	}
	lib.BaseLibrary.BookNameToId = NewBookNameToId
	lib.Database = NewDatabase
}

func (lib *SliceLibrary) ChangeBooks(books []book.Book) {
	for k := range lib.BaseLibrary.BookNameToId {
		delete(lib.BaseLibrary.BookNameToId, k)
	}
	lib.Database = nil
	for _, book := range books {
		lib.AddBook(book)
	}
}

func (lib SliceLibrary) GetBook(name string) book.Book {
	book_id := lib.BaseLibrary.BookNameToId[name]
	for _, id_and_book := range lib.Database {
		if id_and_book.Id == book_id {
			return id_and_book.Book
		}
	}
	fmt.Println("Такой книги нет!")
	return book.Book{}
}

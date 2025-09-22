package main

import "fmt"

const kPrime = 19

type Book struct {
	Name    string
	Filling string
}

type Library struct {
	BookNameToId  map[string]int
	Database      map[int]Book
	Identificator func(book Book) int
}

func NewLibrary() *Library {
	lib := &Library{
		BookNameToId: make(map[string]int),
		Database:     make(map[int]Book),
	}
	lib.Identificator = func(book Book) int {
		return len(lib.Database)
	}
	return lib
}

func (lib *Library) AddBook(book Book) {
	var book_id = lib.Identificator(book)
	lib.BookNameToId[book.Name] = book_id
	lib.Database[book_id] = book
}

func (lib *Library) ChangeIdentificator(idf func(book Book) int) {
	lib.Identificator = idf
	OldToNewId := make(map[int]int)
	for id, book := range lib.Database {
		OldToNewId[id] = idf(book)
	}
	NewBookNameToId := make(map[string]int)
	NewDatabase := make(map[int]Book)
	for name, id := range lib.BookNameToId {
		NewBookNameToId[name] = OldToNewId[id]
		NewDatabase[OldToNewId[id]] = lib.Database[id]
	}
	lib.BookNameToId = NewBookNameToId
	lib.Database = NewDatabase
}

func (lib *Library) ChangeBooks(books []Book) {
	for k := range lib.BookNameToId {
		delete(lib.BookNameToId, k)
	}
	for k := range lib.Database {
		delete(lib.Database, k)
	}
	for _, book := range books {
		lib.AddBook(book)
	}
}

func (lib Library) GetBook(name string) Book {
	if id, exists := lib.BookNameToId[name]; exists {
		return lib.Database[id]
	}
	fmt.Println("Такой книги нет!")
	return Book{}
}

func main() {
	book_lib := NewLibrary()
	books := []Book{{Name: "Tom and Jerry", Filling: "Jerry has won..."}, {Name: "1984", Filling: "I remember these days"}, {Name: "Africa", Filling: "We survive"}}
	for _, book := range books {
		book_lib.AddBook(book)
	}
	fmt.Println(book_lib)
	tmj_book := book_lib.GetBook("Tom and Jerry")
	fmt.Println(tmj_book)
	some_book := book_lib.GetBook("Tom and Lololoshka")
	fmt.Println(some_book)
	book_lib.ChangeIdentificator(func(book Book) int {
		hash := 0
		for _, letter := range book.Name {
			hash += int(letter) * kPrime
		}
		return hash
	})
	new_tmj_book := book_lib.GetBook("Tom and Jerry")
	fmt.Println(new_tmj_book)
	fmt.Println(book_lib)
	other_books := []Book{{Name: "Pressly", Filling: "Will he press it..."}, {Name: "L", Filling: "TAKE IT"}}
	book_lib.ChangeBooks(other_books)
	fmt.Println(book_lib)
	l_book := book_lib.GetBook("L")
	fmt.Println(l_book)
	g_book := book_lib.GetBook("G")
	fmt.Println(g_book)
}

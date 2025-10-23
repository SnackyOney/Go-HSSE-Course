package main

import (
	"book-lib/base_library"
	"book-lib/book"
	"book-lib/map_library"
	"book-lib/slice_library"
	"fmt"
)

const kPrime = 19 // used for proper hash compution

func MapLibTest() {
	book_lib := map_library.NewMapLibrary()
	LibTests(book_lib)
}

func SliceLibTest() {
	book_lib := slice_library.NewSliceLibrary()
	LibTests(book_lib)
}

func LibTests(book_lib base_library.Library) {
	books := []book.Book{{Name: "Tom and Jerry", Filling: "Jerry has won..."}, {Name: "1984", Filling: "I remember these days"}, {Name: "Africa", Filling: "We survive"}}
	for _, book := range books {
		book_lib.AddBook(book)
	}
	fmt.Println(book_lib)
	tmj_book := book_lib.GetBook("Tom and Jerry")
	fmt.Println(tmj_book)
	some_book := book_lib.GetBook("Tom and Lololoshka")
	fmt.Println(some_book)
	book_lib.ChangeIdentificator(func(book book.Book) int {
		hash := 0
		for _, letter := range book.Name {
			hash += int(letter) * kPrime
		}
		return hash
	})
	new_tmj_book := book_lib.GetBook("Tom and Jerry")
	fmt.Println(new_tmj_book)
	fmt.Println(book_lib)
	other_books := []book.Book{{Name: "Pressly", Filling: "Will he press it..."}, {Name: "L", Filling: "TAKE IT"}}
	book_lib.ChangeBooks(other_books)
	fmt.Println(book_lib)
	l_book := book_lib.GetBook("L")
	fmt.Println(l_book)
	g_book := book_lib.GetBook("G")
	fmt.Println(g_book)
}

func main() {
	MapLibTest()
	SliceLibTest()
}

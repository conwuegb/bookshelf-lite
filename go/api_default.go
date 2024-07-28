/*
 * Bookshelf API
 *
 * An API which manages the inventory of books on a bookshelf.
 *
 */
package swagger

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Set up "DB"
var bookshelf []Book = nil

func getDB() *[]Book {
	if bookshelf == nil {
		bookshelf = []Book{}
	}
	return &bookshelf
}

// Alias db access func for testing purposes
var GetDB = getDB

// Handler functions
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, and welcome to your bookshelf!")
}

func BooksGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Fetch books from DB and encode to JSON in response
	var books *[]Book = GetDB()
	json.NewEncoder(w).Encode(books)
}

func BooksPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Decode book from request params
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	// Add book to DB
	createBook(&book, w)

	// Add new (or existing) book to response
	json.NewEncoder(w).Encode(book)

}

func createBook(book *Book, w http.ResponseWriter) {
	// Get DB
	var db *[]Book = GetDB()

	// Check if book isbn already exists in DB
	for _, shelfBook := range *db {
		if shelfBook.Isbn == book.Isbn {
			w.WriteHeader(http.StatusOK)
			fmt.Fprintf(w, "This book is already on your shelf!\n")
			return
		}
	}

	// Add book
	*db = append(*db, *book)
	w.WriteHeader(http.StatusCreated)
}

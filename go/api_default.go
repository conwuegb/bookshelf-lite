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
	"strings"

	"github.com/gorilla/mux"
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
	fmt.Fprintln(w, "Hello, and welcome to your bookshelf!")
}

func BooksGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Fetch books from DB and encode to JSON in response
	json.NewEncoder(w).Encode(getAllBooks())
}

func BooksPost(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	// Decode book from request body
	var book Book
	_ = json.NewDecoder(r.Body).Decode(&book)

	// Add book to DB
	respCode, msg := createBook(&book)

	// Add new (or existing) book to response
	w.WriteHeader(respCode)
	if msg != "" {
		fmt.Fprintln(w, msg)
	}
	json.NewEncoder(w).Encode(book)

}

func BooksDeleteByIsbn(w http.ResponseWriter, r *http.Request) {
	// Decode book from request params
	params := mux.Vars(r)

	// Find and delete book
	respCode, msg := deleteBook(strings.TrimSpace(params["isbn"]))

	// Add new (or existing) book to response
	w.WriteHeader(respCode)
	if msg != "" {
		fmt.Fprintln(w, msg)
	}
}

// Database functions
func getAllBooks() *[]Book {
	return GetDB()
}

func createBook(book *Book) (int, string) {
	var db *[]Book = GetDB()

	// Check if book isbn already exists in DB
	for _, shelfBook := range *db {
		if shelfBook.Isbn == book.Isbn {
			return http.StatusOK, "This book is already on your shelf!"
		}
	}

	// Add book
	*db = append(*db, *book)
	return http.StatusCreated, ""
}

func deleteBook(isbn string) (int, string) {
	var db *[]Book = GetDB()

	// Find isbn in DB and delete
	for idx, shelfBook := range *db {
		if shelfBook.Isbn == isbn {
			temp := append((*db)[:idx], (*db)[idx+1:]...)
			*db = temp
			return http.StatusNoContent, ""
		}
	}

	return http.StatusNotFound, "This book was not found on your shelf."
}

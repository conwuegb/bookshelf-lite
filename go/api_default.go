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

// HANDLER FUNCTIONS

// Index handles any requests to the root url.
// It takes no input and responds with a standard welcome message.
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, and welcome to your bookshelf!")
}

// BooksGet handles a GET request to the /books endpoint.
// It expects no input and responds with a JSON object representing
// all the books on the bookshelf.
func BooksGet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	// Fetch books from DB and encode to JSON in response
	json.NewEncoder(w).Encode(getAllBooks())
}

// BooksPost handles a POST request to the /books endpoint.
// It expects a request body containing a complete Book represented
// in either JSON or a URL-encoding.
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

// BooksDeleteByIsbn handles a DELETE request to the /books/{isbn} endpoint.
// It expects a request param containing the ISBN of the book to be deleted.
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

// DATABASE FUNCTIONS

// getAllBooks fetches the full collection of Book objects.
func getAllBooks() *[]Book {
	return GetDB()
}

// createBook takes a pointer to a Book object and writes it to the database.
// It returns an HTTP Status Code and an optional message for the response body.
// In the case that the book is already in the shelf, the status code returned
// will be 200 (OK). Otherwise the status code will be 204 (Created).
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

// deleteBook takes the ISBN of the book to be deleted as a string and removes
// it from the database if it can be found.
// It returns an HTTP Status Code and an optional message for the response body.
// In the case that the book is found on the shelf, it returns
// a status code of 204 (No Content) to indicate the deletion.
// Otherwise the status code will be 404 (Not Found).
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

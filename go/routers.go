/*
 * Bookshelf API
 *
 * An API which manages the inventory of books on a bookshelf.
 *
 */
package swagger

import (
	"github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)

	router.HandleFunc("/", Index).Methods("GET")
	router.HandleFunc("/books", BooksGet).Methods("GET")
	router.HandleFunc("/books", BooksPost).Methods("POST")

	return router
}

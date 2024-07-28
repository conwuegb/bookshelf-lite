package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sw "bookshelf/go"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

// Stub data
var books = []sw.Book{
	{Isbn: "9781627792127", Title: "Six of Crows", Author: "Leigh Bardugo"},
	{Isbn: "9781785652509", Title: "The Invisible Life of Addie Larue", Author: "V. E. Schwab"},
	{Isbn: "9781250105714", Title: "Sadie", Author: "Courtney Summers"},
	{Isbn: "9788408045076", Title: "El Alquimista", Author: "Paulo Coehlo"},
}

var test_bookshelf []sw.Book = nil

// Mock the database
func setupTestDB(baseBooks ...sw.Book) {
	sw.GetDB = func() *[]sw.Book {
		// Initialize or reset db
		test_bookshelf = []sw.Book{}
		// Add books for current test
		test_bookshelf = append(test_bookshelf, baseBooks...)
		return &test_bookshelf
	}
}

// Helper function to set up the router
func setupRouter() *mux.Router {
	return sw.NewRouter()
}

func TestGetBooks_CleanDB(t *testing.T) {
	router := setupRouter()

	setupTestDB()

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	expected, _ := json.Marshal([]sw.Book{})
	actual := strings.TrimRight(rr.Body.String(), "\n")
	assert.Equal(t, string(expected), actual)
}

func TestGetBooks_ExistingData(t *testing.T) {
	router := setupRouter()

	setupTestDB(books[2], books[3])

	req, _ := http.NewRequest("GET", "/books", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "Sadie")
	assert.Contains(t, rr.Body.String(), "El Alquimista")
}

func TestCreateBook_NewBook(t *testing.T) {
	router := setupRouter()

	setupTestDB(books...)

	// Build Request
	newBook := sw.Book{Isbn: "9780735211292", Title: "Atomic Habits", Author: "James Clear"}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Serve & Record Request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusCreated, rr.Code)
	assert.Contains(t, rr.Body.String(), "Atomic Habits")
}

func TestCreateBook_DuplicateBook(t *testing.T) {
	router := setupRouter()

	setupTestDB(books...)

	// Build Request
	newBook := sw.Book{Isbn: "9781785652509", Title: "The Invisible Life of Addie Larue", Author: "V. E. Schwab"}
	body, _ := json.Marshal(newBook)
	req, _ := http.NewRequest("POST", "/books", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// Serve & Record Request
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "already on your shelf")
	assert.Contains(t, rr.Body.String(), "Addie Larue")
}

func TestDeleteBookByIsbn_BookExists(t *testing.T) {
	router := setupRouter()

	setupTestDB(books...)

	// Delete from start edge
	req, _ := http.NewRequest("DELETE", "/books/9781627792127", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body)

	// Delete from middle
	req, _ = http.NewRequest("DELETE", "/books/9781250105714", nil)
	rr = httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNoContent, rr.Code)
	assert.Empty(t, rr.Body)
}

func TestDeleteBookByIsbn_BookDoesNotExists(t *testing.T) {
	router := setupRouter()

	setupTestDB(books[1:]...)

	req, _ := http.NewRequest("DELETE", "/books/9781627792127", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusNotFound, rr.Code)
	assert.Contains(t, rr.Body.String(), "not found")
}

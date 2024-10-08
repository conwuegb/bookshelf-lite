/*
 * Bookshelf API
 *
 * An API which manages the inventory of books on a bookshelf.
 *
 * API version: 1.0
 * Generated by: Swagger Codegen (https://github.com/swagger-api/swagger-codegen.git)
 */
package swagger

type Book struct {

	Author string `json:"author"`

	Isbn string `json:"isbn"`

	Title string `json:"title"`
}

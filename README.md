# Bookshelf REST API in Go

An lightweight API which manages the inventory of books on a bookshelf.

## Overview
The API is specified with OpenAPI spec, which can be found in `/api`.

It allows a user to view, add, and remove books on a bookshelf. 

Books are specified as JSON containing the books title, author, and isbn13 number (no special characters).

Example:
```json
{
  "author": "Fyodor Dostoyevsky",
  "isbn": "9780679734529",
  "title": "Notes From the Underground"
}
```

### Running the server
To run the server, first build it using:

```bash
$ go build .
```

Then run using:

```bash
$ go run .
```
or 
```bash
$ ./bookshelf
```

### Limitations
- Currently, the server does not persist data between runs.
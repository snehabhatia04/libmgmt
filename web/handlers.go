// package handlers

// import (
// 	"encoding/json"
// 	"net/http"
// 	"strconv"

//     "github.com/snehabhatia04/libmgmt/model"
// )

// type Handler struct {
// 	UserStore       model.UserStore
// 	BookStore       model.BookStore
// 	IssuedBookStore model.IssuedBookStore
// 	FineStore       model.FineStore
// }

// // User Handlers
// func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
// 	user, err := h.UserStore.GetUser(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(user)
// }

// func (h *Handler) GetUsers(w http.ResponseWriter, r *http.Request) {
// 	users, err := h.UserStore.GetUsers()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(users)
// }

// // Book Handlers
// func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
// 	book, err := h.BookStore.GetBook(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(book)
// }

// func (h *Handler) GetBooks(w http.ResponseWriter, r *http.Request) {
// 	books, err := h.BookStore.GetBooks()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(books)
// }

// // Book Issue Handlers
// func (h *Handler) IssueBook(w http.ResponseWriter, r *http.Request) {
// 	var issue model.BookIssue
// 	if err := json.NewDecoder(r.Body).Decode(&issue); err != nil {
// 		http.Error(w, "Invalid input", http.StatusBadRequest)
// 		return
// 	}
// 	if err := h.IssuedBookStore.IssueBook(&issue); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusCreated)
// }

// func (h *Handler) ReturnBook(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(r.URL.Query().Get("book_id"))
// 	fine, err := h.IssuedBookStore.ReturnBook(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(map[string]float64{"fine": fine})
// }

// // Fine Handlers
// func (h *Handler) GetFine(w http.ResponseWriter, r *http.Request) {
// 	id, _ := strconv.Atoi(r.URL.Query().Get("id"))
// 	fine, err := h.FineStore.GetFine(id)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(fine)
// }

// func (h *Handler) GetFines(w http.ResponseWriter, r *http.Request) {
// 	fines, err := h.FineStore.GetFines()
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	json.NewEncoder(w).Encode(fines)
// }

package web

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	// "github.com/snehabhatia04/libmgmt/model"

	"github.com/snehabhatia04/libmgmt/model"
)

// Handler struct to interact with storage interfaces
type Handler struct {
	UserStore       model.UserStore
	BookStore       model.BookStore
	IssuedBookStore model.IssuedBookStore
	FineStore       model.FineStore
}


// NewHandler initializes a new Handler instance
func NewHandler(userStore model.UserStore, bookStore model.BookStore, issuedBookStore model.IssuedBookStore, fineStore model.FineStore) *Handler {
	return &Handler{
		UserStore:       userStore,
		BookStore:       bookStore,
		IssuedBookStore: issuedBookStore,
		FineStore:       fineStore,
	}
}
// GetUsers retrieves all users
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.UserStore.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}

// GetUser retrieves a single user by ID
func (h *Handler) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	user, err := h.UserStore.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUser adds a new user
func (h *Handler) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.UserStore.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// DeleteUser removes a user by ID
func (h *Handler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.UserStore.DeleteUser(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}

// GetBooks retrieves all books
func (h *Handler) GetBooks(c *gin.Context) {
	books, err := h.BookStore.GetBooks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}

// GetBook retrieves a book by ID
func (h *Handler) GetBook(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	book, err := h.BookStore.GetBook(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, book)
}

// CreateBook handles adding a new book
func (h *Handler) CreateBook(c *gin.Context) {
	var book model.Book
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	err := h.BookStore.CreateBook(&book)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create book"})
		return
	}

	c.JSON(201, gin.H{"message": "Book created successfully", "book": book})
}


// UpdateBook handles updating a book
func (h *Handler) UpdateBook(c *gin.Context) {
	//bookID := c.Param("id") // Get book ID from URL parameter
	var book model.Book

	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	err := h.BookStore.UpdateBook(&book)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(200, gin.H{"message": "Book updated successfully", "book": book})
}


// IssueBook issues a book to a user
func (h *Handler) IssueBook(c *gin.Context) {
	var issue model.BookIssue
	if err := c.ShouldBindJSON(&issue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.IssuedBookStore.IssueBook(&issue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, issue)
}

// ReturnBook processes a book return and calculates fine
func (h *Handler) ReturnBook(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("book_id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}
	fine, err := h.IssuedBookStore.ReturnBook(bookID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"fine": fine})
}

// DeleteBook handles deleting a book
func (h *Handler) DeleteBook(c *gin.Context) {
	bookID, err := strconv.Atoi(c.Param("id")) // Get book ID from URL
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid book ID"})
		return
	}

	err = h.BookStore.DeleteBook(bookID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(200, gin.H{"message": "Book deleted successfully"})
}


// GetFines retrieves all fines
func (h *Handler) GetFines(c *gin.Context) {
	fines, err := h.FineStore.GetFines()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fines)
}

// GetFine retrieves a fine by ID
func (h *Handler) GetFine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	fine, err := h.FineStore.GetFine(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, fine)
}

// DeleteFine removes a fine by ID
func (h *Handler) DeleteFine(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	if err := h.FineStore.DeleteFine(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Fine deleted successfully"})
}

// // RegisterRoutes sets up routes for the handlers
// func (h *Handler) RegisterRoutes(router *gin.Engine) {
// 	api := router.Group("/api")
// 	{
// 		api.GET("/users", h.GetUsers)
// 		api.GET("/users/:id", h.GetUser)
// 		api.POST("/users", h.CreateUser)
// 		api.DELETE("/users/:id", h.DeleteUser)

// 		api.GET("/books", h.GetBooks)
// 		api.GET("/books/:id", h.GetBook)
// 		api.POST("/books", h.CreateBook)
// 		api.PUT("/books/:id", h.UpdateBook)
// 		api.DELETE("/books/:id", h.DeleteBook)

// 		api.POST("/issue", h.IssueBook)
// 		api.PUT("/return/:book_id", h.ReturnBook)

// 		api.GET("/fines", h.GetFines)
// 		api.GET("/fines/:id", h.GetFine)
// 		api.DELETE("/fines/:id", h.DeleteFine)
// 	}
// }
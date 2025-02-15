package model

// User struct
type User struct {
    ID   int    `json:"id"`
    Name string `json:"name"`
}

// Book struct
type Book struct {
    ID       int    `json:"id"`
    Title    string `json:"title"`
    AuthorID int    `json:"author_id"`
    Location string `json:"location"`
}

// BookIssue struct
type BookIssue struct {
    ID         int    `json:"id"`
    UserID     int    `json:"user_id"`
    BookID     int    `json:"book_id"`
    IssueDate  string `json:"issue_date"`
    ReturnDate string `json:"return_date,omitempty"`
    DueDate    string `json:"due_date"`
}

// Fine struct
type Fine struct {
    ID         int     `json:"id"`
    UserID     int     `json:"user_id"`
    FineAmount float64 `json:"fine_amount"`
    Date       string  `json:"date"`
}

// UserStore defines methods for user management
type UserStore interface {
    GetUser(id int) (User, error)
    GetUsers() ([]User, error)
    CreateUser(user *User) error
    UpdateUser(user *User) error
    DeleteUser(id int) error
}

// BookStore defines methods for book management
type BookStore interface {
    GetBook(id int) (Book, error)
    GetBooks() ([]Book, error)
    CreateBook(book *Book) error
    UpdateBook(book *Book) error
    DeleteBook(id int) error
}

// IssuedBookStore defines methods for issuing and returning books
type IssuedBookStore interface {
    IssueBook(issue *BookIssue) error
    ReturnBook(bookID int) (float64, error) // Returns fine amount if applicable
    GetIssuedBookByBookID(bookID int) (BookIssue, error)
    GetIssuedBooks() ([]BookIssue, error)
}

// FineStore defines methods for managing fines
type FineStore interface {
    AddFine(fine *Fine) error
    GetFine(id int) (Fine, error)
    GetFines() ([]Fine, error)
    DeleteFine(id int) error
}

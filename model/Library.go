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

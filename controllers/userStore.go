package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/snehabhatia04/libmgmt/model"
)

// UserStore struct
type UserStore struct {
	DB *sqlx.DB
}

// Ensure UserStore implements model.UserStore
var _ model.UserStore = (*UserStore)(nil)

// GetUser retrieves a user by ID
func (us *UserStore) GetUser(id int) (model.User, error) {
	var user model.User
	err := us.DB.Get(&user, "SELECT * FROM users WHERE id=$1", id)
	return user, err
}

// GetUsers retrieves all users
func (us *UserStore) GetUsers() ([]model.User, error) {
	var users []model.User
	err := us.DB.Select(&users, "SELECT * FROM users")
	return users, err
}

// CreateUser adds a new user
func (us *UserStore) CreateUser(user *model.User) error {
	_, err := us.DB.Exec("INSERT INTO users (name) VALUES ($1)", user.Name)
	return err
}

// UpdateUser updates user details
func (us *UserStore) UpdateUser(user *model.User) error {
	_, err := us.DB.Exec("UPDATE users SET name=$1 WHERE id=$2", user.Name, user.ID)
	return err
}

// DeleteUser removes a user
func (us *UserStore) DeleteUser(id int) error {
	_, err := us.DB.Exec("DELETE FROM users WHERE id=$1", id)
	return err
}

// NewDBUserStore initializes a new UserStore
func NewDBUserStore(db *sqlx.DB) *UserStore {
	return &UserStore{DB: db}
}

type UserController struct {
	Store model.UserStore
}

func NewUserController(store model.UserStore) *UserController {
	return &UserController{Store: store}
}

func (uc *UserController) GetUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	user, err := uc.Store.GetUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (uc *UserController) GetUsers(c *gin.Context) {
	users, err := uc.Store.GetUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uc *UserController) CreateUser(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	err := uc.Store.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, user)
}

func (uc *UserController) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	err = uc.Store.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

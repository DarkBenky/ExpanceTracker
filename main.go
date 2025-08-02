package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	_ "github.com/mattn/go-sqlite3"
)

type Expense struct {
	ID           int     `json:"id" db:"id"`
	Description  string  `json:"description" db:"description"`
	Amount       float64 `json:"amount" db:"amount"`
	Category     string  `json:"category" db:"category"`
	Date         string  `json:"date" db:"date"`
	OwnerGroupID int     `json:"owner_group_id" db:"owner_group_id"`
}

type UsersGroup struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name"`
	OwnerID  int    `json:"owner_id" db:"owner_id"`   // ID of the user who owns the group
	UsersIDs []int  `json:"users_ids" db:"users_ids"` // List of user IDs in the group
}

type User struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}

var db *sql.DB

func initDB() {
	var err error
	db, err = sql.Open("sqlite3", "./expenses.db")
	if err != nil {
		log.Fatal(err)
	}

	// Crete Users table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT UNIQUE NOT NULL,
		password TEXT NOT NULL
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Create Expenses table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS expenses (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		amount REAL NOT NULL,
		category TEXT NOT NULL,
		date TEXT NOT NULL,
		owner_group_id INTEGER NOT NULL,
		FOREIGN KEY (owner_group_id) REFERENCES users_groups(id)
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Create UsersGroups table
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users_groups (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		owner_id INTEGER NOT NULL,
		FOREIGN KEY (owner_id) REFERENCES users(id)
	)`)
	if err != nil {
		log.Fatal(err)
	}

	// Create group_members table for many-to-many relationship
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS group_members (
		group_id INTEGER NOT NULL,
		user_id INTEGER NOT NULL,
		PRIMARY KEY (group_id, user_id),
		FOREIGN KEY (group_id) REFERENCES users_groups(id),
		FOREIGN KEY (user_id) REFERENCES users(id)
	)`)
	if err != nil {
		log.Fatal(err)
	}
}

func isUserInGroup(userID, groupID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM group_members WHERE user_id = ? AND group_id = ?)", userID, groupID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func userExists(userID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id = ?)", userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type AddExpenseRequest struct {
	Token       string  `json:"token"`    // JWT token for authentication
	UserID      int     `json:"user_id"`  // ID of the user adding the expense
	GroupID     int     `json:"group_id"` // ID of the group to which the expense belongs
	Description string  `json:"description"`
	Amount      float64 `json:"amount"`
	Category    string  `json:"category"`
	Date        string  `json:"date"`
}

func AddExpense(c echo.Context) error {
	var req AddExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.Description == "" || req.Amount <= 0 || req.Category == "" || req.Date == "" || req.GroupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "All fields are required"})
	}

	// check if the user exists
	exists, err := userExists(req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// check token validity (this is a placeholder, actual JWT validation should be implemented)
	TODO

	// Check if the user is part of the group
	isMember, err := isUserInGroup(req.UserID, req.GroupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !isMember {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User is not part of the group"})
	}

	// Insert the expense into the database
	_, err = db.Exec(`INSERT INTO expenses (description, amount, category, date, owner_group_id) VALUES (?, ?, ?, ?, ?)`,
		req.Description, req.Amount, req.Category, req.Date, req.GroupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add expense"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Expense added successfully"})
}

func RemoveExpense(c echo.Context) error {

func main() {
	initDB()

}

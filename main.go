package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

const salt = "1afdcf647dbd07d353987550168"

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

// JWT secret key
var jwtSecret = []byte("380cde94f76c8d37d6c3946dff11dffdc5f469bdca3aff5931f1ae1a445856a9776c7624bfbdd8017806c64af182337e51a10c746321ac689cafef2a5bc1e90e")

// JWT Claims structure
type JWTClaims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// Generate JWT token for user
func generateToken(userID int) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Token expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// Validate JWT token and extract user ID
func validateToken(tokenString string) (int, error) {
	if tokenString == "" {
		return 0, errors.New("token is required")
	}

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims.UserID, nil
	}

	return 0, errors.New("invalid token")
}

// Hash password using bcrypt
func hashPassword(password string) (string, error) {
	// Add salt to password before hashing
	saltedPassword := password + salt

	// Generate bcrypt hash with cost 12 (recommended for production)
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), 12)
	if err != nil {
		return "", err
	}

	return string(hashedBytes), nil
}

// Verify password using bcrypt
func verifyPassword(password, hashedPassword string) bool {
	// Add salt to password before verification
	saltedPassword := password + salt

	// Compare password with hash
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
	return err == nil
}

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

	// Validate JWT token and get user ID
	userID, err := validateToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
	}

	// check if the user exists
	exists, err := userExists(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Check if the user is part of the group
	isMember, err := isUserInGroup(userID, req.GroupID)
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


func checkExpenseExists(expenseID int) (bool, error) {
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM expenses WHERE id = ?)", expenseID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

type removeExpenseRequest struct {
	Token     string `json:"token"`      // JWT token for authentication
	GroupID   int    `json:"group_id"`   // ID of the group to which the expense belongs
	ExpenseID int    `json:"expense_id"` // ID of the expense to be removed
}


func RemoveExpense(c echo.Context) error {
	fmt.Println("RemoveExpense called")
	var req removeExpenseRequest
	if err := c.Bind(&req); err != nil {
		fmt.Println("Error binding request:", err)
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.ExpenseID <= 0 || req.GroupID <= 0 {
		fmt.Println("Invalid request: Expense ID and Group ID are required")
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Expense ID and Group ID are required"})
	}

	// Validate JWT token and get user ID
	userID, err := validateToken(req.Token)
	if err != nil {
		fmt.Println("Error validating token:", err)
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
	}

	// check if the user exists
	exists, err := userExists(userID)
	if err != nil {
		fmt.Println("Error checking user existence:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		fmt.Println("User not found")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	// Check if the user is part of the group
	isMember, err := isUserInGroup(userID, req.GroupID)
	if err != nil {
		fmt.Println("Error checking group membership:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !isMember {
		fmt.Println("User is not part of the group")
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User is not part of the group"})
	}

	// Check if the expense exists
	exists, err = checkExpenseExists(req.ExpenseID)
	if err != nil {
		fmt.Println("Error checking expense existence:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		fmt.Println("Expense not found")
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Expense not found"})
	}

	// Delete the expense from the database
	_, err = db.Exec(`DELETE FROM expenses WHERE id = ?`, req.ExpenseID)
	if err != nil {
		fmt.Println("Error deleting expense:", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to remove expense"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Expense removed successfully"})
}

// User Registration
type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
	}

	// Validate password strength (optional but recommended)
	if len(req.Password) < 8 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Password must be at least 8 characters long"})
	}

	// Check if username already exists
	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)", req.Username).Scan(&exists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"error": "Username already exists"})
	}

	// Hash the password using bcrypt
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to process password"})
	}

	// Insert new user with hashed password
	result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", req.Username, hashedPassword)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	userID, _ := result.LastInsertId()
	token, err := generateToken(int(userID))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
		"user_id": userID,
		"token":   token,
	})
}

// User Login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Username and password are required"})
	}

	// Check user credentials
	var userID int
	var storedHashedPassword string
	err := db.QueryRow("SELECT id, password FROM users WHERE username = ?", req.Username).Scan(&userID, &storedHashedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	// Verify password using bcrypt
	if !verifyPassword(req.Password, storedHashedPassword) {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	// Generate JWT token
	token, err := generateToken(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Login successful",
		"user_id": userID,
		"token":   token,
	})
}

type CreateGroupRequest struct {
	Token string `json:"token"` // JWT token for authentication
	Name  string `json:"name"`  // Name of the group
}

func CreateGroup(c echo.Context) error {
	var req CreateGroupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group name is required"})
	}

	// Validate JWT token
	userID, err := validateToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
	}

	// Insert new group into the database
	result, err := db.Exec("INSERT INTO users_groups (name, owner_id) VALUES (?, ?)", req.Name, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create group"})
	}

	groupID, _ := result.LastInsertId()

	// Add the owner to the group as a member
	_, err = db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add owner to group"})
	}

	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message":  "Group created successfully",
		"group_id": groupID,
	})
}

type AddUserToGroupRequest struct {
	Token   string `json:"token"`    // JWT token for authentication
	GroupID int    `json:"group_id"` // ID of the group to which the user will be added
	UserID  int    `json:"user_id"`  // ID of the user to be added to the group
}

func AddUserToGroup(c echo.Context) error {
	var req AddUserToGroupRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate required fields
	if req.GroupID <= 0 || req.UserID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Group ID and User ID are required"})
	}

	// Validate JWT token
	validUserID, err := validateToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
	}

	// Check if the user is the owner of the group
	var ownerID int
	err = db.QueryRow("SELECT owner_id FROM users_groups WHERE id = ?", req.GroupID).Scan(&ownerID)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Group not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	if ownerID != validUserID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Only the group owner can add users"})
	}

	// Check if the user already exists in the group
	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM group_members WHERE group_id = ? AND user_id = ?)", req.GroupID, req.UserID).Scan(&exists)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if exists {
		return c.JSON(http.StatusConflict, map[string]string{"error": "User already exists in the group"})
	}

	// Add the user to the group
	_, err = db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", req.GroupID, req.UserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to add user to group"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User added to group successfully"})
}

type GetExpensesRequest struct {
	Token   string `json:"token"`              // JWT token for authentication
	GroupID int    `json:"group_id,omitempty"` // Optional: filter by group, -1 for all groups
}

func GetExpenses(c echo.Context) error {
	var req GetExpensesRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate JWT token
	UserID, err := validateToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
	}

	// Build query based on group filter
	var query string
	var args []interface{}

	if req.GroupID > 0 {
		// Check if user is member of the specific group
		isMember, err := isUserInGroup(UserID, req.GroupID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		if !isMember {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "User is not part of the group"})
		}

		// Get expenses from specific group
		query = `SELECT e.id, e.description, e.amount, e.category, e.date, e.owner_group_id 
                 FROM expenses e 
                 WHERE e.owner_group_id = ?
                 ORDER BY e.date DESC`
		args = []interface{}{req.GroupID}
	} else {
		// Get all expenses from groups where user is a member (GroupID = 0, -1, or not provided)
		query = `SELECT e.id, e.description, e.amount, e.category, e.date, e.owner_group_id 
                 FROM expenses e 
                 INNER JOIN group_members gm ON e.owner_group_id = gm.group_id 
                 WHERE gm.user_id = ?
                 ORDER BY e.date DESC`
		args = []interface{}{UserID}
	}

	rows, err := db.Query(query, args...)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	var expenses []Expense
	for rows.Next() {
		var expense Expense
		err := rows.Scan(&expense.ID, &expense.Description, &expense.Amount,
			&expense.Category, &expense.Date, &expense.OwnerGroupID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		expenses = append(expenses, expense)
	}

	return c.JSON(http.StatusOK, expenses)
}

type UpdateExpenseRequest struct {
	Token       string  `json:"token"`       // JWT token for authentication
	GroupID     int     `json:"group_id"`    // ID of the group to which the expense belongs
	ExpenseID   int     `json:"expense_id"`  // ID of the expense to be updated
	Description string  `json:"description"` // Updated description of the expense
	Amount      float64 `json:"amount"`      // Updated amount of the expense
	Category    string  `json:"category"`    // Updated category of the expense
	Date        string  `json:"date"`        // Updated date of the expense
}

func UpdateExpense(c echo.Context) error {
	var req UpdateExpenseRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	// Validate required fields
	if req.ExpenseID <= 0 || req.GroupID <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Expense ID and Group ID are required"})
	}
	// Validate JWT token and get user ID
	userID, err := validateToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
	}
	// check if the user exists
	exists, err := userExists(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}
	// Check if the user is part of the group
	isMember, err := isUserInGroup(userID, req.GroupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !isMember {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "User is not part of the group"})
	}

	// Check if the expense exists
	exists, err = checkExpenseExists(req.ExpenseID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Expense not found"})
	}
	// Update the expense in the database
	_, err = db.Exec(`UPDATE expenses SET description = ?, amount = ?, category = ?, date = ?
		WHERE id = ? AND owner_group_id = ?`, req.Description, req.Amount, req.Category, req.Date, req.ExpenseID, req.GroupID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update expense"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Expense updated successfully"})
}

type GetGroupsRequest struct {
	Token string `json:"token"` // JWT token for authentication
}

func GetUserGroups(c echo.Context) error {
	var req GetGroupsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate JWT token
	validUserID, err := validateToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
	}

	// Get all groups where user is a member
	query := `SELECT ug.id, ug.name, ug.owner_id 
              FROM users_groups ug 
              INNER JOIN group_members gm ON ug.id = gm.group_id 
              WHERE gm.user_id = ?
              ORDER BY ug.name`

	rows, err := db.Query(query, validUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	var groups []UsersGroup
	for rows.Next() {
		var group UsersGroup
		err := rows.Scan(&group.ID, &group.Name, &group.OwnerID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		groups = append(groups, group)
	}

	return c.JSON(http.StatusOK, groups)
}

type ValidateUserTokenRequest struct {
	Token string `json:"token"` // JWT token for authentication
}

func ValidateUserToken(c echo.Context) error {
	var req ValidateUserTokenRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request", "valid": "false"})
	}

	// Validate JWT token
	userID, err := validateToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error(), "valid": "false"})
	}

	// Check if user exists
	exists, err := userExists(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error", "valid": "false"})
	}
	if !exists {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{"message": "Token is valid", "user_id": userID, "valid": "true"})
}

func populateFakeData() error {
	// Check if data already exists to avoid duplicates
	var userCount int
	err := db.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		return err
	}

	// If users already exist, skip population
	if userCount > 0 {
		log.Println("Database already contains data, skipping fake data population")
		return nil
	}

	log.Println("Populating database with fake data...")

	// Create fake users
	fakeUsers := []struct {
		username string
		password string
	}{
		{"alice", "password123"},
		{"bob", "password123"},
		{"charlie", "password123"},
		{"diana", "password123"},
		{"eve", "password123"},
	}

	var userIDs []int
	for _, user := range fakeUsers {
		hashedPassword, err := hashPassword(user.password)
		if err != nil {
			fmt.Println("Error hashing password for user:", user.username, err)
			return err
		}

		result, err := db.Exec("INSERT INTO users (username, password) VALUES (?, ?)", user.username, hashedPassword)
		if err != nil {
			fmt.Println("Error inserting user:", user.username, err)
			return err
		}

		userID, _ := result.LastInsertId()
		userIDs = append(userIDs, int(userID))
	}

	// Create fake groups
	fakeGroups := []struct {
		name    string
		ownerID int
		members []int
	}{
		{"Family", userIDs[0], []int{userIDs[0], userIDs[1], userIDs[2]}},
		{"Work Team", userIDs[1], []int{userIDs[1], userIDs[3], userIDs[4]}},
		{"Friends", userIDs[2], []int{userIDs[0], userIDs[2], userIDs[3]}},
		{"Roommates", userIDs[3], []int{userIDs[3], userIDs[4]}},
	}

	var groupIDs []int
	for _, group := range fakeGroups {
		result, err := db.Exec("INSERT INTO users_groups (name, owner_id) VALUES (?, ?)", group.name, group.ownerID)
		if err != nil {
			return err
		}

		groupID, _ := result.LastInsertId()
		groupIDs = append(groupIDs, int(groupID))

		// Add members to group
		for _, memberID := range group.members {
			_, err = db.Exec("INSERT INTO group_members (group_id, user_id) VALUES (?, ?)", groupID, memberID)
			if err != nil {
				return err
			}
		}
	}

	// Create fake expenses
	fakeExpenses := []struct {
		description string
		amount      float64
		category    string
		date        string
		groupID     int
	}{
		{"Grocery shopping", 85.50, "Food", "2024-01-15", groupIDs[0]},
		{"Coffee meeting", 12.30, "Food", "2024-01-16", groupIDs[1]},
		{"Movie tickets", 45.00, "Entertainment", "2024-01-17", groupIDs[2]},
		{"Electricity bill", 120.75, "Utilities", "2024-01-18", groupIDs[3]},
		{"Gas station", 60.00, "Transportation", "2024-01-19", groupIDs[0]},
		{"Restaurant dinner", 78.90, "Food", "2024-01-20", groupIDs[1]},
		{"Netflix subscription", 15.99, "Entertainment", "2024-01-21", groupIDs[2]},
		{"Internet bill", 55.00, "Utilities", "2024-01-22", groupIDs[3]},
		{"Bus pass", 25.00, "Transportation", "2024-01-23", groupIDs[0]},
		{"Lunch", 18.50, "Food", "2024-01-24", groupIDs[1]},
		{"Concert tickets", 95.00, "Entertainment", "2024-01-25", groupIDs[2]},
		{"Water bill", 45.25, "Utilities", "2024-01-26", groupIDs[3]},
		{"Uber ride", 22.80, "Transportation", "2024-01-27", groupIDs[0]},
		{"Pizza order", 35.60, "Food", "2024-01-28", groupIDs[1]},
		{"Gaming subscription", 9.99, "Entertainment", "2024-01-29", groupIDs[2]},
	}

	for _, expense := range fakeExpenses {
		_, err = db.Exec("INSERT INTO expenses (description, amount, category, date, owner_group_id) VALUES (?, ?, ?, ?, ?)",
			expense.description, expense.amount, expense.category, expense.date, expense.groupID)
		if err != nil {
			return err
		}
	}

	log.Println("Fake data populated successfully!")
	return nil
}

type GetGroupMembersRequest struct {
    Token   string `json:"token"`    // JWT token for authentication
    GroupID int    `json:"group_id"` // ID of the group to get members for
}

func GetGroupMembers(c echo.Context) error {
    var req GetGroupMembersRequest
    if err := c.Bind(&req); err != nil {
        return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
    }

    // Validate JWT token
    userID, err := validateToken(req.Token)
    if err != nil {
        return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
    }

    // Check if user is member of the group
    isMember, err := isUserInGroup(userID, req.GroupID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    if !isMember {
        return c.JSON(http.StatusForbidden, map[string]string{"error": "User is not part of the group"})
    }

    // Get all users in the group
    query := `SELECT u.id, u.username 
              FROM users u 
              INNER JOIN group_members gm ON u.id = gm.user_id 
              WHERE gm.group_id = ?
              ORDER BY u.username`

    rows, err := db.Query(query, req.GroupID)
    if err != nil {
        return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
    }
    defer rows.Close()

    var users []struct {
        ID       int    `json:"id"`
        Username string `json:"username"`
    }

    for rows.Next() {
        var user struct {
            ID       int    `json:"id"`
            Username string `json:"username"`
        }
        err := rows.Scan(&user.ID, &user.Username)
        if err != nil {
            return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
        }
        users = append(users, user)
    }

    return c.JSON(http.StatusOK, map[string]interface{}{
        "users": users,
    })
}

func GetUsers(c echo.Context) error {
	var req GetGroupsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	// Validate JWT token
	validUserID, err := validateToken(req.Token)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid token: " + err.Error()})
	}

	// Get all users
	rows, err := db.Query("SELECT id, username FROM users WHERE id != ?", validUserID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Username)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Database error"})
	}

	fmt.Println("Returning users:", users)
	return c.JSON(http.StatusOK, users)
}

func Ping(c echo.Context) error {
	return c.String(http.StatusOK, "pong")
}

func main() {
	initDB()
	defer db.Close()

	// // Populate fake data if database is empty
	// if err := populateFakeData(); err != nil {
	// 	log.Printf("Error populating fake data: %v", err)
	// }

	e := echo.New()

	// Add CORS middleware with more permissive settings
    e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
        AllowOrigins: []string{
            "*", // Allow all origins for development
        },
        AllowMethods: []string{
            http.MethodGet, 
            http.MethodPost, 
            http.MethodPut, 
            http.MethodDelete, 
            http.MethodOptions,
        },
        AllowHeaders: []string{
            echo.HeaderOrigin, 
            echo.HeaderContentType, 
            echo.HeaderAccept, 
            echo.HeaderAuthorization,
            "X-Requested-With",
            "Access-Control-Allow-Origin",
        },
        AllowCredentials: true,
    }))

    // Add logging middleware to see requests
    e.Use(middleware.Logger())

	// Routes
	e.GET("/ping", Ping)
	e.POST("/register", Register)
	e.POST("/login", Login)
	e.POST("/expenses", AddExpense)
	e.POST("/expenses/get", GetExpenses)
	e.POST("/expenses/update", UpdateExpense)
	e.DELETE("/expenses", RemoveExpense)
	e.POST("/groups", CreateGroup)
	e.POST("/groups/get", GetUserGroups)
	e.POST("/groups/members/get", GetGroupMembers)
	e.POST("/groups/members/add", AddUserToGroup)
	e.POST("/validate/token", ValidateUserToken)
	e.POST("/users/get", GetUsers)

	// Start server
	log.Println("Server starting on :1234")
	log.Fatal(e.Start(":1234"))
}

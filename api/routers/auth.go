package routers

import (
	"database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"os"
	"regexp"
	"strings"
)

func RegisterAuthRoutes(router *gin.Engine) {
	router.POST("/auth/login", login)
	router.GET("/auth/status", authStatus)
	router.Group("/", AuthRequired).POST("/auth/logout", logout)
	router.POST("/auth/register", register)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func checkUserTable() bool {
	// Check if user table exists
	stmt, err := db.Prepare("SHOW TABLES LIKE 'users'")
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var table string
	err = stmt.QueryRow().Scan(&table)
	return err == nil
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func createUserTable() {
	if checkUserTable() {
		return
	}
	// create table
	stmt, err := db.Prepare("CREATE TABLE users (id VARCHAR(12) PRIMARY KEY, username VARCHAR(255), password VARCHAR(255), user_type VARCHAR(255), real_name VARCHAR(255), email VARCHAR(255), address VARCHAR(255), phone_number VARCHAR(255))")
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec()
	if err != nil {
		panic(err.Error())
	}
	createAdminUser()
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func createAdminUser() {
	id := createRandomID()
	for checkUserID(id) {
		id = createRandomID()
	}
	stmt, err := db.Prepare("INSERT INTO users (id, username, password, user_type, real_name, email, address, phone_number) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	adminUserName := os.Getenv("ADMIN_USERNAME")
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	_, err = stmt.Exec(id, adminUserName, hashPassword(adminPassword), "admin", "admin", "", "", "")
	if err != nil {
		panic(err.Error())
	}
}

func authStatus(c *gin.Context) {
	// Check if user is logged in
	session := sessions.Default(c)
	user := session.Get("user")
	if user != nil && checkUserID(user.(string)) {
		c.JSON(200, gin.H{
			"message": "OK",
			"login":   true,
		})
		return
	}
	c.JSON(401, gin.H{
		"error": "unauthorized",
		"login": false,
	})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func login(c *gin.Context) {
	var loginForm struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	// check if already logged in
	session := sessions.Default(c)
	user := session.Get("user")
	if user != nil && checkUserID(user.(string)) {
		c.JSON(200, gin.H{"message": "already logged in"})
		return
	}
	if c.ShouldBindJSON(&loginForm) != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	stmt, err := db.Prepare("SELECT id FROM users WHERE username = ? AND password = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var id string
	err = stmt.QueryRow(loginForm.Username, hashPassword(loginForm.Password)).Scan(&id)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid username or password"})
		return
	}
	session.Set("user", id)
	err = session.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "login successful"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func logout(c *gin.Context) {
	// Logout
	session := sessions.Default(c)
	session.Delete("user")
	err := session.Save()
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "logout successful"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func register(c *gin.Context) {
	// Register
	var registerRequest struct {
		Username    string `json:"username"`
		Password    string `json:"password"`
		RealName    string `json:"real_name"`
		Email       string `json:"email"`
		Address     string `json:"address"`
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&registerRequest); err != nil {
		c.JSON(400, gin.H{"error": "invalid request"})
		return
	}
	if checkUserName(registerRequest.Username) {
		c.JSON(400, gin.H{"error": "username already exists"})
		return
	}
	// check email format
	if !strings.Contains(registerRequest.Email, "@") || !strings.Contains(registerRequest.Email, ".") {
		c.JSON(400, gin.H{"error": "invalid email format"})
		return
	}
	// check phone number format(only +, 0-9)
	regx := regexp.MustCompile(`^[+0-9]+$`)
	if !regx.MatchString(registerRequest.PhoneNumber) {
		c.JSON(400, gin.H{"error": "invalid phone number format"})
		return
	}
	stmt, err := db.Prepare("INSERT INTO users (id, username, password, user_type, real_name, email, address, phone_number) VALUES (?, ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var id = createRandomID()
	for checkUserID(id) {
		id = createRandomID()
	}
	_, err = stmt.Exec(id, registerRequest.Username, hashPassword(registerRequest.Password), "user", registerRequest.RealName, registerRequest.Email, registerRequest.Address, registerRequest.PhoneNumber)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(201, gin.H{"message": "user registered"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func checkUserID(id string) bool {
	// Check if user id exists
	stmt, err := db.Prepare("SELECT id FROM users WHERE id = ?")
	if err != nil {
		return false
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var userID string
	err = stmt.QueryRow(id).Scan(&userID)
	return err == nil
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func checkUserName(username string) bool {
	// Check if username exists
	stmt, err := db.Prepare("SELECT username FROM users WHERE username = ?")
	if err != nil {
		return false
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var user string
	err = stmt.QueryRow(username).Scan(&user)
	return err == nil
}

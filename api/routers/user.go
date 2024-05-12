package routers

import (
	"database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

func RegisterUserRoutes(router *gin.Engine) {
	authorized := router.Group("/user", AuthRequired)
	authorized.POST("", editMyUserInformation)
	authorized.GET("", getMyUserInformation)
	authorized.DELETE("", deleteUser)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func editMyUserInformation(c *gin.Context) {
	// Get request body
	var editUserRequest struct {
		RealName    string `json:"real_name"`
		Email       string `json:"email"`
		Address     string `json:"address"`
		OldPassword string `json:"old_password"`
		NewPassword string `json:"new_password"`
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&editUserRequest); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	session := sessions.Default(c)
	user := session.Get("user").(string)
	// Check if the password is correct
	stmt, err := db.Prepare("SELECT id FROM users WHERE id = ? AND password = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var id string
	err = stmt.QueryRow(user, hashPassword(editUserRequest.OldPassword)).Scan(&id)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid password"})
		return
	}

	// check email format
	if !strings.Contains(editUserRequest.Email, "@") || !strings.Contains(editUserRequest.Email, ".") {
		c.JSON(400, gin.H{"error": "invalid email format"})
		return
	}
	// check phone number format(only +, 0-9)
	regx := regexp.MustCompile(`^[+0-9]+$`)
	if !regx.MatchString(editUserRequest.PhoneNumber) {
		c.JSON(400, gin.H{"error": "invalid phone number format"})
		return
	}
	// Change password
	if editUserRequest.NewPassword != "" {
		stmt, err = db.Prepare("UPDATE users SET password = ? WHERE id = ?")
		if err != nil {
			c.JSON(500, gin.H{"error": serverErrorMessage})
			return
		}
		defer func(stmt *sql.Stmt) {
			_ = stmt.Close()
		}(stmt)
		_, err = stmt.Exec(editUserRequest.NewPassword, user)
		if err != nil {
			c.JSON(500, gin.H{"error": serverErrorMessage})
			return
		}
	}
	// Edit user information
	stmt, err = db.Prepare("UPDATE users SET real_name = ?, email = ?, address = ?, phone_number = ? WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(editUserRequest.RealName, editUserRequest.Email, editUserRequest.Address, editUserRequest.PhoneNumber, user)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "User information updated"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getMyUserInformation(c *gin.Context) {
	// Get user information
	session := sessions.Default(c)
	user := session.Get("user").(string)
	stmt, err := db.Prepare("SELECT real_name, email, address, phone_number FROM users WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var realName, email, address, phoneNumber string
	err = stmt.QueryRow(user).Scan(&realName, &email, &address, &phoneNumber)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{
		"user_id":      user,
		"real_name":    realName,
		"email":        email,
		"address":      address,
		"phone_number": phoneNumber,
	})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func deleteUser(c *gin.Context) {
	// Delete user
	session := sessions.Default(c)
	user := session.Get("user").(string)
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(user)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	session.Delete("user")
	_ = session.Save()
	c.JSON(200, gin.H{"message": "User deleted"})
}

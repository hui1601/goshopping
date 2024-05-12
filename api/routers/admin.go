package routers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"regexp"
	"strings"
)

func RegisterAdminRoutes(router *gin.Engine) {
	admin := router.Group("/admin", AuthRequired).Use(AdminRequired)
	admin.POST("/user/:id", editUserInformationAdmin)
	admin.GET("/user/:id", getUserInformationAdmin)
	admin.GET("/user", getAllUsersAdmin)
	admin.DELETE("/user/:id", deleteUserAdmin)
	admin.GET("/user/purchase/:id", getUserPurchaseHistoryAdmin)
	admin.GET("/permission/grant/:id", grantAdmin)
	admin.GET("/permission/revoke/:id", revokeAdmin)
	admin.GET("/permission", getUserPermissionsAdmin)
	admin.GET("/purchase", getPurchaseHistoryAdmin)
	admin.POST("/purchase/:id", editPurchaseAdmin)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func editUserInformationAdmin(c *gin.Context) {
	// Get request body
	var editUserRequest struct {
		RealName    string `json:"real_name"`
		Email       string `json:"email"`
		Address     string `json:"address"`
		PhoneNumber string `json:"phone_number"`
	}
	if err := c.ShouldBindJSON(&editUserRequest); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
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
	// Edit user information
	stmt, err := db.Prepare("UPDATE users SET real_name = ?, email = ?, address = ?, phone_number = ? WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(editUserRequest.RealName, editUserRequest.Email, editUserRequest.Address, editUserRequest.PhoneNumber, c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "User information updated"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getUserInformationAdmin(c *gin.Context) {
	// Get user information
	stmt, err := db.Prepare("SELECT username, real_name, email, address, phone_number FROM users WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var username, realName, email, address, phoneNumber string
	err = stmt.QueryRow(c.Param("id")).Scan(&username, &realName, &email, &address, &phoneNumber)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"username": username, "real_name": realName, "email": email, "address": address, "phone_number": phoneNumber})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getAllUsersAdmin(c *gin.Context) {
	// Get all users
	rows, err := db.Query("SELECT id, username, real_name, email, address, phone_number FROM users")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var users []gin.H
	for rows.Next() {
		var id, username, realName, email, address, phoneNumber string
		err = rows.Scan(&id, &username, &realName, &email, &address, &phoneNumber)
		if err != nil {
			c.JSON(500, gin.H{"error": serverErrorMessage})
			return
		}
		users = append(users, gin.H{"id": id, "username": username, "real_name": realName, "email": email, "address": address, "phone_number": phoneNumber})
	}
	c.JSON(200, users)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func deleteUserAdmin(c *gin.Context) {
	// Delete user
	stmt, err := db.Prepare("DELETE FROM users WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func grantAdmin(c *gin.Context) {
	// Grant admin
	stmt, err := db.Prepare("UPDATE users SET user_type = 'admin' WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "Admin granted"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func revokeAdmin(c *gin.Context) {
	// Revoke admin
	stmt, err := db.Prepare("UPDATE users SET user_type = 'user' WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "Admin revoked"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getUserPermissionsAdmin(c *gin.Context) {
	// Get user permissions
	rows, err := db.Query("SELECT id, username, user_type FROM users")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var users []gin.H
	for rows.Next() {
		var id, username, userType string
		err = rows.Scan(&id, &username, &userType)
		if err != nil {
			c.JSON(500, gin.H{"error": serverErrorMessage})
			return
		}
		users = append(users, gin.H{"id": id, "username": username, "user_type": userType})
	}
	c.JSON(200, users)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getUserPurchaseHistoryAdmin(c *gin.Context) {
	// Get purchase history
	stmt, err := db.Prepare("SELECT product_id, request_id, amount, created_at, payment_status FROM purchase_history WHERE user_id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	rows, err := stmt.Query(c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	var purchaseHistory []gin.H
	for rows.Next() {
		var productID, requestID, paymentStatus string
		var amount int
		var createdAt string
		err = rows.Scan(&productID, &requestID, &amount, &createdAt, &paymentStatus)
		if err != nil {
			c.JSON(500, gin.H{"error": serverErrorMessage})
			return
		}
		purchaseHistory = append(purchaseHistory, gin.H{"product_id": productID, "request_id": requestID, "amount": amount, "created_at": createdAt, "payment_status": paymentStatus})
	}
	c.JSON(200, purchaseHistory)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getPurchaseHistoryAdmin(c *gin.Context) {
	// Get purchase history
	rows, err := db.Query("SELECT product_id, request_id, amount, created_at, request_message, address, payment_status FROM purchase_history")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)
	var purchaseHistory []gin.H
	for rows.Next() {
		var productID, requestID, paymentStatus string
		var amount int
		var createdAt, requestMessage, address string
		err = rows.Scan(&productID, &requestID, &amount, &createdAt, &requestMessage, &address, &paymentStatus)
		if err != nil {
			c.JSON(500, gin.H{"error": serverErrorMessage})
			return
		}
		purchaseHistory = append(purchaseHistory, gin.H{"product_id": productID, "request_id": requestID, "amount": amount, "created_at": createdAt, "request_message": requestMessage, "address": address, "payment_status": paymentStatus})
	}
	c.JSON(200, purchaseHistory)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func editPurchaseAdmin(c *gin.Context) {
	// Get request body
	var purchaseReq struct {
		Amount         int    `json:"amount"`
		RequestMessage string `json:"request_message"`
		Address        string `json:"address"`
		PaymentStatus  string `json:"payment_status"`
	}
	if err := c.ShouldBindJSON(&purchaseReq); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}

	// Edit purchase history
	stmt, err := db.Prepare("UPDATE purchase_history SET amount = ?, request_message = ?, address = ?, payment_status = ? WHERE request_id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(purchaseReq.Amount, purchaseReq.RequestMessage, purchaseReq.Address, purchaseReq.PaymentStatus, c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "Purchase updated"})
}

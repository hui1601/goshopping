package routers

import (
	"database/sql"
	"encoding/json"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type purchaseRequest struct {
	ProductID      string `json:"product_id"`
	RequestID      string `json:"request_id"`
	Amount         int    `json:"amount"` // Amount to pay
	RequestMessage string `json:"request_message"`
	Address        string `json:"address"`
}

func RegisterPurchaseRoutes(router *gin.Engine) {
	authorized := router.Group("/purchase", AuthRequired)
	authorized.POST("/confirm", confirmPurchase)
	authorized.POST("/request", requestPurchase)
	authorized.GET("", getPurchaseHistory)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func createPurchaseTables() {
	// Create purchase_history table
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS purchase_history (product_id VARCHAR(12), request_id VARCHAR(36) PRIMARY KEY, amount INT, request_message TEXT, address TEXT, user_id VARCHAR(12) REFERENCES users(id), created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, payment_status VARCHAR(255) DEFAULT 'pending')")
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
}

// requestPurchase: Request a purchase from the user(product_id, request_id, amount)
//
//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func requestPurchase(c *gin.Context) {
	// Get request body
	var purchaseReq purchaseRequest
	if err := c.ShouldBindJSON(&purchaseReq); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	// Request product info
	stmt, err := db.Prepare("SELECT price FROM products WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var price int
	err = stmt.QueryRow(purchaseReq.ProductID).Scan(&price)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	// Check if amount is same as price
	if purchaseReq.Amount != price {
		c.JSON(400, gin.H{"error": "Invalid amount"})
		return
	}
	// Set purchase request to session
	session := sessions.Default(c)
	session.Set("purchase_request", purchaseReq)
	err = session.Save()
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "Purchase requested"})
}

// confirmPurchase: Confirm a purchase from the user(paymentKey, orderId, amount)
//
//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func confirmPurchase(c *gin.Context) {
	// Get request body
	var paymentRequest struct {
		OrderID    string `json:"orderId"`
		Amount     int    `json:"amount"`
		PaymentKey string `json:"paymentKey"`
	}
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	// compare with purchase_requests
	session := sessions.Default(c)
	purchaseReq := session.Get("purchase_request").(purchaseRequest)
	// Check if amount is same as requested
	if purchaseReq.Amount != paymentRequest.Amount {
		c.JSON(400, gin.H{
			"error":            "Invalid amount",
			"requested_amount": purchaseReq.Amount,
			"actual_amount":    paymentRequest.Amount,
		})
		return
	}
	secretKey := os.Getenv("TOSS_PAYMENTS_SECRET_KEY")
	client := &http.Client{}
	jsonBody, err := json.Marshal(paymentRequest)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	req, err := http.NewRequest("POST", "https://api.tosspayments.com/v1/payments/confirm", nil)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	encryptedSecretKey := secretKey + ":"
	encryptedSecretKey = encodeBase64([]byte(encryptedSecretKey))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "Basic "+encryptedSecretKey)
	jsonBodyStr := string(jsonBody)
	req.Body = io.NopCloser(strings.NewReader(jsonBodyStr))
	res, err := client.Do(req)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	if res.StatusCode != 200 {
		// read response body as string
		body, _ := io.ReadAll(res.Body)
		c.JSON(500, gin.H{
			"error":     serverErrorMessage,
			"api_error": body,
		})
		return
	}
	// Insert purchase history
	stmt, err := db.Prepare("INSERT INTO purchase_history (user_id, product_id, request_id, amount, request_message, address) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	user := session.Get("user").(string)
	_, err = stmt.Exec(user, purchaseReq.ProductID, purchaseReq.RequestID, purchaseReq.Amount, purchaseReq.RequestMessage, purchaseReq.Address)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "Purchase confirmed"})
}

// getPurchaseHistory: Get purchase history of the user
//
//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getPurchaseHistory(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user").(string)
	// Get purchase history
	stmt, err := db.Prepare("SELECT product_id, request_id, amount, created_at, payment_status, request_message, address FROM purchase_history WHERE user_id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	rows, err := stmt.Query(user)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		log.Println(err.Error())
		return
	}
	var purchaseHistory = make([]map[string]interface{}, 0)
	for rows.Next() {
		var productID, requestID, paymentStatus, createdAt, requestMessage, address string
		var amount int
		err = rows.Scan(&productID, &requestID, &amount, &createdAt, &paymentStatus, &requestMessage, &address)
		if err != nil {
			c.JSON(500, gin.H{"error": serverErrorMessage})
			return
		}
		purchaseHistory = append(purchaseHistory, map[string]interface{}{
			"product_id":      productID,
			"request_id":      requestID,
			"amount":          amount,
			"created_at":      createdAt,
			"payment_status":  paymentStatus,
			"request_message": requestMessage,
			"address":         address,
		})
	}
	c.JSON(200, purchaseHistory)
}

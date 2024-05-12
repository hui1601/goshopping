package routers

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"log"
)

func RegisterProductsRoutes(router *gin.Engine) {
	authorized := router.Group("/products", AuthRequired)
	authorized.GET("", getProducts)
	authorized.GET("/:id", getProductDetail)

	admin := authorized.Use(AdminRequired)
	admin.POST("", addProduct)
	admin.PUT("/:id", editProductDetail)
	admin.DELETE("/:id", deleteProduct)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func checkProductsTable() bool {
	// Check if products table exists
	stmt, err := db.Prepare("SHOW TABLES LIKE 'products'")
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
func createProductsTable() {
	if checkProductsTable() {
		return
	}
	// Create products table
	stmt, err := db.Prepare("CREATE TABLE IF NOT EXISTS products (id VARCHAR(12) PRIMARY KEY, name VARCHAR(255), image VARCHAR(12) REFERENCES images(id), price INT)")
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

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func checkProductID(id string) bool {
	// Check if product ID exists
	stmt, err := db.Prepare("SELECT id FROM products WHERE id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var productID string
	err = stmt.QueryRow(id).Scan(&productID)
	return err == nil
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func addProduct(c *gin.Context) {
	// Get request body
	var addProductRequest struct {
		Name  string `json:"name"`
		Image string `json:"image"`
		Price int    `json:"price"`
	}
	if err := c.ShouldBindJSON(&addProductRequest); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	// Add product
	stmt, err := db.Prepare("INSERT INTO products (id, name, image, price) VALUES (?, ?, ?, ?)")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var id = createRandomID()
	for checkProductID(id) {
		id = createRandomID()
	}
	_, err = stmt.Exec(id, addProductRequest.Name, addProductRequest.Image, addProductRequest.Price)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "Product added"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getProducts(c *gin.Context) {
	// Get all products
	rows, err := db.Query("SELECT id, name, image, price FROM products")
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	var products = make([]map[string]interface{}, 0)
	for rows.Next() {
		var id, name, image string
		var price int
		err = rows.Scan(&id, &name, &image, &price)
		if err != nil {
			c.JSON(500, gin.H{"error": serverErrorMessage})
			return
		}
		products = append(products, map[string]interface{}{
			"id":    id,
			"name":  name,
			"image": image,
			"price": price,
		})
	}
	c.JSON(200, products)
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func getProductDetail(c *gin.Context) {
	// Get product detail
	stmt, err := db.Prepare("SELECT name, image, price FROM products WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	var name, image string
	var price int
	err = stmt.QueryRow(c.Param("id")).Scan(&name, &image, &price)
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, map[string]interface{}{
		"name":  name,
		"image": image,
		"price": price,
	})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func editProductDetail(c *gin.Context) {
	// Get request body
	var editProductRequest struct {
		Name  string `json:"name"`
		Image string `json:"image"`
		Price int    `json:"price"`
	}
	if err := c.ShouldBindJSON(&editProductRequest); err != nil {
		c.JSON(400, gin.H{
			"message": "Invalid request",
		})
		return
	}
	// Edit product detail
	stmt, err := db.Prepare("UPDATE products SET name = ?, image = ?, price = ? WHERE id = ?")
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)
	_, err = stmt.Exec(editProductRequest.Name, editProductRequest.Image, editProductRequest.Price, c.Param("id"))
	if err != nil {
		c.JSON(500, gin.H{"error": serverErrorMessage})
		return
	}
	c.JSON(200, gin.H{"message": "Product detail updated"})
}

//goland:noinspection SqlDialectInspection,SqlNoDataSourceInspection
func deleteProduct(c *gin.Context) {
	// Delete product
	stmt, err := db.Prepare("DELETE FROM products WHERE id = ?")
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
	c.JSON(200, gin.H{"message": "Product deleted"})
}

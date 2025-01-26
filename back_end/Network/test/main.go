package test

//
//import (
//	"fmt"
//	"net/http"
//
//	"github.com/gin-gonic/gin"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//)
//
//type User struct {
//	ID       uint   `json:"id" gorm:"primaryKey"`
//	Username string `json:"username" gorm:"unique;not null"`
//	Password string `json:"password" gorm:"not null"`
//}
//
//var db *gorm.DB
//
//func main() {
//	dsn := "root:root@tcp(127.0.0.1:3307)/testdb?charset=utf8mb4&parseTime=True&loc=Local"
//	var err error
//	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
//	if err != nil {
//		panic("failed to connect to the database")
//	}
//
//	if err := db.AutoMigrate(&User{}); err != nil {
//		panic("failed to migrate database")
//	}
//
//	insertTestData()
//
//	r := gin.Default()
//
//	// 设置 CORS 中间件
//	r.Use(corsMiddleware())
//
//	// 路由
//	r.POST("/register", registerHandler)
//	r.POST("/login", loginHandler)
//
//	// 运行服务
//	if err := r.Run("192.168.31.36:8080"); err != nil {
//		panic("failed to start server")
//	}
//}
//
//// CORS 中间件
//func corsMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// 设置跨域头
//		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173") // 限定前端来源
//		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
//		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
//		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true") // 支持跨域 Cookies
//
//		// 如果是 OPTIONS 请求，提前返回 204 状态码
//		if c.Request.Method == http.MethodOptions {
//			c.AbortWithStatus(http.StatusNoContent)
//			return
//		}
//
//		c.Next()
//	}
//}
//
//func insertTestData() {
//	var count int64
//	db.Model(&User{}).Count(&count)
//	if count == 0 {
//		users := []User{
//			{Username: "user1", Password: "password1"},
//			{Username: "user2", Password: "password2"},
//		}
//		if err := db.Create(&users).Error; err != nil {
//			fmt.Println("Failed to insert test data:", err)
//		}
//	}
//}
//
//func registerHandler(c *gin.Context) {
//	var user User
//	if err := c.ShouldBindJSON(&user); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//		return
//	}
//
//	if err := db.Create(&user).Error; err != nil {
//		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Registration successful"})
//}
//
//func loginHandler(c *gin.Context) {
//	var loginData User
//	if err := c.ShouldBindJSON(&loginData); err != nil {
//		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
//		return
//	}
//
//	var user User
//	if err := db.Where("username = ?", loginData.Username).First(&user).Error; err != nil || user.Password != loginData.Password {
//		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
//		return
//	}
//
//	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
//}

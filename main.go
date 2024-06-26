package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var db *sql.DB

func main() {
	var err error

	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Get the database URL from the environment variable
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	// Connect to the database
	db, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping the database to ensure a connection is established
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}

	// Set up the Gin router
	router := gin.Default()

	// Define a simple GET endpoint
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/users", GetEmployee)

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	router.Run(":" + port)
}

type Employee struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	Department  string `json:"department"`
}

func GetEmployee(c *gin.Context) {
	var employees []Employee
	name := c.Query("name")
	var rows *sql.Rows
	var err error

	query := "SELECT id, name, email, address, phoneNumber, department FROM ms_employee"
	if name == "" {
		rows, err = db.Query(query)
	} else {
		query += " WHERE name = $1;"
		rows, err = db.Query(query, name)
	}

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	defer rows.Close()
	for rows.Next() {
		var emp Employee
		err = rows.Scan(&emp.Id, &emp.Name, &emp.Email, &emp.Address, &emp.PhoneNumber, &emp.Department)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		employees = append(employees, emp)
	}
	if err = rows.Err(); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(employees) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "employee.name tidak ditemukan"})
		return
	} else {
		c.JSON(200, gin.H{
			"message": "data employee",
			"data":    employees,
		})
	}
}

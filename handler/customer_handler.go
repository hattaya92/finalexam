package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hattaya92/finalexam/middleware"

	"github.com/gin-gonic/gin"
	"github.com/hattaya92/finalexam/database"

	_ "github.com/lib/pq"
)

type Customers struct {
	ID     int    `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Status string `json:"status"`
}

//var db *sql.DB

func CreateTable() {
	createDB := "CREATE TABLE IF NOT EXISTS customers (id SERIAL PRIMARY KEY,name TEXT,email TEXT,status TEXT)"
	_, err := database.ConnDB().Exec(createDB)
	if err != nil {
		log.Fatal("create database error:", err)
		return
	}
	fmt.Println("created database success")

}

func createCustomersHandler(c *gin.Context) {
	var item Customers
	err := c.ShouldBindJSON(&item)
	if err != nil {
		log.Fatal("input error: ", err)
		c.JSON(http.StatusBadRequest, gin.H{"status: ": err.Error()})
		return
	}
	row := database.ConnDB().QueryRow("INSERT INTO customers (name,email,status) VALUES ($1,$2,$3) RETURNING id", item.Name, item.Email, item.Status)
	var id int
	err = row.Scan(&id)
	if err != nil {
		log.Fatal("scan error: ", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status:": err.Error()})
		return
	}
	item.ID = id
	c.JSON(201, item)

}

func getCustomersByIDHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := database.ConnDB().Prepare("SELECT * FROM customers WHERE id=$1")
	if err != nil {
		log.Fatal("prepare statment error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status :": err.Error()})
		return
	}
	cust := Customers{}
	row := stmt.QueryRow(id)
	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		log.Fatal("query error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status :": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cust)
}

func getAllCustomersHandler(c *gin.Context) {
	stmt, err := database.ConnDB().Prepare("SELECT * FROM customers")
	if err != nil {
		log.Fatal("prepare statment error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status :": err.Error()})
		return
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("cannnot query all customers", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status :": err.Error()})
		return
	}

	var allCust []Customers

	for rows.Next() {
		item := Customers{}
		err := rows.Scan(&item.ID, &item.Name, &item.Email, &item.Status)
		if err != nil {
			log.Fatal("scan error :", err)
			c.JSON(http.StatusInternalServerError, gin.H{"status :": err.Error()})
			return
		}
		allCust = append(allCust, item)
	}

	c.JSON(http.StatusOK, allCust)

}

func updateCustomerByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var item Customers
	err := c.ShouldBindJSON(&item)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": err.Error()})
		return
	}

	stmt, err := database.ConnDB().Prepare("UPDATE customers SET name=$2,email=$3,status=$4 WHERE id=$1")
	if err != nil {
		log.Fatal("prepare statment error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"ERROR :": err.Error()})
		return
	}

	if _, err := stmt.Exec(id, item.Name, item.Email, item.Status); err != nil {
		log.Fatal("error execute update ", err)
	}

	stmt, err = database.ConnDB().Prepare("SELECT * FROM customers WHERE id=$1")
	if err != nil {
		log.Fatal("prepare statment error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status :": err.Error()})
		return
	}
	cust := Customers{}

	row := stmt.QueryRow(id)
	err = row.Scan(&cust.ID, &cust.Name, &cust.Email, &cust.Status)
	if err != nil {
		log.Fatal("scan error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status :": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cust)

}

func delCustomerByIDHandler(c *gin.Context) {
	id := c.Param("id")
	stmt, err := database.ConnDB().Prepare("DELETE FROM customers WHERE id=$1")
	if err != nil {
		log.Fatal("prepare statment error :", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status:": err.Error()})
		return
	}
	if _, err := stmt.Exec(id); err != nil {
		log.Fatal("delete error", err)
		c.JSON(http.StatusInternalServerError, gin.H{"status :": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "customer deleted"})

}

func loginMiddleware(c *gin.Context) {
	authKey := c.GetHeader("Authorization")
	if authKey != "token2019" {
		c.JSON(http.StatusUnauthorized, "Status code is 401 Unauthorized")
		c.Abort()
		return
	}

	c.Next()
	log.Println("ending middleware")
}

func NewRouter() *gin.Engine {
	r := gin.Default()
	r.Use(middleware.LoginMiddleware)
	r.POST("/customers", createCustomersHandler)
	r.GET("/customers/:id", getCustomersByIDHandler)
	r.GET("/customers/", getAllCustomersHandler)
	r.PUT("/customers/:id", updateCustomerByIDHandler)
	r.DELETE("/customers/:id", delCustomerByIDHandler)
	return r
}

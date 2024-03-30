package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type Fibonacci struct {
	Number int `json:"number"`
	Result int `json:"result"`
}

func main() {

	db, err := sql.Open("postgres", "user=pguser password=pgpasswd host=pgservice.fibnamespace.svc.cluster.local port=5432 dbname=fibonnaci sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	router := gin.Default()

	_, err = db.Exec(`
	CREATE TABLE IF NOT EXISTS fibonacci (
		id SERIAL PRIMARY KEY,
		number INT,
		result INT
	)
	`)
	router.LoadHTMLGlob("index.html")
	router.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", gin.H{
			"appName":   "Fibonacci App",
			"endpoint1": "/fibonacci/:n",
			"endpoint2": "/fibonacci/",
		})
	})

	if err != nil {
		log.Fatal(err)
	}
	router.GET("/fibonacci/:n", func(c *gin.Context) {
		n, err := strconv.Atoi(c.Param("n"))

		if err != nil {
			c.JSON(200, gin.H{"error": "Invalid input"})
			return
		}
		fibonacci := fibonacciSequence(n)
		c.JSON(200, gin.H{"fibonacci": fibonacci})
		last_element := fibonacci[len(fibonacci)-1]

		fmt.Println("the last element is : ", last_element)
		_, err = db.Exec("INSERT INTO fibonacci (number, result) VALUES ($1, $2)", n, last_element)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

	})

	router.GET("/fibonacci/", func(c *gin.Context) {
		var fibonacci []Fibonacci
		rows, err := db.Query("SELECT number, result FROM fibonacci")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		defer rows.Close()
		for rows.Next() {
			var f Fibonacci
			err := rows.Scan(&f.Number, &f.Result)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
				return
			}
			fibonacci = append(fibonacci, f)
		}
		c.JSON(http.StatusOK, gin.H{"fibonacci": fibonacci})
	})

	router.Run(":8000")
}

func fibonacciSequence(n int) []int {
	if n <= 1 {
		return []int{1}
	}
	fibonacci := make([]int, n)
	fibonacci[0] = 1
	fibonacci[1] = 1
	for i := 2; i < n; i++ {
		fibonacci[i] = fibonacci[i-1] + fibonacci[i-2]
	}
	return fibonacci
}

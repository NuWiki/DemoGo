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

type Todo struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Status string `json:"status"`
}

func helloHandler(c *gin.Context) {
	log.Println("ln hellohandler")
	c.JSON(http.StatusOK, gin.H{
		"message": "hello",
	})
}

func todoPostHandler(c *gin.Context) {
	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Printf("dddd%v", todo)
	todo.insert()
	c.JSON(http.StatusOK, gin.H{
		"message": todo,
	})
}

func (todo *Todo) insert() {
	postgresDb, err := sql.Open("postgres", "postgres://mydglzdw:3dtcdLkiIO93u2CQDE8SC2hIFP80WDjh@suleiman.db.elephantsql.com:5432/mydglzdw")
	if err != nil {
		log.Fatal("Connect to database error", err)

	}
	defer postgresDb.Close()

	row := postgresDb.QueryRow("INSERT INTO todos (title, status) values ($1,$2) RETURNING id", todo.Title, todo.Status)
	var id int

	err = row.Scan(&id)
	if err != nil {
		fmt.Println("cant scan id", err)
		return
	}
	todo.ID = id
	fmt.Println("insert todo success : ", id)
	// return Todo{id, todo.Title, todo.Status}

}

func todoGetHandler(c *gin.Context) {
	var todos = get()
	c.JSON(http.StatusOK, gin.H{
		"message": todos,
	})
}

func get() []Todo {
	db, err := sql.Open("postgres", "postgres://mydglzdw:3dtcdLkiIO93u2CQDE8SC2hIFP80WDjh@suleiman.db.elephantsql.com:5432/mydglzdw")
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, status FROM todos")
	if err != nil {
		log.Fatal("cant prepare one all todos statement", err)
	}

	rows, err := stmt.Query()
	if err != nil {
		log.Fatal("cant prepare one all todos statement", err)
	}
	var s = []Todo{}
	// for rows.Next() {
	// 	var id int
	// 	var title, status string
	// 	err := rows.Scan(&id, &title, &status)

	// 	s = append(s, Todo{id, title, status})

	// 	if err != nil {
	// 		log.Fatal("cant scan row into variable", err)
	// 	}
	// 	fmt.Println(id, title, status)
	// }
	for rows.Next() {

		t := Todo{}
		err := rows.Scan(&t.ID, &t.Title, &t.Status)

		s = append(s, t)

		if err != nil {
			log.Fatal("cant scan row into variable", err)
		}

	}
	fmt.Printf("%v", s)
	fmt.Println("query all todos success")
	return s
}

func todoGetByIdHandler(c *gin.Context) {
	userid := c.Param("userid")
	i, err := strconv.Atoi(userid)
	if err != nil {
		log.Fatal("cast fail", err)
	}
	var todos = gett(i)
	c.JSON(http.StatusOK, gin.H{
		"messadddge": todos,
	})
}

func gett(id int) Todo {
	db, err := sql.Open("postgres", "postgres://mydglzdw:3dtcdLkiIO93u2CQDE8SC2hIFP80WDjh@suleiman.db.elephantsql.com:5432/mydglzdw")
	if err != nil {
		log.Fatal("Connect to database error", err)
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT id, title, status FROM todos where id=$1")
	if err != nil {
		log.Fatal("cant prepare one row statement", err)
	}

	row := stmt.QueryRow(id)

	todo := Todo{}
	err = row.Scan(&todo.ID, &todo.Title, &todo.Status)
	if err != nil {
		log.Fatal("cant scan row into variables", err)
	}

	return todo
}

func main() {
	r := gin.Default()
	r.Use(middleWare)
	r.Use(middleWare2)
	r.POST("/todos", todoPostHandler)
	r.GET("/todos", todoGetHandler)
	r.GET("/todos/:id", todoGetByIdHandler)
	// port := os.Getenv("PORT")
	// r.Run(port)
	r.Run(":1234")
}

func middleWare(c *gin.Context) {
	log.Println("start middleware")
	c.Next()
}
func middleWare2(c *gin.Context) {
	log.Println("2 start middleware")
	c.Next()
	log.Println("2 end middleware")
}

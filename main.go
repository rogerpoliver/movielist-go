package main

import (
	"database/sql"
	"net/http"
	"text/template"

	_ "github.com/lib/pq"
)

type Category struct {
	Id   int
	Name string
}

type StreamingService struct {
	Id   int
	Name string
}

type Show struct {
	Id       int
	Year     int
	Name     string
	Category Category
	Service  StreamingService
}

func getDbConnection() *sql.DB {
	connectionString := "user=postgres dbname=movielist password=postgres host=localhost sslmode=disable"
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var templates = template.Must(template.ParseGlob("templates/*.html"))

func main() {
	db := getDbConnection()
	defer db.Close()
	http.HandleFunc("/", index)
	http.ListenAndServe(":8000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	services := []StreamingService{}
	db := getDbConnection()

	data, err := db.Query("select * from streaming_services")
	if err != nil {
		panic(err.Error())
	}

	for data.Next() {
		var id int
		var name string

		err = data.Scan(&id, &name)

		service := StreamingService{
			Id:   id,
			Name: name,
		}

		services = append(services, service)
	}

	templates.ExecuteTemplate(w, "Index", services)
	defer db.Close()
}

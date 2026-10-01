package main

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/joho/godotenv"
)

// func main() {
// 	connString := "postgres://sevia:sevia%402024@localhost:5432/school"
// 	conn, err := pgx.Connect(context.Background(), connString)
// 	if err != nil {
// 		log.Fatal("Error connecting to database:", err)
// 		return
// 	}
// 	defer conn.Close(context.Background())
// 	fmt.Println("Connected to PostgreSQL database successfully!")
// 	err = conn.Ping(context.Background())
// 	if err != nil {
// 		log.Fatal("Error pinging database:", err)
// 	}
// 	fmt.Println("Database is working fine!")
// }

func main() {
	//load  .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	//get database credential
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_NAME")

	//encode password for url
	password = url.QueryEscape(password)

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s",
		user,
		password,
		host,
		port,
		database,
	)
	fmt.Println("Connecting to PostgreSQL.......")

	//connect to database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
		return
	}
	defer conn.Close(context.Background())

	fmt.Println("Connected to PostgreSQL database successfully!")

	//test database connection
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}
	fmt.Println("Database is working fine!")

	//query to create table
	rows, err := conn.Query(context.Background(),
		"SELECT id,name,age,email,created_at FROM students")
	if err != nil {
		log.Fatal("Error querying database:", err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int
		var name string
		var age int
		var email string
		var createdAt time.Time

		err := rows.Scan(&id, &name, &age, &email, &createdAt)
		if err != nil {
			log.Fatal("Error scanning row:", err)
		}
		////convert timestamp to Nairobi time zone
		location, err := time.LoadLocation("Africa/Nairobi")
		if err != nil {
			log.Fatal("Error loading location:", err)
		}
		createdAt = createdAt.In(location)
		fmt.Println("-----------------------------")
		fmt.Println("ID:", id)
		fmt.Println("Name:", name)
		fmt.Println("Age:", age)
		fmt.Println("Email:", email)
		fmt.Println("Created:", createdAt.Format("2006-01-02 15:04:05"))
	}
	if err := rows.Err(); err != nil {
		log.Fatal("Error iterating over rows:", err)
	}
}

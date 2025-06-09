package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	"github.com/e6a5/learning/backend/02-mysql-crud/internal/handlers"
	"github.com/e6a5/learning/backend/02-mysql-crud/internal/repository"
)

func main() {
	// Initialize database connection
	db, err := initializeDatabase()
	if err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer db.Close()

	// Initialize dependencies
	userRepo := repository.NewUserRepository(db)
	userHandler := handlers.NewUserHandler(userRepo)

	// Setup HTTP server
	router := setupRoutes(userHandler)

	log.Println("🛠️  Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

func initializeDatabase() (*sql.DB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		dsn = "user:password@tcp(localhost:3306)/testdb"
	}

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func setupRoutes(userHandler *handlers.UserHandler) *mux.Router {
	router := mux.NewRouter()

	// User CRUD routes
	router.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	router.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")

	return router
}

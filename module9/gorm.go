package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Blog struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"unique"`
	Content   string
	Timestamp time.Time `gorm:"autoCreateTime"`
}

func SetUpDBConnection() *gorm.DB {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file: ", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("You must provide a database name as an argument")
	}

	dbName := os.Args[1]
	dbUser := os.Getenv("User")
	dbPassword := os.Getenv("Password")

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8mb4&parseTime=True&loc=Local", dbUser, dbPassword, dbName)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to database: ", err)
	}

	return db
}

func main() {
	db := SetUpDBConnection()
	if db != nil {
		fmt.Println("Successfully connected to the database")
	}

	r := mux.NewRouter()
	r.HandleFunc("/blog", CreateBlog(db)).Methods("POST")
	r.HandleFunc("/blogs", GetAllBlogs(db)).Methods("GET")
	r.HandleFunc("/blog/{id}", GetBlog(db)).Methods("GET")
	r.HandleFunc("/blog/{id}", UpdateBlog(db)).Methods("PUT")
	r.HandleFunc("/blog/{id}", DeleteBlog(db)).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":8080", r))
}

func CreateBlog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blog Blog
		if err := json.NewDecoder(r.Body).Decode(&blog); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if err := db.Create(&blog).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(blog)
	}
}

func DeleteBlog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		if err := db.Where("id = ?", id).Delete(&Blog{}).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func UpdateBlog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var blog Blog
		if err := db.First(&blog, id).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		var updateData Blog
		if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		blog.Content = updateData.Content
		blog.Title = updateData.Title

		if err := db.Save(&blog).Error; err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(blog)
	}
}

func GetAllBlogs(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var blogs []Blog

		if err := db.Find(&blogs).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(blogs)
	}
}

func GetBlog(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]

		var blog Blog
		if err := db.First(&blog, id).Error; err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(blog)
	}
}

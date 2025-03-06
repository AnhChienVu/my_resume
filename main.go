package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

// Database Connection Details
const (
	DB_USER     = "postgres"
	DB_PASSWORD = "Anhchien12345@"
	DB_NAME     = "my_resume"
	DB_HOST     = "localhost"
	DB_PORT     = "5432"
)

// ResumeData represents the data to be displayed on the resume page
type ResumeData struct {
	ObjectiveCount int
	EducationCount int
	SkillsCount    int
	ProjectCount   int
}

// Database connection details
var db *sql.DB

// Initialize database connection
func initDB() {
	var err error
	connStr := fmt.Sprintf("host= %s port= %s user= %s password= %s dbname= %s sslmode=disable", DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
}

// Get click counts for all sections
func getClickCounts() (ResumeData, error) {
	counts := ResumeData{}

	query := "SELECT section, count FROM clicks"
	rows, err := db.Query(query)
	if err != nil {
		return counts, err
	}
	defer rows.Close()

	for rows.Next() {
		var section string
		var count int
		if err := rows.Scan(&section, &count); err != nil {
			return counts, err
		}

		switch section {
		case "objective":
			counts.ObjectiveCount = count
		case "education":
			counts.EducationCount = count
		case "skills":
			counts.SkillsCount = count
		case "projects":
			counts.ProjectCount = count
		}
	}
	return counts, nil
}

func incrementClickCount(w http.ResponseWriter, r *http.Request) {
	section := r.URL.Query().Get("section")
	if section == "" {
		http.Error(w, "Section parameter is required", http.StatusBadRequest)
		return
	}

	_, err := db.Exec("UPDATE clicks SET count = count + 1 WHERE section = $1", section)
	if err != nil {
		http.Error(w, "Error updating click count", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

func main() {
	initDB()
	defer db.Close()

	// Serve static files (JavaScript, CSS)
	http.Handle("/js/", http.StripPrefix("/js/", http.FileServer(http.Dir("js"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/increment", incrementClickCount)

	fmt.Println("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Server is not running:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	counts, err := getClickCounts()
	if err != nil {
		http.Error(w, "Error getting click counts", http.StatusInternalServerError)
		return
	}

	tpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Error loading template", http.StatusInternalServerError)
		return
	}
	tpl.Execute(w, counts)
}

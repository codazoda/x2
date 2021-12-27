package main

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	// Set some port variables
	const port = "8002"
	const address = ":" + port
	// Setup a file server to serve the www directory
	fileServer := http.FileServer(http.Dir("./www"))
	http.Handle("/", fileServer)
	// Setup functions to handle other URI's
	http.HandleFunc("/root", rootHandler)
	http.HandleFunc("/upload", uploadHandler)
	http.HandleFunc("/image/", imageHandler)
	http.HandleFunc("/create", createHandler)
	// Start a server
	fmt.Printf("Server started on port %s\n", port)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, "Hello World.\n")
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Reject if the method wasn't post
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Restrict the uploads to a maximum size
	megabyte := 1048576 // A megabyte is 1024 * 1024 (1048576) bytes
	maxMegs := 10       // 1048576 bytes is 1024 * 1024 or 1MB
	maxUpload := int64(maxMegs * megabyte)
	r.Body = http.MaxBytesReader(w, r.Body, maxUpload)
	if err := r.ParseMultipartForm(maxUpload); err != nil {
		http.Error(w, "The uploaded file is too big. Please choose an file that's less than "+strconv.Itoa(maxMegs)+"MB in size", http.StatusBadRequest)
		return
	}
	// Grab the file
	file, _, err := r.FormFile("image")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	// Close the file in the future
	defer file.Close()
	// Create a big buffer to read the file into
	fileBuffer := make([]byte, maxUpload)
	// Base64 encode the file into the big buffer
	readTotal, err := file.Read(fileBuffer)
	if err != nil {
		if err != io.EOF {
			fmt.Println("Error reading the file")
			fmt.Println(err)
		}
	}
	// TODO: Encrypt the file
	// Write the file to the DB
	database, _ :=
		sql.Open("sqlite3", "./images.db")
	statement, _ :=
		database.Prepare("INSERT INTO images (image) VALUES (?)")
	statement.Exec(base64.StdEncoding.EncodeToString(fileBuffer[:readTotal]))
	// Say something
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, "File saved")
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the vars from the request (including the URI)
	idString := r.URL.Path[len("/image/"):]
	// Query for the image from the file
	database, _ :=
		sql.Open("sqlite3", "./images.db")
	id, _ := strconv.Atoi(idString)
	row := database.QueryRow("SELECT image FROM images WHERE id = ?", id)
	var image string
	row.Scan(&image)
	// TODO: Decrypt the file
	// Show the image
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, "<img src='data:image/jpeg;base64,"+image+"' style='width: 100%'>")
}

func createHandler(w http.ResponseWriter, r *http.Request) {
	// Create a database an add a record
	database, _ :=
		sql.Open("sqlite3", "./images.db")
	statement, _ :=
		database.Prepare("CREATE TABLE IF NOT EXISTS images (id INTEGER PRIMARY KEY, image BLOB)")
	statement.Exec()
}

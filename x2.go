package main

import (
	"database/sql"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

type application struct {
	database struct {
		file string
	}
}

func main() {
	app := new(application)
	app.database.file = "./images.db"

	// Initialize the application
	app.init()
	// Set some port variables
	const port = "8002"
	const address = ":" + port
	// Setup a file server to serve the www directory
	fileServer := http.FileServer(http.Dir("./www"))
	http.Handle("/", fileServer)
	// Setup functions to handle other URI's
	http.HandleFunc("/root", app.rootHandler)
	http.HandleFunc("/upload", app.uploadHandler)
	http.HandleFunc("/album", app.albumHandler)
	http.HandleFunc("/image/", app.imageHandler)
	// Start a server
	fmt.Printf("Server started on port %s\n", port)
	if err := http.ListenAndServe(address, nil); err != nil {
		log.Fatal(err)
	}
}

func (app *application) init() {
	// Check if the database file exists and create it if it doesn't
	_, err := os.Open(app.database.file)
	if errors.Is(err, os.ErrNotExist) {
		// Create a database an add a record
		database, _ :=
			sql.Open("sqlite3", app.database.file)
		statement, _ :=
			database.Prepare("CREATE TABLE IF NOT EXISTS images (id INTEGER PRIMARY KEY, image BLOB)")
		statement.Exec()
	}
}

func (app *application) rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, "Hello World.\n")
}

func (app *application) uploadHandler(w http.ResponseWriter, r *http.Request) {
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
		sql.Open("sqlite3", app.database.file)
	statement, _ :=
		database.Prepare("INSERT INTO images (image) VALUES (?)")
	statement.Exec(base64.StdEncoding.EncodeToString(fileBuffer[:readTotal]))
	// Say something
	w.Header().Add("content-type", "text/html")
	io.WriteString(w, "File saved")
}

func (app *application) imageHandler(w http.ResponseWriter, r *http.Request) {
	// Grab the vars from the request (including the URI)
	idString := r.URL.Path[len("/image/"):]
	// Query for the image from the file
	database, _ :=
		sql.Open("sqlite3", app.database.file)
	id, _ := strconv.Atoi(idString)
	row := database.QueryRow("SELECT image FROM images WHERE id = ?", id)
	var image string
	row.Scan(&image)
	// TODO: Decrypt the file
	// Show the image
	w.Header().Add("content-type", "image/jpeg")
	imageBinary, _ := base64.StdEncoding.DecodeString(image)
	w.Write(imageBinary)
}

func (app *application) albumHandler(w http.ResponseWriter, r *http.Request) {
	// Setup variables
	var id int
	// Query for all the images in the album
	database, _ :=
		sql.Open("sqlite3", app.database.file)
	rows, _ := database.Query("SELECT id FROM images ORDER BY id DESC")
	// Loop through the images showing them
	w.Header().Add("content-type", "text/html")
	for rows.Next() {
		rows.Scan(&id)
		io.WriteString(w, "<img src='./image/"+strconv.Itoa(id)+"' style='width: calc(100% - 40px); padding: 20px;' loading='lazy'>\n")
	}
	// TODO: Decrypt the images
}

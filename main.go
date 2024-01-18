package main

import (
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	// creating the Bot instance
	//bot := structs.Bot{}

	// RUN BOT!
	//bot.Go()

	http.HandleFunc("/upload", uploadHandler)
	log.Println("Starting the server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests are allowed", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file: "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	log.Printf("Uploaded file: %+v\n", header.Filename)

	// Save the file
	dst, err := os.Create("D:/uploads/" + header.Filename)
	if err != nil {
		http.Error(w, "Error saving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error saving the file: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("File successfully uploaded"))
}

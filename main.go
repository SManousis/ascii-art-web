package main

import (
	"fmt"
	"html/template"
	"net/http"

	"AsciiArtWeb/funcs"
)

type PageData struct {
	InputText string
	Banner    string
	Result    string
	Error     string
}

func main() {
	http.HandleFunc("/", homeHandler)

	// Start the server on port 8080
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
		return
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	data := PageData{}
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method == http.MethodPost {
		// Process form submission
		err := r.ParseForm()
		if err != nil {
			data.Error = "Bad Request: Unable to parse form data"
			renderTemplate(w, data)
			return
		}

		inputText := r.FormValue("inputText")
		banner := r.FormValue("banner")

		// Validate inputText and banner
		if inputText == "" || banner == "" {
			data.Error = "Missing input text or banner selection"
			data.InputText = inputText
			data.Banner = banner
			renderTemplate(w, data)
			return
		}

		if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
			data.Error = "Invalid banner selection"
			data.InputText = inputText
			data.Banner = banner
			renderTemplate(w, data)
			return
		}

		// Generate ASCII art
		asciiArtResult, err := funcs.GenerateAsciiArt(inputText, banner)
		if err != nil {
			data.Error = err.Error()
			data.InputText = inputText
			data.Banner = banner
			renderTemplate(w, data)
			return
		}

		// Prepare data for template
		data.InputText = inputText
		data.Banner = banner
		data.Result = asciiArtResult
	}

	// Render the template
	renderTemplate(w, data)
}

func renderTemplate(w http.ResponseWriter, data PageData) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

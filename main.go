package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PageData struct {
	InputText string
	Banner    string
	Result    string
	Error     string
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/ascii-art", asciiArtHandler)

	// Start the server on port 8080
	fmt.Println("Server started at http://localhost:8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, nil)
}

func asciiArtHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	inputText := r.FormValue("inputText")
	banner := r.FormValue("banner")

	// Validate inputText and banner
	if inputText == "" || banner == "" {
		http.Error(w, "Bad Request: Missing input text or banner selection", http.StatusBadRequest)
		return
	}

	if banner != "standard" && banner != "shadow" && banner != "thinkertoy" {
		http.Error(w, "Bad Request: Invalid banner selection", http.StatusBadRequest)
		return
	}

	// Generate ASCII art
	asciiArtResult, err := generateAsciiArt(inputText, banner)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Prepare data for template
	data := PageData{
		InputText: inputText,
		Banner:    banner,
		Result:    asciiArtResult,
	}

	tmpl, err := template.ParseFiles("templates/result.html")
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	tmpl.Execute(w, data)
}

func generateAsciiArt(inputText, banner string) (string, error) {
	bannerPath := filepath.Join("banners", banner+".txt")
	file, err := os.ReadFile(bannerPath)
	if err != nil {
		return "", fmt.Errorf("banner not found")
	}
	str := string(file)
	str = strings.Replace(str, "\r\n", "\n", -1)
	lines := strings.Split(str, "\n")
	text := strings.Split(inputText, "\n")

	asciiArt, err := printAsciiArt(text, lines)
	if err != nil {
		return "", err
	}

	return asciiArt, nil
}

func printAsciiArt(sentences []string, textFile []string) (string, error) {
	var result strings.Builder
	for i, word := range sentences {
		if word == "" {
			// If the sentence is empty (i.e., there was a '\n'), print an empty line
			if i != 0 {
				result.WriteString("\n")
			}
			continue
		}
		// Loop through each of the 8 lines of ASCII art for the characters
		for h := 1; h < 9; h++ {
			// Loop through each character in the word
			for i := 0; i < len(word); i++ {
				char := word[i]
				if char < 32 || char > 126 {
					continue
				}
				for lineIndex, line := range textFile {
					if lineIndex == (int(char)-32)*9+h {
						result.WriteString(line)
					}
				}
			}

			result.WriteString("\n")
		}
	}
	return result.String(), nil
}

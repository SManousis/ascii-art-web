package funcs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func GenerateAsciiArt(inputText, banner string) (string, error) {
	bannerPath := filepath.Join("banners", banner+".txt")
	file, err := os.ReadFile(bannerPath)
	if err != nil {
		return "", fmt.Errorf("banner not found")
	}
	str := string(file)
	str = strings.Replace(str, "\r\n", "\n", -1)
	lines := strings.Split(str, "\n")
	text := strings.Split(inputText, "\n")
	asciiArt, err := ConvertToAsciiArt(text, lines)
	if err != nil {
		return "", err
	}
	return asciiArt, nil
}

func ConvertToAsciiArt(sentences []string, textFile []string) (string, error) {
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

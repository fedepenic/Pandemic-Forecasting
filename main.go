package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

func main() {
	url := "https://www.worldometers.info/coronavirus/country/south-africa/" // Replace this with the URL of the webpage you want to parse

	// Fetch the webpage
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching the webpage: %v", err)
	}
	defer resp.Body.Close()

	// Parse the HTML content
	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	// Get the content of <script type="text/javascript"> tags containing "graph-cases-daily"
	content := getScriptContents(doc, "graph-cases-daily")

	// Extract content using the specified regular expression
	extractedContent := extractContent(content, `data:\s*\[(.*?)\]`)
	extractedContentDates := extractContent(content, `categories:\s*\[(.*?)\]`)

	arr := convertStringToArray(extractedContent)
	arrDates, err := convertStringDatesToArray(extractedContentDates)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Verify that the length of arr and arrDates match
	if len(arr) != len(arrDates) {
		fmt.Println("Error: Lengths of arrays do not match.")
		return
	}

	// Create and open the CSV file
	file, err := os.Create("data.csv")
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write the CSV header
	writer.Write([]string{"Date", "Value"})

	// Write data rows to CSV
	for i := 0; i < len(arr); i++ {
		writer.Write([]string{arrDates[i], strconv.Itoa(arr[i])})
	}

	fmt.Println("CSV data successfully written to data.csv.")
}

func convertStringDatesToArray(str string) ([]string, error) {
	var arr []string
	err := json.Unmarshal([]byte(str), &arr)
	if err != nil {
		return nil, err
	}

	return arr, nil
}

func convertStringToArray(str string) []int {
	// Step 1: Parse the string into a slice of strings
	str = strings.Trim(str, "[]") // Remove square brackets from the string
	strArr := strings.Split(str, ",")

	// Step 2: Convert the slice of strings into a slice of integers
	var intArr []int
	for _, s := range strArr {
		// Replace "null" with "0" and parse the string to an integer
		if s == "null" {
			intArr = append(intArr, 0)
		} else {
			if num, err := strconv.Atoi(s); err == nil {
				intArr = append(intArr, num)
			}
		}
	}

	return intArr
}

// Recursive function to get the contents of <script type="text/javascript"> tags
// if they contain the specified string and return it as a string.
func getScriptContents(n *html.Node, targetString string) string {
	if n.Type == html.ElementNode && n.Data == "script" {
		for _, attr := range n.Attr {
			if attr.Key == "type" && attr.Val == "text/javascript" {
				content := getTextContent(n)
				if strings.Contains(content, targetString) {
					// Return the content of the script tag
					return content
				}
			}
		}
	}
	var scriptContent string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scriptContent = getScriptContents(c, targetString)
		if scriptContent != "" {
			return scriptContent
		}
	}
	return ""
}

// Helper function to get the text content of an HTML node and its descendants
func getTextContent(n *html.Node) string {
	if n.Type == html.TextNode {
		return n.Data
	}
	var textContent string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		textContent += getTextContent(c)
	}
	return textContent
}

// Function to extract content using a regular expression
func extractContent(content, regex string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(content)

	var finalContent string

	if len(match) > 1 {
		// Extracted content is in match[1]
		categoriesContent := match[1]
		finalContent = "[" + categoriesContent + "]"
	} else {
		finalContent = "Content not found in the input."
	}

	return finalContent
}

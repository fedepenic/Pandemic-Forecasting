package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
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

    re := regexp.MustCompile(`categories:\s*\[(.*?)\]`)

    match := re.FindStringSubmatch(content)

	var finalContent string

	if len(match) > 1 {
		// Extracted content is in match[1]
		categoriesContent := match[1]
		finalContent = "[" + categoriesContent + "]"
	} else {
		finalContent = "Categories not found in the input."
	}

	fmt.Println(finalContent)
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

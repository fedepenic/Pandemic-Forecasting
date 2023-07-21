package main

import (
	"fmt"
	"log"
	"net/http"
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

    // Print the contents of <script type="text/javascript"> tags containing "graph-cases-daily"
    printScriptContents(doc, "graph-cases-daily")
}

// Recursive function to print the contents of <script type="text/javascript"> tags
// if they contain the specified string.
func printScriptContents(n *html.Node, targetString string) {
    if n.Type == html.ElementNode && n.Data == "script" {
        for _, attr := range n.Attr {
            if attr.Key == "type" && attr.Val == "text/javascript" {
                content := getTextContent(n)
                if strings.Contains(content, targetString) {
                    // Print the content of the script tag
                    fmt.Println(content)
                }
            }
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        printScriptContents(c, targetString)
    }
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

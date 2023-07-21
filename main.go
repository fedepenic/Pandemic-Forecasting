package main

import (
	"fmt"
	"log"
	"net/http"

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

    // Find and print the content of divs with class "newsdate_div"
    printDivContentByClass(doc, "newsdate_div")
}

// Recursive function to find and print the content of divs by class
func printDivContentByClass(n *html.Node, targetClass string) {
    if n.Type == html.ElementNode && n.Data == "div" {
        classVal := getAttributeValue(n, "class")
        if classVal == targetClass {
            fmt.Println(getTextContent(n))
        }
    }
    for c := n.FirstChild; c != nil; c = c.NextSibling {
        printDivContentByClass(c, targetClass)
    }
}

// Helper function to get the value of an attribute for an HTML node
func getAttributeValue(n *html.Node, attrName string) string {
    for _, a := range n.Attr {
        if a.Key == attrName {
            return a.Val
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

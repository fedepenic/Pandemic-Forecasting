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
	url := "https://www.worldometers.info/coronavirus/country/south-africa/" 

	resp, err := http.Get(url)
	if err != nil {
		log.Fatalf("Error fetching the webpage: %v", err)
	}
	defer resp.Body.Close()

	doc, err := html.Parse(resp.Body)
	if err != nil {
		log.Fatalf("Error parsing HTML: %v", err)
	}

	//The data from the daily cases graph is obtained. 
	scriptContent := getScriptContent(doc, "graph-cases-daily")

	//Using regular expressions, the dates and their respective new Covid cases are obtained. 
	newCases := extractContent(scriptContent, `data:\s*\[(.*?)\]`)
	dates := extractContent(scriptContent, `categories:\s*\[(.*?)\]`)

	newCasesArray := convertStringToArray(newCases)
	datesArray, err := convertStringDatesToArray(dates)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	if len(newCasesArray) != len(datesArray) {
		fmt.Println("Error: Lengths of the arrays do not match.")
		return
	}

	file, err := os.Create("data.csv")
	if err != nil {
		log.Fatalf("Error creating CSV file: %v", err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"Date", "Value"})

	for i := 0; i < len(newCasesArray); i++ {
		writer.Write([]string{datesArray[i], strconv.Itoa(newCasesArray[i])})
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
	str = strings.Trim(str, "[]") 
	strArr := strings.Split(str, ",")

	var intArr []int
	for _, s := range strArr {
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

func getScriptContent(n *html.Node, targetString string) string {
	if n.Type == html.ElementNode && n.Data == "script" {
		for _, attr := range n.Attr {
			if attr.Key == "type" && attr.Val == "text/javascript" {
				content := getTextContent(n)
				if strings.Contains(content, targetString) {
					return content
				}
			}
		}
	}
	var scriptContent string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		scriptContent = getScriptContent(c, targetString)
		if scriptContent != "" {
			return scriptContent
		}
	}
	return ""
}

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

func extractContent(content, regex string) string {
	re := regexp.MustCompile(regex)
	match := re.FindStringSubmatch(content)

	var finalContent string

	if len(match) > 1 {
		categoriesContent := match[1]
		finalContent = "[" + categoriesContent + "]"
	} else {
		finalContent = "Content not found in the input."
	}

	return finalContent
}

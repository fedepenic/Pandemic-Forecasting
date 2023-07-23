package main

import (
	"encoding/csv"
	"net/http"
	"os"
)

func csvHandler(w http.ResponseWriter, r *http.Request) {
	// Locate the data.csv file relative to the project directory
	filePath := "../../data.csv"

	// Open the CSV file
	file, err := os.Open(filePath)
	if err != nil {
		http.Error(w, "Failed to open the CSV file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read the CSV data
	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		http.Error(w, "Failed to read CSV data", http.StatusInternalServerError)
		return
	}

	// Set the appropriate headers for CSV response
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", "inline;filename=data.csv")

	// Write the CSV data to the response
	writer := csv.NewWriter(w)
	defer writer.Flush()

	for _, row := range data {
		writer.Write(row)
	}
}

func main() {
	router := http.NewServeMux()
	router.HandleFunc("/csv", csvHandler)

	http.ListenAndServe(":8080", router)
}

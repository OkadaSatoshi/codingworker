package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	// Open the CSV file
	file, err := os.Open("sample.csv")
	if err != nil {
		fmt.Println("Error opening CSV file:", err)
		return
	}
	defer file.Close()

	// Create a CSV reader
	reader := csv.NewReader(file)

	// Read all records from the CSV file
	allRecords, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV file:", err)
		return
	}

	// Sort the records by the "name" column using integer index 0
	sortedRecords := sortRecords(allRecords, 0)

	// Print the sorted records to stdout
	for _, record := range sortedRecords {
		fmt.Printf("%s\n", strings.Join(record, ","))
	}
}

func sortRecords(records [][]string, key int) [][]string {
	sort.Slice(records, func(i, j int) bool {
		return records[i][key] < records[j][key]
	})
	return records
}

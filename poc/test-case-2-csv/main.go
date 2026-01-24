package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
)

type Person struct {
	Name string
	Age  int
	City string
}

func (p Person) Less(other Person) bool {
	return p.Name < other.Name
}

func main() {
	file, err := os.Open("sample.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Println("Error reading CSV:", err)
		return
	}

	people := make([]Person, len(records)-1) // Skip the header row
	for i, record := range records[1:] {
		age, _ := strconv.Atoi(record[1])
		people[i] = Person{Name: record[0], Age: age, City: record[2]}
	}

	sort.Slice(people, func(i, j int) bool {
		return people[i].Less(people[j])
	})

	for _, person := range people {
		fmt.Printf("%s,%d,%s\n", person.Name, person.Age, person.City)
	}
}

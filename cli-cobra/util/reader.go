package util

import (
	"encoding/csv"
	"fmt"
	"os"
)

type Details struct {
}

func ReadCSV() string {
	// Open the file
	// Read the file
	// Return the data

	fileName := "soccer_data.csv"

	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Error opening the file %s: %v\n", file.Name(), err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		fmt.Printf("Error reading the file %s: %v\n", file.Name(), err)
	}

	return fmt.Sprint(records)

}

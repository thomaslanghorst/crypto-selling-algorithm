package common

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
)

func ReadFile(csvFile string) ([][]string, error) {

	csvfile, err := os.Open(csvFile)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read csv file. Error: %s", err.Error()))
	}

	defer csvfile.Close()

	reader := csv.NewReader(csvfile)
	reader.Comma = ','

	allLines, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("unable to read csv file. Error: %s", err.Error()))
	}

	return allLines, nil

}

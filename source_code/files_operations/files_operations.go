package files_operations

import (
	"aca_hw1/data_structures"
	"encoding/json"
	"fmt"
	"os"
)

// Function to read the input JSON file and return the decoded data
func ReadInputFile(inputFile string) (*data_structures.InputInstructions, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("error reading input file: %v", err)
	}

	var instructionsList []string
	err = json.Unmarshal(data, &instructionsList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	inputInstructions := data_structures.InputInstructions{
		Instructions: instructionsList,
	}

	return &inputInstructions, nil
}

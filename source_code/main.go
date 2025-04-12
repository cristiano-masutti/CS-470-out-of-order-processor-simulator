package main

import (
	"aca_hw1/execution"
	"aca_hw1/files_operations"
	"fmt"
	"log"
	"os"
)

/*
Ideas:

if or at the false: stop
*/
func main() {
	if len(os.Args) < 3 {
		log.Fatal("Please provide path to input and output file")
	}

	inputFile := os.Args[1]
	outputFile := os.Args[2]

	decodedInputInstructions, err := files_operations.ReadInputFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	err = files_operations.CreateOrCleanOutputFile(outputFile)

	processorState := execution.NewProcessorState(decodedInputInstructions)

	err = processorState.SaveState(outputFile)
	if err != nil {
		log.Fatal(err)
	}

	//for i := 0; i < 23; i++ {
	for {
		processorState.Propagate()
		processorState.Latch()

		err = processorState.SaveState(outputFile)
		if err != nil {
			log.Fatal(err)
		}

		if processorState.Exception ||
			(int(processorState.PCP.GetCurrentValue()) == len(decodedInputInstructions) &&
				len(processorState.ActiveList.GetActiveList()) == 0) {
			break
		}
	}

	// Deal with exception if present
	if processorState.Exception {
		for {
			processorState.RecoverExceptionState()
			processorState.Latch()

			err = processorState.SaveState(outputFile)
			if err != nil {
				log.Fatal(err)
			}

			if !processorState.Exception {
				break
			}

			if len(processorState.ActiveList.GetActiveList()) == 0 {
				processorState.Exception = false
			}
		}
	}

	fmt.Println("Output written to", outputFile)
}

package main

import (
	"aca_hw1/execution"
	"aca_hw1/files_operations"
	"fmt"
	"log"
)

/*
Ideas:

if or at the false: stop
*/
func main() {
	//if len(os.Args) < 3 {
	//	log.Fatal("Please provide path to input and output file")
	//}
	//
	//inputFile := os.Args[1]
	//outputFile := os.Args[2]

	inputFile := "../given_tests/01/input.json"
	outputFile := "output.json"

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

	for i := 0; i < 22; i++ {
		log.Printf("cycle: ", i)
		err = processorState.Propagate()
		if err != nil {
			log.Fatal(err)
		}

		err = processorState.Latch()
		if err != nil {
			log.Fatal(err)
		}

		err = processorState.SaveState(outputFile)
		if err != nil {
			log.Fatal(err)
		}
	}

	//if processorState.PC == uint64(len(program)) && len(processorState.ActiveList) == 0 {
	//	break
	//}

	//err = files_operations.WriteOutputFile(outputFile, decodedInputInstructions)
	//if err != nil {
	//	log.Fatal(err)
	//}

	fmt.Println("Output written to", outputFile)
}

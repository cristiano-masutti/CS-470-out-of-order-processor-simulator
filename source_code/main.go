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

	inputInstructions, err := files_operations.ReadInputFile(inputFile)
	if err != nil {
		log.Fatal(err)
	}

	processorState := execution.NewProcessorState(inputInstructions)

	for {
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

		break

		//if processorState.PC == uint64(len(program)) && len(processorState.ActiveList) == 0 {
		//	break
		//}
	}

	//err = files_operations.WriteOutputFile(outputFile, inputInstructions)
	//if err != nil {
	//	log.Fatal(err)
	//}

	fmt.Println("Output written to", outputFile)
}

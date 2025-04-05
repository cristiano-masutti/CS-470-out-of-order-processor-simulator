package files_operations

import (
	"aca_hw1/execution"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// ReadInputFile Function to read the input JSON file and return the decoded instructions
func ReadInputFile(inputFile string) ([]execution.Instruction, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("error reading input file: %v", err)
	}

	var instructionList []string
	err = json.Unmarshal(data, &instructionList)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	var decodedInstructionList []execution.Instruction

	for _, instruction := range instructionList {
		parts := strings.SplitN(instruction, " ", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid instruction format: %s", instruction)
		}
		opcode := parts[0]
		rawOperands := strings.Split(parts[1], ",")

		if len(rawOperands) < 2 {
			return nil, fmt.Errorf("invalid operands in: %s", instruction)
		}
		for i := range rawOperands {
			rawOperands[i] = strings.TrimSpace(rawOperands[i])
		}

		switch opcode {
		case "add", "sub", "mulu", "divu", "remu":
			if len(rawOperands) != 3 {
				return nil, fmt.Errorf("expected 3 operands for %s: %s", opcode, instruction)
			}

			dest, err1 := strconv.Atoi(strings.TrimPrefix(rawOperands[0], "x"))
			opA, err2 := strconv.Atoi(strings.TrimPrefix(rawOperands[1], "x"))
			opB, err3 := strconv.Atoi(strings.TrimPrefix(rawOperands[2], "x"))

			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("error parsing registers in: %s", instruction)
			}

			base := execution.BaseInstruction{Dest: dest, OpA: opA}
			var inst execution.Instruction

			switch opcode {
			case "add":
				inst = &execution.Add{BaseInstruction: base, OpB: opB}
			case "sub":
				inst = &execution.Sub{BaseInstruction: base, OpB: opB}
			case "mulu":
				inst = &execution.Mulu{BaseInstruction: base, OpB: opB}
			case "divu":
				inst = &execution.Divu{BaseInstruction: base, OpB: opB}
			case "remu":
				inst = &execution.Remu{BaseInstruction: base, OpB: opB}
			}

			decodedInstructionList = append(decodedInstructionList, inst)

		case "addi":
			if len(rawOperands) != 3 {
				return nil, fmt.Errorf("expected 3 operands for addi: %s", instruction)
			}

			dest, err1 := strconv.Atoi(strings.TrimPrefix(rawOperands[0], "x"))
			opA, err2 := strconv.Atoi(strings.TrimPrefix(rawOperands[1], "x"))
			imm, err3 := strconv.Atoi(rawOperands[2])

			if err1 != nil || err2 != nil || err3 != nil {
				return nil, fmt.Errorf("error parsing addi instruction: %s", instruction)
			}

			inst := &execution.Addi{
				BaseInstruction: execution.BaseInstruction{Dest: dest, OpA: opA},
				Imm:             imm,
			}

			decodedInstructionList = append(decodedInstructionList, inst)

		default:
			return nil, fmt.Errorf("unrecognized instruction: %s", opcode)
		}
	}

	return decodedInstructionList, nil
}

func CreateOrCleanOutputFile(outputFile string) error {
	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open or create file: %w", err)
	}
	defer file.Close()
	return nil
}

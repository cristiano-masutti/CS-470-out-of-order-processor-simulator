package execution

import (
	"fmt"
	"strconv"
	"strings"
)

func (ps *ProcessorState) Propagate() error {
	return nil
}

func (ps *ProcessorState) FetchAndDecode() error {
	// Determine how many
	numberOfInstructions := 4

	//get current PC
	pc := int(ps.PCP.CurrentValue)
	ps.PCP.SetNextValue(uint64(numberOfInstructions))

	var nextInstructions []Instruction

	// Decode
	for i := 0; i < numberOfInstructions; i++ {
		actualPc := pc + i
		currentStringInstruction := ps.InputInstructions.GetInstruction(actualPc)
		parts := strings.Split(currentStringInstruction, " ")

		identifier := parts[0]
		operands := strings.Split(parts[1], ", ")

		switch identifier {
		case "add":
			currentInstruction := &Add{BaseInstruction: BaseInstruction{PC: actualPc}, Dest: operands[0], OpA: operands[1], OpB: operands[2]}
			nextInstructions = append(nextInstructions, currentInstruction)
		case "addi":
			imm64, _ := strconv.ParseInt(operands[2], 10, 64)
			imm := int(imm64)
			currentInstruction := &Addi{BaseInstruction: BaseInstruction{PC: actualPc}, Dest: operands[0], OpA: operands[1], Imm: imm}
			nextInstructions = append(nextInstructions, currentInstruction)
		case "sub":
			currentInstruction := &Sub{BaseInstruction: BaseInstruction{PC: actualPc}, Dest: operands[0], OpA: operands[1], OpB: operands[2]}
			nextInstructions = append(nextInstructions, currentInstruction)
		case "mulu":
			currentInstruction := &Mulu{BaseInstruction: BaseInstruction{PC: actualPc}, Dest: operands[0], OpA: operands[1], OpB: operands[2]}
			nextInstructions = append(nextInstructions, currentInstruction)
		case "divu":
			currentInstruction := &Divu{BaseInstruction: BaseInstruction{PC: actualPc}, Dest: operands[0], OpA: operands[1], OpB: operands[2]}
			nextInstructions = append(nextInstructions, currentInstruction)
		case "remu":
			currentInstruction := &Remu{BaseInstruction: BaseInstruction{PC: actualPc}, Dest: operands[0], OpA: operands[1], OpB: operands[2]}
			nextInstructions = append(nextInstructions, currentInstruction)
		default:
			return fmt.Errorf("Unrecognized instruction \"%s\"", identifier)
		}
	}

	ps.DPR.SetNextValue(nextInstructions)

	return nil
}

func (ps *ProcessorState) Latch() error {
	return nil
}

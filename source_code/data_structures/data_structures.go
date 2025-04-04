package data_structures

import (
	"encoding/json"
	"fmt"
	"os"
)

// Input

type InputInstructions struct {
	Instructions []string
}

// SaveState output file format
type SaveState struct {
	PC                   uint64              `json:"PC"`
	PhysicalRegisterFile [64]uint64          `json:"PhysicalRegisterFile"`
	DecodedPCs           []uint64            `json:"DecodedPCs"`
	Exception            bool                `json:"Exception"`
	ExceptionPC          uint64              `json:"ExceptionPC"`
	RegisterMapTable     [32]uint8           `json:"RegisterMapTable"`
	FreeList             []uint8             `json:"FreeList"`
	BusyBitTable         [64]bool            `json:"BusyBitTable"`
	ActiveList           []ActiveListEntry   `json:"ActiveList"`
	IntegerQueue         []IntegerQueueEntry `json:"IntegerQueue"`
}

// ProcessorState represents state of the processor
type ProcessorState struct {
	InputInstructions    *InputInstructions
	PCP                  *PCPipelineRegister
	PhysicalRegisterFile [64]uint64          `json:"PhysicalRegisterFile"`
	DecodedPCs           []uint64            `json:"DecodedPCs"`
	Exception            bool                `json:"Exception"`
	ExceptionPC          uint64              `json:"ExceptionPC"`
	RegisterMapTable     [32]uint8           `json:"RegisterMapTable"`
	FreeList             []uint8             `json:"FreeList"`
	BusyBitTable         [64]bool            `json:"BusyBitTable"`
	ActiveList           []ActiveListEntry   `json:"ActiveList"`
	IntegerQueue         []IntegerQueueEntry `json:"IntegerQueue"`
}

// NewProcessorState create new ProcessorState
func NewProcessorState(instructions *InputInstructions) *ProcessorState {
	ps := &ProcessorState{
		InputInstructions: instructions,
		PCP: &PCPipelineRegister{
			CurrentValue: 0,
			NewValue:     0,
		},
		DecodedPCs:   make([]uint64, 0),
		Exception:    false,
		ExceptionPC:  0,
		FreeList:     make([]uint8, 32),
		ActiveList:   make([]ActiveListEntry, 0),
		IntegerQueue: make([]IntegerQueueEntry, 0),
	}

	for i := range ps.PhysicalRegisterFile {
		ps.PhysicalRegisterFile[i] = 0
	}

	for i := range ps.RegisterMapTable {
		ps.RegisterMapTable[i] = uint8(i)
	}

	for i := range ps.FreeList {
		ps.FreeList[i] = uint8(i + 32)
	}

	for i := range ps.BusyBitTable {
		ps.BusyBitTable[i] = false
	}

	return ps
}

// SaveState writes processor state to JSON file
func (ps *ProcessorState) SaveState(filename string) error {
	// Step 1: Extract the current state into SaveState struct
	var pc uint64
	if ps.PCP != nil {
		pc = ps.PCP.CurrentValue
	}

	newState := SaveState{
		PC:                   pc,
		PhysicalRegisterFile: ps.PhysicalRegisterFile,
		DecodedPCs:           ps.DecodedPCs,
		Exception:            ps.Exception,
		ExceptionPC:          ps.ExceptionPC,
		RegisterMapTable:     ps.RegisterMapTable,
		FreeList:             ps.FreeList,
		BusyBitTable:         ps.BusyBitTable,
		ActiveList:           ps.ActiveList,
		IntegerQueue:         ps.IntegerQueue,
	}

	// Step 2: Read existing data (if file exists)
	var existingData []SaveState
	if _, err := os.Stat(filename); err == nil {
		fileData, err := os.ReadFile(filename)
		if err != nil {
			return fmt.Errorf("error reading file: %v", err)
		}
		if len(fileData) > 0 {
			if err := json.Unmarshal(fileData, &existingData); err != nil {
				return fmt.Errorf("error decoding existing JSON: %v", err)
			}
		}
	}

	// Step 3: Append new state
	updatedData := append(existingData, newState)

	// Step 4: Write back to file (indented for readability)
	jsonData, err := json.MarshalIndent(updatedData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	return os.WriteFile(filename, jsonData, 0644)
}

func (ps *ProcessorState) Propagate() error {
	return nil
}

func (ps *ProcessorState) Latch() error {
	return nil
}

// ActiveListEntry represents entry in Active List
type ActiveListEntry struct {
	Done               bool   `json:"Done"`
	Exception          bool   `json:"Exception"`
	LogicalDestination uint8  `json:"LogicalDestination"`
	OldDestination     uint8  `json:"OldDestination"`
	PC                 uint64 `json:"PC"`
}

// IntegerQueueEntry represents entry in Integer Queue
type IntegerQueueEntry struct {
	DestRegister uint8  `json:"DestRegister"`
	OpAIsReady   bool   `json:"OpAIsReady"`
	OpARegTag    uint8  `json:"OpARegTag"`
	OpAValue     uint64 `json:"OpAValue"`
	OpBIsReady   bool   `json:"OpBIsReady"`
	OpBRegTag    uint8  `json:"OpBRegTag"`
	OpBValue     uint64 `json:"OpBValue"`
	OpCode       string `json:"OpCode"`
	PC           uint64 `json:"PC"`
}

// Pipeline registers:

// PCPipelineRegister represents pipeline register for pc stage
type PCPipelineRegister struct {
	CurrentValue uint64
	NewValue     uint64
}

func (cpl *PCPipelineRegister) FetchedOperations(incNumber uint64) {
	cpl.NewValue = cpl.CurrentValue + incNumber
}

func (cpl *PCPipelineRegister) LatchPCPipelineRegister() {
	cpl.CurrentValue = cpl.NewValue
}

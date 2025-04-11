package execution

import (
	"encoding/json"
	"fmt"
	"os"
)

// ProcessorState represents state of the execution
type ProcessorState struct {
	InputInstructions                 []Instruction
	PCP                               *PCPipelineRegister
	DPR                               *DirPipelineRegister
	PhysicalRegisterFile              [64]uint64
	Exception                         *ExceptionFlag
	ExceptionPC                       *PCPipelineRegister
	RegisterMapTable                  [32]uint64
	FreeList                          *FreeList
	BusyBitTable                      *BusyBitTable
	ActiveList                        *ActiveList
	IntegerQueue                      *IntegerQueue
	IssuedInstructionPipelineRegister *InstructionsPipelineRegister
	AluPipelineRegisters              *InstructionsPipelineRegister
	ForwardingPaths                   *ForwardingPaths
	CommitPipeline                    *CommitPipelineRegister
}

// NewProcessorState create new ProcessorState
func NewProcessorState(instructions []Instruction) *ProcessorState {
	ps := &ProcessorState{
		InputInstructions:                 instructions,
		PCP:                               NewPCPipelineRegister(),
		DPR:                               NewDirPipelineRegister(),
		Exception:                         NewExceptionFlag(),
		ExceptionPC:                       NewPCPipelineRegister(),
		FreeList:                          NewFreeList(),
		ActiveList:                        NewActiveList(),
		IntegerQueue:                      NewIntegerQueue(),
		BusyBitTable:                      NewBusyTable(),
		IssuedInstructionPipelineRegister: NewInstructionPipelineRegister(),
		AluPipelineRegisters:              NewInstructionPipelineRegister(),
		ForwardingPaths:                   NewForwardingPaths(),
		CommitPipeline:                    NewCommitPipelineRegister(),
	}

	for i := range ps.PhysicalRegisterFile {
		ps.PhysicalRegisterFile[i] = 0
	}

	for i := range ps.RegisterMapTable {
		ps.RegisterMapTable[i] = uint64(i)
	}

	return ps
}

///////////////////////////////////////////////////////////////////////////

// SaveState output file format
type SaveState struct {
	ActiveList           []ActiveListEntry   `json:"ActiveList"`
	BusyBitTable         []bool              `json:"BusyBitTable"`
	DecodedPCs           []uint64            `json:"DecodedPCs"`
	Exception            bool                `json:"Exception"`
	ExceptionPC          uint64              `json:"ExceptionPC"`
	FreeList             []uint64            `json:"FreeList"`
	IntegerQueue         []IntegerQueueEntry `json:"IntegerQueue"`
	PC                   uint64              `json:"PC"`
	PhysicalRegisterFile [64]uint64          `json:"PhysicalRegisterFile"`
	RegisterMapTable     [32]uint64          `json:"RegisterMapTable"`
}

func (ps *ProcessorState) SaveState(filename string) error {
	newState := SaveState{
		ActiveList:           ps.ActiveList.GetActiveList(),
		BusyBitTable:         ps.BusyBitTable.GetBusyBitTable(),
		DecodedPCs:           ps.DPR.GetCurrentValue(),
		Exception:            ps.Exception.GetCurrentStatus(),
		ExceptionPC:          ps.ExceptionPC.GetCurrentValue(),
		FreeList:             ps.FreeList.GetFreeList(),
		IntegerQueue:         ps.IntegerQueue.GetCurrentIntegerQueue(),
		PC:                   ps.PCP.GetCurrentValue(),
		PhysicalRegisterFile: ps.PhysicalRegisterFile,
		RegisterMapTable:     ps.RegisterMapTable,
	}

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

	updatedData := append(existingData, newState)

	jsonData, err := json.MarshalIndent(updatedData, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshaling data: %v", err)
	}

	return os.WriteFile(filename, jsonData, 0644)
}

package execution

import (
	"encoding/json"
	"fmt"
	"os"
)

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

// ProcessorState represents state of the execution
type ProcessorState struct {
	InputInstructions    []Instruction
	PCP                  *PCPipelineRegister
	DPR                  *DirPipelineRegister
	PhysicalRegisterFile [64]uint64
	Exception            bool
	ExceptionPC          uint64
	RegisterMapTable     [32]uint64
	FreeList             *FreeList
	BusyBitTable         *BusyBitTable
	ActiveList           *ActiveList
	IntegerQueue         *IntegerQueue
}

// NewProcessorState create new ProcessorState
func NewProcessorState(instructions []Instruction) *ProcessorState {
	ps := &ProcessorState{
		InputInstructions: instructions,
		PCP: &PCPipelineRegister{
			CurrentValue: 0,
			NewValue:     0,
		},
		DPR: &DirPipelineRegister{
			CurrentDecodedInstructions: make([]uint64, 0),
			NewDecodedInstructions:     make([]uint64, 0),
			BackPressure:               false,
		},
		Exception:    false,
		ExceptionPC:  0,
		FreeList:     NewFreeList(),
		ActiveList:   NewActiveList(),
		IntegerQueue: NewIntegerQueue(),
		BusyBitTable: NewBusyTable(),
	}

	for i := range ps.PhysicalRegisterFile {
		ps.PhysicalRegisterFile[i] = 0
	}

	for i := range ps.RegisterMapTable {
		ps.RegisterMapTable[i] = uint64(i)
	}

	return ps
}

// SaveState writes execution state to JSON file
func (ps *ProcessorState) SaveState(filename string) error {
	var pc uint64

	if ps.PCP != nil {
		pc = ps.PCP.CurrentValue
	}

	newState := SaveState{
		ActiveList:           ps.ActiveList.GetActiveList(),
		BusyBitTable:         ps.BusyBitTable.GetBusyBitTable(),
		DecodedPCs:           ps.DPR.GetCurrentValue(),
		Exception:            ps.Exception,
		ExceptionPC:          ps.ExceptionPC,
		FreeList:             ps.FreeList.GetFreeList(),
		IntegerQueue:         ps.IntegerQueue.GetIntegerQueue(),
		PC:                   pc,
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

type BusyBitTable struct {
	BusyTableEntries []bool
}

func NewBusyTable() *BusyBitTable {
	return &BusyBitTable{make([]bool, 64)}
}

func (bt *BusyBitTable) GetBusyBitTable() []bool {
	return bt.BusyTableEntries
}

// GetRegisterState returns the state of register
// if returns true, means that value not ready
// if false, it is ready
func (bt *BusyBitTable) GetRegisterState(index int) bool {
	return bt.BusyTableEntries[index]
}

func (bt *BusyBitTable) SetRegisterState(index int, b bool) {
	bt.BusyTableEntries[index] = b
}

type ActiveList struct {
	ActiveListEntries []ActiveListEntry
}

func NewActiveList() *ActiveList {
	return &ActiveList{make([]ActiveListEntry, 0)}
}

func (al *ActiveList) Append(entry ActiveListEntry) {
	al.ActiveListEntries = append(al.ActiveListEntries, entry)
}

func (al *ActiveList) GetActiveList() []ActiveListEntry {
	return al.ActiveListEntries
}

type IntegerQueue struct {
	IntegerQueueEntries []IntegerQueueEntry
}

func NewIntegerQueue() *IntegerQueue {
	return &IntegerQueue{make([]IntegerQueueEntry, 0)}
}

func (qe *IntegerQueue) Append(entry IntegerQueueEntry) {
	qe.IntegerQueueEntries = append(qe.IntegerQueueEntries, entry)
}

func (qe *IntegerQueue) GetIntegerQueue() []IntegerQueueEntry {
	return qe.IntegerQueueEntries
}

// FreeList to handle free list operation
type FreeList struct {
	registerList []uint64
}

func NewFreeList() *FreeList {
	registerList := make([]uint64, 32)

	for i := range registerList {
		registerList[i] = uint64(i + 32)
	}
	return &FreeList{
		registerList: registerList,
	}
}

func (fl *FreeList) GetFreeList() []uint64 {
	return fl.registerList
}

// GetRegister returns and removes the first register in FIFO order.
func (fl *FreeList) GetRegister() uint64 {
	reg := fl.registerList[0]
	fl.registerList = fl.registerList[1:]
	return reg
}

// FreeRegister appends the register ID to the end of the list.
func (fl *FreeList) FreeRegister(id uint64) {
	fl.registerList = append(fl.registerList, id)
}

// ActiveListEntry represents entry in Active List
type ActiveListEntry struct {
	Done               bool   `json:"Done"`
	Exception          bool   `json:"Exception"`
	LogicalDestination uint64 `json:"LogicalDestination"`
	OldDestination     uint64 `json:"OldDestination"`
	PC                 uint64 `json:"PC"`
}

// IntegerQueueEntry represents entry in Integer Queue
type IntegerQueueEntry struct {
	DestRegister uint64 `json:"DestRegister"`
	OpAIsReady   bool   `json:"OpAIsReady"`
	OpARegTag    uint64 `json:"OpARegTag"`
	OpAValue     uint64 `json:"OpAValue"`
	OpBIsReady   bool   `json:"OpBIsReady"`
	OpBRegTag    uint64 `json:"OpBRegTag"`
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

func (pcpr *PCPipelineRegister) SetNextValue(incNumber uint64) {
	pcpr.NewValue = pcpr.CurrentValue + incNumber
}

func (pcpr *PCPipelineRegister) LatchPCPipelineRegister() {
	pcpr.CurrentValue = pcpr.NewValue
}

func (pcpr *PCPipelineRegister) GetCurrentValue() uint64 {
	return pcpr.CurrentValue
}

type DirPipelineRegister struct {
	CurrentDecodedInstructions []uint64
	NewDecodedInstructions     []uint64
	BackPressure               bool
}

func (dpr *DirPipelineRegister) SetBackPressure(BackPressure bool) {
	dpr.BackPressure = BackPressure
}

func (dpr *DirPipelineRegister) GetBackPressure() bool {
	return dpr.BackPressure
}

func (dpr *DirPipelineRegister) SetNextValue(newDecodedInstructions []uint64) {
	dpr.NewDecodedInstructions = newDecodedInstructions
}

func (dpr *DirPipelineRegister) LatchPCPipelineRegister() {
	dpr.CurrentDecodedInstructions = dpr.NewDecodedInstructions
}

func (dpr *DirPipelineRegister) GetCurrentValue() []uint64 {
	return dpr.CurrentDecodedInstructions
}

// Instructions types

type Instruction interface {
	GetDest() int
	GetOpA() int
	GetOpCode() string
	GetSecondArg() int
}

type BaseInstruction struct {
	Dest int
	OpA  int
}

func (b *BaseInstruction) GetDest() int { return b.Dest }
func (b *BaseInstruction) GetOpA() int  { return b.OpA }
func (b *BaseInstruction) GetOpCode() string {
	return "Base Instruction"
}
func (b *BaseInstruction) GetSecondArg() int {
	return b.OpA
}

type Add struct {
	BaseInstruction
	OpB int
}

func (a *Add) GetOpCode() string { return "add" }
func (a *Add) GetSecondArg() int { return a.OpB }

type Addi struct {
	BaseInstruction
	Imm int
}

func (a *Addi) GetOpCode() string { return "addi" }
func (a *Addi) GetSecondArg() int { return a.Imm }

type Sub struct {
	BaseInstruction
	OpB int
}

func (s *Sub) GetOpCode() string { return "sub" }
func (s *Sub) GetSecondArg() int { return s.OpB }

type Mulu struct {
	BaseInstruction
	OpB int
}

func (m *Mulu) GetOpCode() string { return "mulu" }
func (m *Mulu) GetSecondArg() int { return m.OpB }

type Divu struct {
	BaseInstruction
	OpB int
}

func (d *Divu) GetOpCode() string { return "divu" }
func (d *Divu) GetSecondArg() int { return d.OpB }

type Remu struct {
	BaseInstruction
	OpB int
}

func (r *Remu) GetOpCode() string { return "remu" }
func (r *Remu) GetSecondArg() int { return r.OpB }

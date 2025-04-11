package execution

import (
	"sort"
)

///////////////////////////////////////////////////////////////////////////

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

func (fl *FreeList) GetIfSpaceAvailable() bool {
	return len(fl.registerList) >= 4
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

///////////////////////////////////////////////////////////////////////////

// BusyBitTable represents the busy bit table
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

///////////////////////////////////////////////////////////////////////////

// ActiveListEntry represents entry in Active List
type ActiveListEntry struct {
	Done               bool   `json:"Done"`
	Exception          bool   `json:"Exception"`
	LogicalDestination uint64 `json:"LogicalDestination"`
	OldDestination     uint64 `json:"OldDestination"`
	PC                 uint64 `json:"PC"`
}

type ActiveList struct {
	CurrentActiveListEntries []ActiveListEntry
	NextActiveListEntries    []ActiveListEntry
}

func NewActiveList() *ActiveList {
	return &ActiveList{
		CurrentActiveListEntries: make([]ActiveListEntry, 0),
		NextActiveListEntries:    make([]ActiveListEntry, 0)}
}

func (al *ActiveList) Append(entry ActiveListEntry) {
	al.NextActiveListEntries = append(al.NextActiveListEntries, entry)
}

func (al *ActiveList) GetActiveList() []ActiveListEntry {
	original := al.CurrentActiveListEntries
	copied := make([]ActiveListEntry, len(original))
	copy(copied, original)
	return copied
}

func (al *ActiveList) SetDoneBitForInstructions(executedInstructions []CommitPipelineRegisterEntry) {
	for _, instr := range executedInstructions {
		for i := range al.NextActiveListEntries {
			if al.NextActiveListEntries[i].PC == instr.PC {
				al.NextActiveListEntries[i].Done = instr.Done
				al.NextActiveListEntries[i].Exception = instr.Exception
				break
			}
		}
	}
}

func (al *ActiveList) GetIfSpaceAvailable() bool {
	return len(al.CurrentActiveListEntries)+4 <= 32
}

func (al *ActiveList) Latch() {
	al.CurrentActiveListEntries = al.NextActiveListEntries
}

func (al *ActiveList) RetireInstructions(numberOfInstructions int) {
	al.NextActiveListEntries = al.NextActiveListEntries[numberOfInstructions:]
}

///////////////////////////////////////////////////////////////////////////

// IntegerQueueEntry represents entry in Integer Queue
type IntegerQueueEntry struct {
	DestRegister uint64 `json:"DestRegister"`
	OpAIsReady   bool   `json:"OpAIsReady"`
	OpARegTag    uint64 `json:"OpARegTag"`
	OpAValue     int    `json:"OpAValue"`
	OpBIsReady   bool   `json:"OpBIsReady"`
	OpBRegTag    uint64 `json:"OpBRegTag"`
	OpBValue     int    `json:"OpBValue"`
	OpCode       string `json:"OpCode"`
	PC           uint64 `json:"PC"`
}

type IntegerQueue struct {
	CurrentIntegerQueueEntries []IntegerQueueEntry
	NextIntegerQueueEntries    []IntegerQueueEntry
}

func NewIntegerQueue() *IntegerQueue {
	return &IntegerQueue{
		CurrentIntegerQueueEntries: make([]IntegerQueueEntry, 0),
		NextIntegerQueueEntries:    make([]IntegerQueueEntry, 0),
	}
}

func (qe *IntegerQueue) Reset() {
	qe.NextIntegerQueueEntries = make([]IntegerQueueEntry, 0)
}

func (qe *IntegerQueue) Append(entry IntegerQueueEntry) {
	qe.NextIntegerQueueEntries = append(qe.NextIntegerQueueEntries, entry)
}

func (qe *IntegerQueue) Latch() {
	qe.CurrentIntegerQueueEntries = qe.NextIntegerQueueEntries
}

func (qe *IntegerQueue) GetIfSpaceAvailable() bool {
	return len(qe.CurrentIntegerQueueEntries)+4 <= 32
}

func (qe *IntegerQueue) ForwardResults(results []ForwardingPathsEntry) {
	for i, entry := range qe.CurrentIntegerQueueEntries {
		for _, result := range results {
			if !entry.OpAIsReady && entry.OpARegTag == result.tag {
				entry.OpAIsReady = true
				entry.OpAValue = int(result.value)
				entry.OpARegTag = 0
			}
			if !entry.OpBIsReady && entry.OpBRegTag == result.tag {
				entry.OpBIsReady = true
				entry.OpBValue = int(result.value)
				entry.OpBRegTag = 0
			}
		}
		qe.CurrentIntegerQueueEntries[i] = entry
	}
}

func (qe *IntegerQueue) GetReadyInstructions() []IntegerQueueEntry {
	ready := make([]IntegerQueueEntry, 0)

	for _, entry := range qe.CurrentIntegerQueueEntries {
		if entry.OpAIsReady && entry.OpBIsReady {
			ready = append(ready, entry)
		}
	}

	if len(ready) > 4 {
		sort.Slice(ready, func(i, j int) bool {
			return ready[i].PC < ready[j].PC
		})
		ready = ready[:4]
	}

	filteredNext := make([]IntegerQueueEntry, 0, len(qe.NextIntegerQueueEntries))
	for _, nextEntry := range qe.NextIntegerQueueEntries {
		keep := true
		for _, selected := range ready {
			if nextEntry == selected {
				keep = false
				break
			}
		}
		if keep {
			filteredNext = append(filteredNext, nextEntry)
		}
	}
	qe.NextIntegerQueueEntries = filteredNext

	return ready
}

func (qe *IntegerQueue) GetCurrentIntegerQueue() []IntegerQueueEntry {
	return qe.CurrentIntegerQueueEntries
}

///////////////////////////////////////////////////////////////////////////

type ExceptionFlag struct {
	currentStatus bool
	nextStatus    bool
}

func NewExceptionFlag() *ExceptionFlag {
	return &ExceptionFlag{
		currentStatus: false,
		nextStatus:    false,
	}
}

func (ef *ExceptionFlag) GetCurrentStatus() bool {
	return ef.currentStatus
}

func (ef *ExceptionFlag) SetNextStatus(nextStatus bool) {
	ef.nextStatus = nextStatus
}

func (ef *ExceptionFlag) Latch() {
	ef.currentStatus = ef.nextStatus
}

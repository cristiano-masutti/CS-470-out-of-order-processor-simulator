package execution

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

///////////////////////////////////////////////////////////////////////////

// DirPipelineRegister represents pipeline register for dir stage
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

///////////////////////////////////////////////////////////////////////////

// IssuedInstructionPipelineRegister represents pipeline register between Issue and Execution
type IssuedInstructionPipelineRegister struct {
	currentIssuedInstruction []IntegerQueueEntry
	nextIssuedInstruction    []IntegerQueueEntry
}

func NewIssuedInstructionPipelineRegister() *IssuedInstructionPipelineRegister {
	return &IssuedInstructionPipelineRegister{
		currentIssuedInstruction: make([]IntegerQueueEntry, 0),
		nextIssuedInstruction:    make([]IntegerQueueEntry, 0),
	}
}

func (irp *IssuedInstructionPipelineRegister) SetNextSetOfInstructions(entries []IntegerQueueEntry) {
	irp.nextIssuedInstruction = entries
}

func (irp *IssuedInstructionPipelineRegister) GetCurrentIssuedInstructions() []IntegerQueueEntry {
	return irp.currentIssuedInstruction
}

func (irp *IssuedInstructionPipelineRegister) Latch() {
	irp.currentIssuedInstruction = irp.nextIssuedInstruction
}

///////////////////////////////////////////////////////////////////////////

// AluPipelineRegisters represents pipeline register in Execution stage
type AluPipelineRegisters struct {
	currentExecutingInstructions []IntegerQueueEntry
	nextExecutingInstructions    []IntegerQueueEntry
}

func NewAluPipelineRegisters() *AluPipelineRegisters {
	return &AluPipelineRegisters{
		currentExecutingInstructions: make([]IntegerQueueEntry, 0),
		nextExecutingInstructions:    make([]IntegerQueueEntry, 0),
	}
}

func (apl *AluPipelineRegisters) SetNextExecutingInstructions(entries []IntegerQueueEntry) {
	apl.nextExecutingInstructions = entries
}

func (apl *AluPipelineRegisters) GetNextExecutingInstructions() []IntegerQueueEntry {
	return apl.nextExecutingInstructions
}

func (apl *AluPipelineRegisters) Latch() {
	apl.currentExecutingInstructions = apl.nextExecutingInstructions
}

///////////////////////////////////////////////////////////////////////////

type ForwardingPathsEntry struct {
	tag   uint64
	value uint64
}

type ForwardingPaths struct {
	completedInstructions []ForwardingPathsEntry
}

func NewForwardingPaths() *ForwardingPaths {
	return &ForwardingPaths{
		completedInstructions: make([]ForwardingPathsEntry, 0),
	}
}

func (fq *ForwardingPaths) SetCompletedInstruction(instructions []ForwardingPathsEntry) {
	fq.completedInstructions = instructions
}

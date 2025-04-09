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
	if len(dpr.CurrentDecodedInstructions) == 0 {
		return []uint64{}
	}
	return dpr.CurrentDecodedInstructions
}

///////////////////////////////////////////////////////////////////////////

// InstructionsPipelineRegister represents pipeline register in form of IntegerQueueEntry
type InstructionsPipelineRegister struct {
	currentInstructions []IntegerQueueEntry
	nextInstructions    []IntegerQueueEntry
}

func NewInstructionPipelineRegister() *InstructionsPipelineRegister {
	return &InstructionsPipelineRegister{
		currentInstructions: make([]IntegerQueueEntry, 0),
		nextInstructions:    make([]IntegerQueueEntry, 0),
	}
}

func (irp *InstructionsPipelineRegister) SetNextInstructions(entries []IntegerQueueEntry) {
	irp.nextInstructions = entries
}

func (irp *InstructionsPipelineRegister) GetCurrentInstructions() []IntegerQueueEntry {
	return irp.currentInstructions
}

func (irp *InstructionsPipelineRegister) Latch() {
	irp.currentInstructions = irp.nextInstructions
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

func (fq *ForwardingPaths) GetCompletedInstructions() []ForwardingPathsEntry {
	return fq.completedInstructions
}

///////////////////////////////////////////////////////////////////////////

type CommitPipelineRegisterEntry struct {
	Done      bool
	Exception bool
	PC        uint64
}

type CommitPipelineRegister struct {
	CurrentRegister []CommitPipelineRegisterEntry
	NextRegister    []CommitPipelineRegisterEntry
}

func NewCommitPipelineRegister() *CommitPipelineRegister {
	return &CommitPipelineRegister{
		CurrentRegister: make([]CommitPipelineRegisterEntry, 0),
		NextRegister:    make([]CommitPipelineRegisterEntry, 0),
	}
}

func (cr *CommitPipelineRegister) SetNextRegister(newRegister []CommitPipelineRegisterEntry) {
	cr.NextRegister = newRegister
}

func (cr *CommitPipelineRegister) GetCurrentRegister() []CommitPipelineRegisterEntry {
	return cr.CurrentRegister
}

func (cr *CommitPipelineRegister) Latch() {
	cr.CurrentRegister = cr.NextRegister
}

package execution

func (ps *ProcessorState) FetchAndDecodeExceptionFlow() {
	ps.PCP.SetNextValue(65536)
	ps.DPR.SetNextValue(make([]uint64, 0))
	return
}

func (ps *ProcessorState) RenameAndDispatchExceptionFlow() {
	ps.IntegerQueue.Reset()
}

func (ps *ProcessorState) IssueExceptionFlow() {
	ps.IssuedInstructionPipelineRegister.SetNextInstructions(make([]IntegerQueueEntry, 0))
}

func (ps *ProcessorState) ExecuteExceptionFlow() {
	ps.AluPipelineRegisters.SetNextInstructions(make([]IntegerQueueEntry, 0))

	ps.CommitPipeline.SetNextRegister(make([]CommitPipelineRegisterEntry, 0))
}

func (ps *ProcessorState) CommitExceptionFlow() {
	
}

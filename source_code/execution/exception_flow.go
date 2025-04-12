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

func (ps *ProcessorState) RecoverExceptionState() {
	// get 4 instructions per cycle from the iq
	instructions := ps.ActiveList.GetBottomInstructionsInReverseOrder()

	// for each instruction reset everything
	for _, instr := range instructions {
		// get current assigned register
		logicalDest := instr.LogicalDestination
		physicalDest := ps.RegisterMapTable[logicalDest]

		// set busy table bit to false
		ps.BusyBitTable.SetRegisterState(int(physicalDest), false)

		// free register
		ps.FreeList.FreeRegister(physicalDest)

		// assign in mapping arch - physical the previous
		ps.RegisterMapTable[logicalDest] = instr.OldDestination
	}
}

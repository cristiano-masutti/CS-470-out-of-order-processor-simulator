package execution

func (ps *ProcessorState) Latch() error {
	ps.PCP.LatchPCPipelineRegister()
	ps.DPR.LatchPCPipelineRegister()
	ps.IntegerQueue.Latch()
	ps.IssuedInstructionPipelineRegister.Latch()
	ps.ActiveList.Latch()
	return nil
}

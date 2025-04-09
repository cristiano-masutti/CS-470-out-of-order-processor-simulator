package execution

func (ps *ProcessorState) Latch() error {
	ps.PCP.LatchPCPipelineRegister()
	ps.DPR.LatchPCPipelineRegister()
	ps.ActiveList.Latch()
	ps.IntegerQueue.Latch()
	ps.IssuedInstructionPipelineRegister.Latch()
	ps.AluPipelineRegisters.Latch()
	ps.CommitPipeline.Latch()
	return nil
}

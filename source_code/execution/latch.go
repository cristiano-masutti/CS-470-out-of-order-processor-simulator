package execution

func (ps *ProcessorState) Latch() error {
	ps.PCP.LatchPCPipelineRegister()
	ps.DPR.LatchPCPipelineRegister()
	return nil
}

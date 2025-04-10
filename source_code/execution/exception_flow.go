package execution

func (ps *ProcessorState) FetchAndDecodeExceptionFlow() {
	if ps.DPR.GetBackPressure() {
		return
	}

	pc := ps.PCP.CurrentValue
	numberOfInstructions := uint64(min(4, len(ps.InputInstructions)-int(pc)))
	ps.PCP.SetNextValue(numberOfInstructions)

	var nextDecodedPCs []uint64

	for i := uint64(0); i < numberOfInstructions; i++ {
		nextDecodedPCs = append(nextDecodedPCs, pc+i)
	}

	ps.DPR.SetNextValue(nextDecodedPCs)
}

func (ps *ProcessorState) CommitExceptionFlow() {
	finishedInstructions := ps.CommitPipeline.GetCurrentRegister()

	activeList := ps.ActiveList.GetActiveList()

	for _, instr := range finishedInstructions {
		for i := range activeList {
			if activeList[i].PC == instr.PC {
				activeList[i].Done = instr.Done
				activeList[i].Exception = instr.Exception
				ps.FreeList.FreeRegister(activeList[i].OldDestination)
				break
			}
		}
	}

	retired := 0
	for _, instr := range activeList {
		if !instr.Done {
			break
		}

		if instr.Exception {
			// TODO: handle exception
			break
		}

		retired++

		if retired == 4 {
			break
		}
	}

	ps.ActiveList.RetireInstructions(retired)
}

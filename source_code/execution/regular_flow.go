package execution

func (ps *ProcessorState) FetchAndDecodeRegularFlow() {
	if ps.DPR.GetBackPressure() {
		return
	}

	pc := ps.PCP.CurrentValue
	numberOfInstructions := uint64(min(4, len(ps.InputInstructions)-int(pc)))
	ps.PCP.SetNextValue(pc + numberOfInstructions)

	var nextDecodedPCs []uint64

	for i := uint64(0); i < numberOfInstructions; i++ {
		nextDecodedPCs = append(nextDecodedPCs, pc+i)
	}

	ps.DPR.SetNextValue(nextDecodedPCs)
}

func (ps *ProcessorState) RenameAndDispatchRegularFlow() {
	// forwarding table updates
	completedInstructions := ps.ForwardingPaths.GetCompletedInstructions()

	for _, el := range completedInstructions {
		ps.PhysicalRegisterFile[el.tag] = el.value
		ps.BusyBitTable.SetRegisterState(int(el.tag), false)
	}

	// Determine condition for backpressure (if needed)
	spaceFL := ps.FreeList.GetIfSpaceAvailable()
	spaceIQ := ps.IntegerQueue.GetIfSpaceAvailable()
	spaceAL := ps.ActiveList.GetIfSpaceAvailable()
	if !spaceFL || !spaceIQ || !spaceAL {
		ps.DPR.SetBackPressure(true)
		return
	}

	ps.DPR.SetBackPressure(false)

	decodedInstructions := ps.DPR.GetCurrentValue()
	for _, pc := range decodedInstructions {
		instruction := ps.InputInstructions[pc]

		// Determine status operands
		// === Operand A ===
		opAReady := false
		opATag := ps.RegisterMapTable[instruction.GetOpA()]
		var opAValue uint64

		if !ps.BusyBitTable.GetRegisterState(int(opATag)) {
			opAReady = true
			opAValue = ps.PhysicalRegisterFile[int(opATag)]
			opATag = 0
		} else {
			for _, fwd := range ps.ForwardingPaths.GetCompletedInstructions() {
				if fwd.tag == opATag {
					opAReady = true
					opAValue = fwd.value
					opATag = 0
					break
				}
			}
		}

		// === Operand B ===
		var opBReady bool
		var opBTag uint64
		var opBValue uint64
		var opCode string

		if instruction.GetOpCode() == "addi" {
			opBReady = true
			opB := uint64(instruction.GetSecondArg())
			opBValue = opB
			opBTag = 0
			opCode = "add"
		} else {
			opCode = instruction.GetOpCode()
			opBTag = ps.RegisterMapTable[instruction.GetSecondArg()]
			if !ps.BusyBitTable.GetRegisterState(int(opBTag)) {
				opBReady = true
				opBValue = ps.PhysicalRegisterFile[int(opBTag)]
				opBTag = 0
			} else {
				for _, fwd := range ps.ForwardingPaths.GetCompletedInstructions() {
					if fwd.tag == opBTag {
						opBReady = true
						opBValue = fwd.value
						opBTag = 0
						break
					}
				}
			}
		}

		// Renaming destination
		physicalRegister := ps.FreeList.GetRegister()
		oldDestination := ps.RegisterMapTable[instruction.GetDest()]
		ps.RegisterMapTable[instruction.GetDest()] = physicalRegister

		// Allocate in queues

		// Active list
		activeListEntry := ActiveListEntry{
			Done:               false,
			Exception:          false,
			LogicalDestination: uint64(instruction.GetDest()),
			OldDestination:     oldDestination,
			PC:                 pc,
		}
		ps.ActiveList.Append(activeListEntry)

		// Integer queue

		// Reserve destination register
		ps.BusyBitTable.SetRegisterState(int(physicalRegister), true)

		// Create and append IntegerQueueEntry
		integerQueueEntry := IntegerQueueEntry{
			DestRegister: physicalRegister,
			OpAIsReady:   opAReady,
			OpARegTag:    opATag,
			OpAValue:     int(opAValue),
			OpBIsReady:   opBReady,
			OpBRegTag:    opBTag,
			OpBValue:     int(opBValue),
			OpCode:       opCode,
			PC:           pc,
		}
		ps.IntegerQueue.Append(integerQueueEntry)
	}
}

func (ps *ProcessorState) IssueRegularFlow() {
	instructions := ps.IntegerQueue.GetReadyInstructions()

	ps.IssuedInstructionPipelineRegister.SetNextInstructions(instructions)
}

func (ps *ProcessorState) ExecuteRegularFlow() {
	ps.AluPipelineRegisters.SetNextInstructions(
		ps.IssuedInstructionPipelineRegister.GetCurrentInstructions())

	currentExecuteInstructions := ps.AluPipelineRegisters.GetCurrentInstructions()

	var results []ForwardingPathsEntry
	var toCommitInstructions []CommitPipelineRegisterEntry

	for _, iq := range currentExecuteInstructions {
		instruction := ps.InputInstructions[iq.PC]

		value, exception := instruction.Execute(iq.OpAValue, iq.OpBValue)
		res := uint64(value)

		if !exception {
			results = append(results, ForwardingPathsEntry{
				tag:   iq.DestRegister,
				value: res,
			})
		}

		toCommitInstructions = append(toCommitInstructions, CommitPipelineRegisterEntry{
			Done:      true,
			Exception: exception,
			PC:        iq.PC,
		})
	}

	ps.ActiveList.SetDoneBitForInstructions(toCommitInstructions)
	ps.ForwardingPaths.SetCompletedInstruction(results)
	ps.CommitPipeline.SetNextRegister(toCommitInstructions)
	ps.IntegerQueue.ForwardResults(results)
}

func (ps *ProcessorState) CommitRegularFlow() {
	activeList := ps.ActiveList.GetActiveList()

	retired := 0
	for _, instr := range activeList {
		if !instr.Done {
			break
		}

		if instr.Exception {
			ps.Exception = true
			ps.ExceptionPC = instr.PC
			break
		}

		retired++

		ps.FreeList.FreeRegister(instr.OldDestination)

		if retired == 4 {
			break
		}
	}

	ps.ActiveList.RetireInstructions(retired)
}

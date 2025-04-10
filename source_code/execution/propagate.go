package execution

// TODO: Get values from the forwarding paths
// TODO: Implement exceptions
// TODO: In queues, in previous stages i can see what the next stages did (updates): so read from next?

func (ps *ProcessorState) Propagate() error {
	ps.Commit()
	ps.Execute()
	ps.Issue()
	err := ps.RenameAndDispatch()
	if err != nil {
		return err
	}
	ps.FetchAndDecode()
	return nil
}

func (ps *ProcessorState) FetchAndDecode() {
	if ps.Exception.GetCurrentStatus() {
		ps.FetchAndDecodeExceptionFlow()
	} else {
		ps.FetchAndDecodeRegularFlow()
	}
}

func (ps *ProcessorState) RenameAndDispatch() error {

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
		return nil
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

	return nil
}

func (ps *ProcessorState) Issue() {
	instructions := ps.IntegerQueue.GetReadyInstructions()

	ps.IssuedInstructionPipelineRegister.SetNextInstructions(instructions)
}

func (ps *ProcessorState) Execute() {
	ps.AluPipelineRegisters.SetNextInstructions(
		ps.IssuedInstructionPipelineRegister.GetCurrentInstructions())

	currentExecuteInstructions := ps.AluPipelineRegisters.GetCurrentInstructions()

	var results []ForwardingPathsEntry
	var toCommitInstructions []CommitPipelineRegisterEntry

	for _, iq := range currentExecuteInstructions {
		instruction := ps.InputInstructions[iq.PC]

		value, exception := instruction.Execute(iq.OpAValue, iq.OpBValue)
		res := uint64(value)

		results = append(results, ForwardingPathsEntry{
			tag:   iq.DestRegister,
			value: res,
		})

		toCommitInstructions = append(toCommitInstructions, CommitPipelineRegisterEntry{
			Done:      !exception,
			Exception: exception,
			PC:        iq.PC,
		})
	}

	ps.ForwardingPaths.SetCompletedInstruction(results)
	ps.CommitPipeline.SetNextRegister(toCommitInstructions)
	ps.IntegerQueue.ForwardResults(results)
}

func (ps *ProcessorState) Commit() {
	if ps.Exception.GetCurrentStatus() {
		ps.CommitExceptionFlow()
	} else {
		ps.CommitRegularFlow()
	}
}

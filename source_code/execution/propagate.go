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

//func (ps *ProcessorState) RenameAndDispatch() error {
//	// determine condition for backpressure
//	//if len(ps.FreeList) == 0 || len(ps.ActiveList) >= 32 || len(ps.IntegerQueue) >= 32 {
//	//	ps.DPR.SetBackPressure(true)
//	//	return nil
//	//}
//	// else
//
//	decodedInstructions := ps.DPR.GetCurrentValue()
//
//	for el, _ := range decodedInstructions {
//		instruction := ps.InputInstructions[el]
//		// get physical register from the free list
//		physicalRegister := ps.FreeList.GetRegister()
//
//		// TODO: CHEEEEECK
//		// in register map from architectural to physical
//		oldDestination := ps.RegisterMapTable[instruction.GetDest()]
//		ps.RegisterMapTable[instruction.GetDest()] = physicalRegister
//
//		// allocate entry in Active list, integer queue
//		activeListEntry := ActiveListEntry{
//			Done:               false,
//			Exception:          false,
//			LogicalDestination: uint64(instruction.GetDest()),
//			OldDestination:     oldDestination,
//			PC:                 uint64(el),
//		}
//
//		ps.ActiveList.Append(activeListEntry)
//
//		// determine state of operands
//		opAReady := false
//		opATag := ps.RegisterMapTable[instruction.GetOpA()]
//		var opAValue uint64
//
//		if !ps.BusyBitTable.GetRegisterState(int(ps.RegisterMapTable[instruction.GetOpA()])) {
//			opAReady = true
//			ps.BusyBitTable.SetRegisterState(int(ps.RegisterMapTable[instruction.GetOpA()]), true)
//			opAValue = ps.PhysicalRegisterFile[int(ps.RegisterMapTable[instruction.GetOpA()])]
//			opATag = 0
//		}
//
//		var integerQueueEntry IntegerQueueEntry
//
//		if instruction.GetOpCode() == "addi" {
//			integerQueueEntry = IntegerQueueEntry{
//				DestRegister: physicalRegister,
//				OpAIsReady:   opAReady,
//				OpARegTag:    opATag,
//				OpAValue:     int(opAValue),
//				OpBIsReady:   true,
//				OpBRegTag:    opATag,
//				OpBValue:     instruction.GetSecondArg(),
//				OpCode:       instruction.GetOpCode(),
//				PC:           uint64(el),
//			}
//		} else {
//			opBReady := false
//			opBTag := ps.RegisterMapTable[instruction.GetSecondArg()]
//			var opBValue uint64
//
//			if !ps.BusyBitTable.GetRegisterState(int(ps.RegisterMapTable[instruction.GetSecondArg()])) {
//				opBReady = true
//				ps.BusyBitTable.SetRegisterState(int(ps.RegisterMapTable[instruction.GetOpA()]), true)
//				opAValue = ps.PhysicalRegisterFile[int(ps.RegisterMapTable[instruction.GetSecondArg()])]
//				opBTag = 0
//			}
//
//			integerQueueEntry = IntegerQueueEntry{
//				DestRegister: physicalRegister,
//				OpAIsReady:   opAReady,
//				OpARegTag:    opATag,
//				OpAValue:     int(opAValue),
//				OpBIsReady:   opBReady,
//				OpBRegTag:    opBTag,
//				OpBValue:     int(opBValue),
//				OpCode:       instruction.GetOpCode(),
//				PC:           uint64(el),
//			}
//		}
//
//		ps.IntegerQueue.Append(integerQueueEntry)
//	}
//
//	return nil
//}

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

		res := uint64(instruction.Execute(iq.OpAValue, iq.OpBValue))

		results = append(results, ForwardingPathsEntry{
			tag:   iq.DestRegister,
			value: res,
		})

		toCommitInstructions = append(toCommitInstructions, CommitPipelineRegisterEntry{
			Done:      true,
			Exception: false,
			PC:        iq.PC,
		})
	}

	ps.ForwardingPaths.SetCompletedInstruction(results)
	ps.CommitPipeline.SetNextRegister(toCommitInstructions)
	ps.IntegerQueue.ForwardResults(results)
}

func (ps *ProcessorState) Commit() {
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

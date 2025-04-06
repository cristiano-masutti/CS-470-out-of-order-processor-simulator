package execution

func (ps *ProcessorState) Propagate() error {
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
	numberOfInstructions := uint64(min(4, len(ps.InputInstructions)-int(pc)+1))
	ps.PCP.SetNextValue(numberOfInstructions)

	var nextDecodedPCs []uint64

	for i := uint64(0); i < numberOfInstructions; i++ {
		nextDecodedPCs = append(nextDecodedPCs, pc+i)
	}

	ps.DPR.SetNextValue(nextDecodedPCs)
}

func (ps *ProcessorState) RenameAndDispatch() error {
	// determine condition for backpressure
	//if len(ps.FreeList) == 0 || len(ps.ActiveList) >= 32 || len(ps.IntegerQueue) >= 32 {
	//	ps.DPR.SetBackPressure(true)
	//	return nil
	//}
	// else

	decodedInstructions := ps.DPR.GetCurrentValue()

	for el, _ := range decodedInstructions {
		instruction := ps.InputInstructions[el]
		// get physical register from the free list
		physicalRegister := ps.FreeList.GetRegister()

		// in register map from architectural to physical
		oldDestination := ps.RegisterMapTable[instruction.GetDest()]
		ps.RegisterMapTable[instruction.GetDest()] = physicalRegister

		// allocate entry in Active list, integer queue
		activeListEntry := ActiveListEntry{
			Done:               false,
			Exception:          false,
			LogicalDestination: uint64(instruction.GetDest()),
			OldDestination:     oldDestination,
			PC:                 uint64(el),
		}

		ps.ActiveList.Append(activeListEntry)

		// determine state of operands
		opAReady := false
		opATag := ps.RegisterMapTable[instruction.GetOpA()]
		var opAValue uint64

		if !ps.BusyBitTable.GetRegisterState(int(ps.RegisterMapTable[instruction.GetOpA()])) {
			opAReady = true
			ps.BusyBitTable.SetRegisterState(int(ps.RegisterMapTable[instruction.GetOpA()]), true)
			opAValue = ps.PhysicalRegisterFile[int(ps.RegisterMapTable[instruction.GetOpA()])]
			opATag = 0
		}

		var integerQueueEntry IntegerQueueEntry

		if instruction.GetOpCode() == "addi" {
			integerQueueEntry = IntegerQueueEntry{
				DestRegister: physicalRegister,
				OpAIsReady:   opAReady,
				OpARegTag:    opATag,
				OpAValue:     opAValue,
				OpBIsReady:   true,
				OpBRegTag:    opATag,
				OpBValue:     uint64(instruction.GetSecondArg()),
				OpCode:       instruction.GetOpCode(),
				PC:           uint64(el),
			}
		} else {
			opBReady := false
			opBTag := ps.RegisterMapTable[instruction.GetSecondArg()]
			var opBValue uint64

			if !ps.BusyBitTable.GetRegisterState(int(ps.RegisterMapTable[instruction.GetSecondArg()])) {
				opBReady = true
				ps.BusyBitTable.SetRegisterState(int(ps.RegisterMapTable[instruction.GetOpA()]), true)
				opAValue = ps.PhysicalRegisterFile[int(ps.RegisterMapTable[instruction.GetSecondArg()])]
				opBTag = 0
			}

			integerQueueEntry = IntegerQueueEntry{
				DestRegister: physicalRegister,
				OpAIsReady:   opAReady,
				OpARegTag:    opATag,
				OpAValue:     opAValue,
				OpBIsReady:   opBReady,
				OpBRegTag:    opBTag,
				OpBValue:     opBValue,
				OpCode:       instruction.GetOpCode(),
				PC:           uint64(el),
			}
		}

		ps.IntegerQueue.Append(integerQueueEntry)
	}

	return nil
}

func (ps *ProcessorState) Issue() {

}

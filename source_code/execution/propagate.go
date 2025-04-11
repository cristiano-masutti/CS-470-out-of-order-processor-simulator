package execution

// TODO: Get values from the forwarding paths
// TODO: Implement exceptions
// TODO: In queues, in previous stages i can see what the next stages did (updates): so read from next?

func (ps *ProcessorState) Propagate() {
	ps.Commit()
	ps.Execute()
	ps.Issue()
	ps.RenameAndDispatch()
	ps.FetchAndDecode()
}

func (ps *ProcessorState) FetchAndDecode() {
	if ps.Exception.GetCurrentStatus() {
		ps.FetchAndDecodeExceptionFlow()
	} else {
		ps.FetchAndDecodeRegularFlow()
	}
}

func (ps *ProcessorState) RenameAndDispatch() {
	if ps.Exception.GetCurrentStatus() {
		ps.RenameAndDispatchExceptionFlow()
	} else {
		ps.RenameAndDispatchRegularFlow()
	}
}

func (ps *ProcessorState) Issue() {
	if ps.Exception.GetCurrentStatus() {
		ps.IssueExceptionFlow()
	} else {
		ps.IssueRegularFlow()
	}
}

func (ps *ProcessorState) Execute() {
	if ps.Exception.GetCurrentStatus() {
		ps.ExecuteExceptionFlow()
	} else {
		ps.ExecuteRegularFlow()
	}
}

func (ps *ProcessorState) Commit() {
	if ps.Exception.GetCurrentStatus() {
		ps.CommitExceptionFlow()
	} else {
		ps.CommitRegularFlow()
	}
}

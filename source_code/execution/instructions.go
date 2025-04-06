package execution

// Instruction is an interface so to create datastructures
type Instruction interface {
	GetDest() int
	GetOpA() int
	GetOpCode() string
	GetSecondArg() int
	Execute(a, b int) int
}

type BaseInstruction struct {
	Dest int
	OpA  int
}

func (bi *BaseInstruction) GetDest() int { return bi.Dest }
func (bi *BaseInstruction) GetOpA() int  { return bi.OpA }
func (bi *BaseInstruction) GetOpCode() string {
	return "Base Instruction"
}
func (bi *BaseInstruction) GetSecondArg() int {
	return bi.OpA
}
func (bi *BaseInstruction) Execute(a, b int) int {
	return a + b
}

// Actual Instructions types

type Add struct {
	BaseInstruction
	OpB int
}

func (ad *Add) GetOpCode() string    { return "add" }
func (ad *Add) GetSecondArg() int    { return ad.OpB }
func (ad *Add) Execute(a, b int) int { return a + b }

type Addi struct {
	BaseInstruction
	Imm int
}

func (ad *Addi) GetOpCode() string    { return "addi" }
func (ad *Addi) GetSecondArg() int    { return ad.Imm }
func (ad *Addi) Execute(a, b int) int { return a + b }

type Sub struct {
	BaseInstruction
	OpB int
}

func (s *Sub) GetOpCode() string    { return "sub" }
func (s *Sub) GetSecondArg() int    { return s.OpB }
func (s *Sub) Execute(a, b int) int { return a - b }

type Mulu struct {
	BaseInstruction
	OpB int
}

func (m *Mulu) GetOpCode() string    { return "mulu" }
func (m *Mulu) GetSecondArg() int    { return m.OpB }
func (m *Mulu) Execute(a, b int) int { return int(uint(a) * uint(b)) }

type Divu struct {
	BaseInstruction
	OpB int
}

func (d *Divu) GetOpCode() string    { return "divu" }
func (d *Divu) GetSecondArg() int    { return d.OpB }
func (d *Divu) Execute(a, b int) int { return int(uint(a) / uint(b)) }

type Remu struct {
	BaseInstruction
	OpB int
}

func (r *Remu) GetOpCode() string    { return "remu" }
func (r *Remu) GetSecondArg() int    { return r.OpB }
func (r *Remu) Execute(a, b int) int { return int(uint(a) % uint(b)) }

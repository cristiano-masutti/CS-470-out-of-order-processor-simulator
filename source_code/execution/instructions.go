package execution

// Instruction is an interface so to create datastructures
type Instruction interface {
	GetDest() int
	GetOpA() int
	GetOpCode() string
	GetSecondArg() int
}

type BaseInstruction struct {
	Dest int
	OpA  int
}

func (b *BaseInstruction) GetDest() int { return b.Dest }
func (b *BaseInstruction) GetOpA() int  { return b.OpA }
func (b *BaseInstruction) GetOpCode() string {
	return "Base Instruction"
}
func (b *BaseInstruction) GetSecondArg() int {
	return b.OpA
}

// Actual Instructions types

type Add struct {
	BaseInstruction
	OpB int
}

func (a *Add) GetOpCode() string { return "add" }
func (a *Add) GetSecondArg() int { return a.OpB }

type Addi struct {
	BaseInstruction
	Imm int
}

func (a *Addi) GetOpCode() string { return "addi" }
func (a *Addi) GetSecondArg() int { return a.Imm }

type Sub struct {
	BaseInstruction
	OpB int
}

func (s *Sub) GetOpCode() string { return "sub" }
func (s *Sub) GetSecondArg() int { return s.OpB }

type Mulu struct {
	BaseInstruction
	OpB int
}

func (m *Mulu) GetOpCode() string { return "mulu" }
func (m *Mulu) GetSecondArg() int { return m.OpB }

type Divu struct {
	BaseInstruction
	OpB int
}

func (d *Divu) GetOpCode() string { return "divu" }
func (d *Divu) GetSecondArg() int { return d.OpB }

type Remu struct {
	BaseInstruction
	OpB int
}

func (r *Remu) GetOpCode() string { return "remu" }
func (r *Remu) GetSecondArg() int { return r.OpB }

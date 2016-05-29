package evmopt

import (
	"fmt"
    "math/big"
)

type Instruction struct {
    Op OpCode
    Arg *big.Int
    Reaches map[int]bool 		// List of program addresses that rely on the output of this instruction
    ReachedBy []map[int]bool 	// List of program addresses that may provide the value for each operand
}

func (self Instruction) String() string {
	if self.Arg != nil {
		return fmt.Sprintf("%v 0x%x", self.Op, self.Arg)
	} else {
		return self.Op.String()
	}
}

type Program struct {
	Instructions map[int]*Instruction
}

func NewProgram(bytecode []byte) *Program {
	program := &Program{
		Instructions: make(map[int]*Instruction),
	}

	for i := 0; i < len(bytecode); i++ {
		op := OpCode(bytecode[i])
		size := op.OperandSize()
		var arg *big.Int
		if size > 0 {
			arg = big.NewInt(0)
			for j := 1; j <= size; j++ {
				arg.Lsh(arg, 8)
				if i + j < len(bytecode) {
					arg.Or(arg, big.NewInt(int64(bytecode[i + j])))
				}
			}
		}
		program.Instructions[i] = &Instruction{
			Op: op,
			Arg: arg,
			Reaches: make(map[int]bool),
			ReachedBy: make([]map[int]bool, op.StackReads()),
		}
		i += size
	}

	program.buildReachings()

	return program
}

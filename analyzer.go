package evmopt

import (
    "fmt"
    "log"
    "math/big"
    "strings"
)

type ReachingPool []map[int]bool

func (self ReachingPool) Combine(other ReachingPool) ReachingPool {
    size := len(self)
    if len(other) > len(self) {
        size = len(other)
    }

    ret := make(ReachingPool, size)
    for i := 0; i < size; i++ {
        ret[i] = make(map[int]bool)
        if i < len(self) {
            for elt := range self[i] {
                ret[i][elt] = true
            }
        }
        if i < len(other) {
            for elt := range other[i] {
                ret[i][elt] = true
            }
        }
    }

    return ret
}

func (self ReachingPool) Equal(other ReachingPool) bool {
    if len(self) != len(other) {
        return false
    }
    for i := range self {
        a, b := self[i], other[i]
        if len(a) != len(b) {
            return false
        }
        for k := range a {
            if !b[k] {
                return false
            }
        }
    }
    return true
}

func (self ReachingPool) String() string {
    frames := make([]string, len(self))
    for i := 0; i < len(self); i++ {
        frame := make([]int, len(self[i]))
        j := 0
        for val := range self[i] {
            frame[j] = val
            j++
        }
        frames[i] = fmt.Sprintf("%X", frame)
    }
    return strings.Join(frames, " ")
}

func (self ReachingPool) Copy() ReachingPool {
    ret := make(ReachingPool, len(self))
    for i, m := range self {
        ret[i] = make(map[int]bool)
        for k, v := range m {
            ret[i][k] = v
        }
    }
    return ret
}

type Operation struct {
    instruction *Instruction
    source int
}

func (self *Operation) Source() int { return self.source }
func (self *Operation) Value() *big.Int { return self.instruction.Arg }
func (self *Operation) String() string { return self.instruction.String() }

type StackFrame struct {
    Up *StackFrame
    Height int
    Value *Operation
}

func NewFrame(up *StackFrame, value *Operation) *StackFrame {
    if up != nil {
        return &StackFrame{up, up.Height + 1, value}
    } else {
        return &StackFrame{nil, 0, value}
    }
}

func (self *StackFrame) UpBy(num int) *StackFrame {
    ret := self
    for i := 0; i < num; i++ {
        ret = ret.Up
    }
    return ret
}

func (self *StackFrame) Replace(num int, value *Operation) (*StackFrame, *Operation) {
    if num == 0 {
        return NewFrame(self.Up, value), self.Value
    }
    up, old := self.Up.Replace(num - 1, value)
    return NewFrame(up, self.Value), old
}

func (self *StackFrame) Swap(num int) *StackFrame {
    up, old := self.Up.Replace(num - 1, self.Value)
    return NewFrame(up, old)
}

func (self *StackFrame) String() string {
    if self.Up != nil {
        return fmt.Sprintf("%v %v", self.Value, self.Up)
    } else {
        return fmt.Sprintf("%v", self.Value)
    }
}

func (self *StackFrame) Popn(n int) (values []*StackFrame, stack *StackFrame) {
    stack = self
    values = make([]*StackFrame, n)
    for i := 0; i < n; i++ {
        values[i] = stack
        stack = stack.Up
    }
    return values, stack
}

type programState struct {
    pc int
    stack *StackFrame
}

func (self *Program) buildReachings() {
    pools := make(map[int]ReachingPool)
    states := []*programState{
        &programState{0, nil},
    }

    i := 0
    for len(states) > 0 {
        var state *programState
        state, states = states[len(states) - 1], states[:len(states) - 1]
        log.Printf("%v PC: 0x%X, op: %v, stack: %v", i, state.pc, self.Instructions[state.pc], state.stack)
        i += 1
        //log.Printf("PC: 0x%X, op: %v, pool: %v", state.pc, self.Instructions[state.pc], pools[state.pc])
        successors := processInstruction(self, state)

        for _, successor := range successors {
            result := getValue(self, successor, pools[state.pc])
            if pools[successor.pc] == nil {
                pools[successor.pc] = result
                states = append(states, successor)
            } else {
                newPool := result.Combine(pools[successor.pc])
                if !newPool.Equal(pools[successor.pc]) {
                    pools[successor.pc] = newPool
                    states = append(states, successor)
                }
            }
        }
    }

    for pc, instruction := range self.Instructions {
        reachedBy := pools[pc]

        // Build the list of instructions that can be the input for each arg, and vice-versa
        for i := 0; i < instruction.Op.StackReads(); i++ {
            instruction.ReachedBy[i] = make(map[int]bool, len(reachedBy[i]))
            for j := range reachedBy[i] {
                instruction.ReachedBy[i][j] = true
                if !instruction.Op.IsDup() && !instruction.Op.IsSwap() {
                    self.Instructions[j].Reaches[pc] = true
                }
            }
        }
    }
}

func getValue(prog *Program, state *programState, inpool ReachingPool) (pool ReachingPool) {
    for s := state.stack; s != nil; s = s.Up {
        pool = append(pool, map[int]bool{s.Value.Source(): true})
    }
    return pool
}

func processInstruction(prog *Program, state *programState) (nextstates []*programState) {
    inst := prog.Instructions[state.pc]
    op := inst.Op
    stack := state.stack

    operandFrames, stack := stack.Popn(op.StackReads())
    operands := make([]*Operation, len(operandFrames))
    for i, frame := range operandFrames {
        operands[i] = frame.Value
    }

    switch op {
    // Ops that terminate execution
    case STOP: break
    case RETURN: break
    case SELFDESTRUCT: break

    case PUSH1: fallthrough
    case PUSH2: fallthrough
    case PUSH3: fallthrough
    case PUSH4: fallthrough
    case PUSH5: fallthrough
    case PUSH6: fallthrough
    case PUSH7: fallthrough
    case PUSH8: fallthrough
    case PUSH9: fallthrough
    case PUSH10: fallthrough
    case PUSH11: fallthrough
    case PUSH12: fallthrough
    case PUSH13: fallthrough
    case PUSH14: fallthrough
    case PUSH15: fallthrough
    case PUSH16: fallthrough
    case PUSH17: fallthrough
    case PUSH18: fallthrough
    case PUSH19: fallthrough
    case PUSH20: fallthrough
    case PUSH21: fallthrough
    case PUSH22: fallthrough
    case PUSH23: fallthrough
    case PUSH24: fallthrough
    case PUSH25: fallthrough
    case PUSH26: fallthrough
    case PUSH27: fallthrough
    case PUSH28: fallthrough
    case PUSH29: fallthrough
    case PUSH30: fallthrough
    case PUSH31: fallthrough
    case PUSH32:
        nextstates = []*programState{
            &programState{state.pc + op.OperandSize() + 1, NewFrame(stack, &Operation{inst, state.pc})},
        }
    case JUMP:
        if operands[0].Value() == nil {
            log.Fatalf("%v: Could not determine jump location statically; source is %v", state.pc, operands[0].Source())
        }
        nextstates = []*programState{
            &programState{int(operands[0].Value().Int64()), stack},
        }
    case JUMPI:
        if operands[0].Value() == nil {
            log.Fatalf("%v: Could not determine jump location statically; source is %v", state.pc, operands[0].Source())
        }
        nextstates = []*programState{
            &programState{int(operands[0].Value().Int64()), stack},
            &programState{state.pc + 1, stack},
        }
    case DUP1: fallthrough
    case DUP2: fallthrough
    case DUP3: fallthrough
    case DUP4: fallthrough
    case DUP5: fallthrough
    case DUP6: fallthrough
    case DUP7: fallthrough
    case DUP8: fallthrough
    case DUP9: fallthrough
    case DUP10: fallthrough
    case DUP11: fallthrough
    case DUP12: fallthrough
    case DUP13: fallthrough
    case DUP14: fallthrough
    case DUP15: fallthrough
    case DUP16:
        // Uses state.stack instead of stack, because we don't actually want to pop all those elements
        nextstates = []*programState{
            &programState{state.pc + 1, NewFrame(state.stack, state.stack.UpBy(op.StackReads() - 1).Value)},
        }
    case SWAP1: fallthrough
    case SWAP2: fallthrough
    case SWAP3: fallthrough
    case SWAP4: fallthrough
    case SWAP5: fallthrough
    case SWAP6: fallthrough
    case SWAP7: fallthrough
    case SWAP8: fallthrough
    case SWAP9: fallthrough
    case SWAP10: fallthrough
    case SWAP11: fallthrough
    case SWAP12: fallthrough
    case SWAP13: fallthrough
    case SWAP14: fallthrough
    case SWAP15: fallthrough
    case SWAP16:
        // Uses state.stack instead of stack, because we don't actually want to pop all those elements
        nextstates = []*programState{
            &programState{state.pc + 1, state.stack.Swap(op.StackReads() - 1)},
        }
    default:
        switch op.StackWrites() {
        case 0:
            nextstates = []*programState{
                &programState{state.pc + 1, stack},
            }
        case 1:
            nextstates = []*programState{
                &programState{state.pc + 1, NewFrame(stack, &Operation{inst, state.pc})},
            }
        default:
            log.Fatalf("Unexpected op %v makes %v writes to the stack", op, op.StackWrites())
        }
    }

    if stack != nil && stack.Height > 1024 {
        // Stack too tall; no future states.
        return nil
    }

    return nextstates
}

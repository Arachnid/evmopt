package evmopt

import (
    "fmt"
    "log"
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
        frames[i] = fmt.Sprintf("%v", frame)
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

type stackFrame struct {
    up *stackFrame
    source int
    height int
}

func NewFrame(up *stackFrame, source int) *stackFrame {
    if up != nil {
        return &stackFrame{up, source, up.height + 1}
    } else {
        return &stackFrame{nil, source, 0}
    }
}

func (self *stackFrame) upBy(num int) *stackFrame {
    ret := self
    for i := 0; i < num; i++ {
        ret = ret.up
    }
    return ret
}

func (self *stackFrame) replace(num int, value int) (*stackFrame, int) {
    if num == 0 {
        return NewFrame(self.up, value), self.source
    }
    up, old := self.up.replace(num - 1, value)
    return NewFrame(up, self.source), old
}

func (self *stackFrame) swap(num int) *stackFrame {
    up, old := self.up.replace(num - 1, self.source)
    return NewFrame(up, old)
}

func (self *stackFrame) String() string {
    if self.up != nil {
        return fmt.Sprintf("%v %v", self.source, self.up)
    } else {
        return fmt.Sprintf("%v", self.source)
    }
}

type programState struct {
    pc int
    stack *stackFrame
}

func Analyze(p *Program) map[int]ReachingPool {
    pools := make(map[int]ReachingPool)
    states := []*programState{
        &programState{0, nil},
    }

    for len(states) > 0 {
        var state *programState
        state, states = states[len(states) - 1], states[:len(states) - 1]
        //log.Printf("PC: %v, op: %v, stack: %v", state.pc, p.Instructions[state.pc], state.stack)
        result, successors := processInstruction(state, p)

        for _, successor := range successors {
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

    return pools
}

func processInstruction(state *programState, prog *Program) (pool ReachingPool, nextstates []*programState) {
    inst := prog.Instructions[state.pc]
    op := inst.Op
    stack := state.stack
    var nextpcs []int

    // PC changes
    switch op {
    case JUMPI:
        sourceInst := prog.Instructions[stack.source]
        if !sourceInst.Op.IsPush() {
            log.Fatalf("Found jump with operand type %v at %v", sourceInst.Op, stack.source)  
        }
        nextpcs = []int{state.pc + 1, int(sourceInst.Arg.Int64())}
    case JUMP:
        sourceInst := prog.Instructions[state.stack.source]
        if !sourceInst.Op.IsPush() {
            log.Fatalf("Found jump with operand type %v at %v", sourceInst.Op, state.stack.source)  
        }
        nextpcs = []int{int(sourceInst.Arg.Int64())}
    case RETURN: break
    case STOP: break
    case SELFDESTRUCT: break
    default:
        nextpcs = []int{state.pc + op.OperandSize() + 1}
    }

    // Stack changes
    switch op {
    case DUP1: stack = NewFrame(stack, stack.source)
    case DUP2: stack = NewFrame(stack, stack.up.source)
    case DUP3: stack = NewFrame(stack, stack.upBy(2).source)
    case DUP4: stack = NewFrame(stack, stack.upBy(3).source)
    case DUP5: stack = NewFrame(stack, stack.upBy(4).source)
    case DUP6: stack = NewFrame(stack, stack.upBy(5).source)
    case DUP7: stack = NewFrame(stack, stack.upBy(6).source)
    case DUP8: stack = NewFrame(stack, stack.upBy(7).source)
    case DUP9: stack = NewFrame(stack, stack.upBy(8).source)
    case DUP10: stack = NewFrame(stack, stack.upBy(9).source)
    case DUP11: stack = NewFrame(stack, stack.upBy(10).source)
    case DUP12: stack = NewFrame(stack, stack.upBy(11).source)
    case DUP13: stack = NewFrame(stack, stack.upBy(12).source)
    case DUP14: stack = NewFrame(stack, stack.upBy(13).source)
    case DUP15: stack = NewFrame(stack, stack.upBy(14).source)
    case DUP16: stack = NewFrame(stack, stack.upBy(15).source)
    case SWAP1: stack = stack.swap(1)
    case SWAP2: stack = stack.swap(2)
    case SWAP3: stack = stack.swap(3)
    case SWAP4: stack = stack.swap(4)
    case SWAP5: stack = stack.swap(5)
    case SWAP6: stack = stack.swap(6)
    case SWAP7: stack = stack.swap(7)
    case SWAP8: stack = stack.swap(8)
    case SWAP9: stack = stack.swap(9)
    case SWAP10: stack = stack.swap(10)
    case SWAP11: stack = stack.swap(11)
    case SWAP12: stack = stack.swap(12)
    case SWAP13: stack = stack.swap(13)
    case SWAP14: stack = stack.swap(14)
    case SWAP15: stack = stack.swap(15)
    case SWAP16: stack = stack.swap(16)
    default:
        stack = stack.upBy(op.StackReads())
        for i := 0; i < op.StackWrites(); i++ {
            stack = NewFrame(stack, state.pc)
        }
    }

    if stack != nil && stack.height > 1024 {
        // Stack too tall; no future states.
        return nil, nil
    }

    nextstates = make([]*programState, len(nextpcs))
    for i, pc := range nextpcs {
        nextstates[i] = &programState{pc, stack}
    }

    for s := stack; s != nil; s = s.up {
        pool = append(pool, map[int]bool{s.source: true})
    }    

    return pool, nextstates
}

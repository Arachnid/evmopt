package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "github.com/arachnid/evmopt"
)

/*func markLive(pc int, program *evmopt.Program, reachings map[int]evmopt.ReachingPool, live map[int]bool) {
    if live[pc] {
        return
    }
    live[pc] = true
    for i := 0; i < program.Instructions[pc].Op.StackReads(); i++ {
        for addr := range reachings[pc][i] {
            markLive(addr, program, reachings, live)
        }
    }
}

func findLive(program *evmopt.Program, reachings map[int]evmopt.ReachingPool) (live map[int]bool) {
    live = make(map[int]bool)
    for idx := 0; ; idx += program.Instructions[idx].Op.OperandSize() + 1 {
        inst, ok := program.Instructions[idx]
        if !ok {
            break
        }

        switch inst.Op {
        case evmopt.CALLDATACOPY: fallthrough
        case evmopt.CODECOPY: fallthrough
        case evmopt.EXTCODECOPY: fallthrough
        case evmopt.MSTORE: fallthrough
        case evmopt.MSTORE8: fallthrough
        case evmopt.SSTORE: fallthrough
        case evmopt.JUMP: fallthrough
        case evmopt.JUMPI: fallthrough
        case evmopt.LOG0: fallthrough
        case evmopt.LOG1: fallthrough
        case evmopt.LOG2: fallthrough
        case evmopt.LOG3: fallthrough
        case evmopt.LOG4: fallthrough
        case evmopt.CREATE: fallthrough
        case evmopt.CALL: fallthrough
        case evmopt.RETURN: fallthrough
        case evmopt.CALLCODE: fallthrough
        case evmopt.DELEGATECALL: fallthrough
        case evmopt.STOP: fallthrough
        case evmopt.JUMPDEST: fallthrough
        case evmopt.SELFDESTRUCT:
            markLive(idx, program, reachings, live)
        }
    }

    return live
}*/

func fetchInstructions(program *evmopt.Program, locations map[int]bool) (ret []*evmopt.Instruction) {
    for source := range locations {
        ret = append(ret, program.Instructions[source])
    }
    return ret
}

func intMapKeys(m map[int]bool) []int {
    ret := make([]int, 0, len(m))
    for k := range m {
        ret = append(ret, k)
    }
    return ret
}

func main() {
    bytecode, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        log.Fatalf("Could not read from stdin: %v", err)
    }

    program := evmopt.NewProgram(bytecode)
    //reachings := evmopt.Analyze(program)
    //live := findLive(program, reachings)
    for idx := 0; ; idx += program.Instructions[idx].Op.OperandSize() + 1 {
        inst, ok := program.Instructions[idx]
        if !ok {
            break
        }
        operands := make([][]*evmopt.Instruction, len(inst.ReachedBy))
        for i, frame := range inst.ReachedBy {
            operands[i] = fetchInstructions(program, frame)
        }
        fmt.Printf("0x%X\t%x\t%v\t%v\n", idx, intMapKeys(inst.Reaches), inst, operands)
        //fmt.Printf("0x%X\t%v\t%v\n", idx, live[idx], inst)
    }
}

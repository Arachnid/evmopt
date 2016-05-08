package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "os"

    "github.com/arachnid/evmopt"
)

func main() {
    bytecode, err := ioutil.ReadAll(os.Stdin)
    if err != nil {
        log.Fatalf("Could not read from stdin: %v", err)
    }

    program := evmopt.NewProgram(bytecode)
    reachings := evmopt.Analyze(program)
    for idx := 0; ; idx += program.Instructions[idx].Op.OperandSize() + 1 {
        inst, ok := program.Instructions[idx]
        if !ok {
            break
        }
        stack := reachings[idx][:inst.Op.StackReads()]
        stackInstructions := make([][]evmopt.Instruction, len(stack))
        for i, frame := range stack {
            var frameInstructions []evmopt.Instruction
            for source := range frame {
                frameInstructions = append(frameInstructions, program.Instructions[source])
            }
            stackInstructions[i] = frameInstructions
        }
        fmt.Printf("%v\t%v\t%v\n", idx, inst, stackInstructions)
    }
}

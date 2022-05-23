package brainfuck

import (
	"github.com/pkg/errors"
	"io"
	"unicode"
)

type UnknownCommandCharPolicy int

const (
	IgnoreWhitespacesPolicy UnknownCommandCharPolicy = iota
	IgnoreUnknownCharsPolicy
	ZeroTolerancePolicy
)

type InterpreterRunner struct {
	state      InterpreterState
	charPolicy UnknownCommandCharPolicy
	isFinished bool
}

func NewInterpreterRunner(
	input io.ByteReader,
	output io.ByteWriter,
	opts ...InterpreterRunnerOption,
) (*InterpreterRunner, error) {
	params := interpreterRunnerOptParams{}
	for _, opt := range opts {
		opt(&params)
	}
	if params.Memory == nil {
		params.Memory = &FixedSizeMemory{}
	}
	if params.Stack == nil {
		params.Stack = &SliceBasedStack{}
	}
	if params.CmdMap == nil {
		params.CmdMap = GetDefaultCommandMap()
	}
	coreStackItem := &StdStackItem{
		CmdMap: params.CmdMap,
	}
	err := params.Stack.Push(coreStackItem)
	if err != nil {
		return nil, WrapError(err)
	}
	return &InterpreterRunner{
		state: InterpreterState{
			Memory: params.Memory,
			Stack:  params.Stack,
			Input:  input,
			Output: output,
		},
		charPolicy: params.UnknownCharPolicy,
		isFinished: false,
	}, nil
}

func (runner *InterpreterRunner) Execute(char CommandChar) error {
	if runner.isFinished {
		return EOR
	}
	execStackIt, ok := runner.state.Stack.Top()
	if !ok {
		return errors.Errorf("stack is empty, cannot execute command %c", char)
	}
	err := execStackIt.Execute(runner.state, char)
	if err != nil {
		if errors.Is(err, EOR) {
			runner.isFinished = true
			return EOR
		}
		if runner.charPolicy != ZeroTolerancePolicy {
			var unknownCharErr UnknownCharErr
			if errors.As(err, &unknownCharErr) {
				isSpacesOk := runner.charPolicy == IgnoreWhitespacesPolicy && unicode.IsSpace(rune(unknownCharErr.char))
				isAnyCharOk := runner.charPolicy == IgnoreUnknownCharsPolicy
				if isSpacesOk || isAnyCharOk {
					return nil
				}
			}
		}
	}
	return WrapError(err)
}

type interpreterRunnerOptParams struct {
	Memory            Memory
	Stack             CommandStack
	CmdMap            map[CommandChar]Command
	UnknownCharPolicy UnknownCommandCharPolicy
}

type InterpreterRunnerOption func(runner *interpreterRunnerOptParams)

func WithMemory(memory Memory) InterpreterRunnerOption {
	return func(params *interpreterRunnerOptParams) {
		params.Memory = memory
	}
}

func WithCommandStack(stack CommandStack) InterpreterRunnerOption {
	return func(params *interpreterRunnerOptParams) {
		params.Stack = stack
	}
}

func WithCommands(cmdMap map[CommandChar]Command) InterpreterRunnerOption {
	return func(params *interpreterRunnerOptParams) {
		params.CmdMap = cmdMap
	}
}

func WithUnknownCharPolicy(policy UnknownCommandCharPolicy) InterpreterRunnerOption {
	return func(params *interpreterRunnerOptParams) {
		params.UnknownCharPolicy = policy
	}
}

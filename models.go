package brainfuck

import "io"

type CommandChar rune
type Cell uint16
type Position int64

const (
	MaxCellVal Cell = 65535
	MinCellVal      = 0
)

type Command func(state InterpreterState) error

type Interpreter interface {
	Execute(char CommandChar) error
}

type Memory interface {
	GetHeadPosition() Position
	SetHeadPosition(position Position) error
	Read() Cell
	Write(value Cell)
}

type CommandStack interface {
	Len() int
	Top() (CommandStackItem, bool)
	Push(CommandStackItem) error
	Pop() (CommandStackItem, error)
}

type CommandStackItem interface {
	CopyCommands() map[CommandChar]Command
	Execute(state InterpreterState, char CommandChar) error
}

type InterpreterState struct {
	Memory Memory
	Stack  CommandStack
	Input  io.ByteReader
	Output io.ByteWriter
}

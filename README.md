# go-brianfuck
Package for [brainfuck](https://en.wikipedia.org/wiki/Brainfuck) interpreter.

`go get github.com/aga3000/go-brianfuck`

Other than the interpreter itself there are 3 main entities:
* memory
* command stack
* command

## Memory
Memory is just an array of cells on which the interpreter operates. This implementation uses 8-bit cells with fixed-size array (by default).
Array size is not restricted. You can change it by implementation of this interface:
```go
type Memory interface {
	GetHeadPosition() Position
	SetHeadPosition(position Position) error
	Read() Cell
	Write(value Cell)
}
```
The cell size is fixed and cannot be changed.

## Command
Command is a function that uses memory, command stack and input/output interfaces.
By default, there are [8 basic commands](https://en.wikipedia.org/wiki/Brainfuck#Commands), but you can implement your own.
```go
type Command func(state InterpreterState) error

type InterpreterState struct {
    Memory Memory
    Stack  CommandStack
    Input  io.ByteReader
    Output io.ByteWriter
}
```

## Command Stack
Command stack is a solution to redeclare commands on demand.
In this implementation it's used for jump commands `[`, `]`.
But it could be useful for many other commands as well.
For example for a command that saves current cell and place it somewhere else in memory afterwards.

Command stack interface:
```go
type CommandStack interface {
    Len() int
    Top() (CommandStackItem, bool)
    Push(CommandStackItem) error
    Pop() (CommandStackItem, error)
}

type CommandStackItem interface {
    // to copy commands while creating a new stack item
    CopyCommands() map[CommandChar]Command
    // interpreter uses this method of item on top of stack to execute commands
    Execute(state InterpreterState, char CommandChar) error
}
```

## Interpreter
Interpreter doesn't parse source code and reads only one command at a time.
```go
type Interpreter interface {
    Execute(char CommandChar) error
}
```
All of mentioned entities can be passed into the interpreter (Memory, Command Stack and set of Commands).
You also can set policy for unknown chars.
```go
runner, err := brainfuck.NewInterpreterRunner(
    reader, writer,
    brainfuck.WithMemory(myMemoryImplementation),
    brainfuck.WithCommands(myCommandMap),
    brainfuck.WithCommandStack(myCommandStack),
    brainfuck.WithUnknownCharPolicy(brainfuck.IgnoreUnknownCharsPolicy),
)
```

Enjoy, I guess ¯\_(ツ)_/¯
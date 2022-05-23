package brainfuck

import (
	"github.com/pkg/errors"
	"io"
)

const (
	IncrementPosCmdChar CommandChar = '>'
	DecrementPosCmdChar             = '<'
	IncrementValCmdChar             = '+'
	DecrementValCmdChar             = '-'
	GetValCmdChar                   = '.'
	SetValCmdChar                   = ','
	JumpForwardCmdChar              = '['
	JumpBackwardCmdChar             = ']'
)

func GetDefaultCommandMap() map[CommandChar]Command {
	res := map[CommandChar]Command{
		IncrementPosCmdChar: IncrementPos,
		DecrementPosCmdChar: DecrementPos,
		IncrementValCmdChar: IncrementVal,
		DecrementValCmdChar: DecrementVal,
		SetValCmdChar:       SetVal,
		GetValCmdChar:       GetVal,
	}
	jumpFuncRepo := JumpFuncRepo{
		JumpForwardCmdChar:  JumpForwardCmdChar,
		JumpBackwardCmdChar: JumpBackwardCmdChar,
	}
	_ = jumpFuncRepo.AddJumpFunctions(res) // all commands' chars were determined beforehand. It should be all right
	return res
}

func IncrementPos(state InterpreterState) error {
	curPos := state.Memory.GetHeadPosition()
	return state.Memory.SetHeadPosition(curPos + 1)

}

func DecrementPos(state InterpreterState) error {
	curPos := state.Memory.GetHeadPosition()
	return state.Memory.SetHeadPosition(curPos - 1)
}

func IncrementVal(state InterpreterState) error {
	mem := state.Memory
	cellVal := mem.Read()
	mem.Write(cellVal + 1)
	return nil
}

func DecrementVal(state InterpreterState) error {
	mem := state.Memory
	cellVal := mem.Read()
	mem.Write(cellVal - 1)
	return nil
}

func GetVal(state InterpreterState) error {
	cellVal := state.Memory.Read()
	err := state.Output.WriteByte(byte(cellVal))
	return WrapError(err)
}

func SetVal(state InterpreterState) error {
	cellByte, err := state.Input.ReadByte()
	if errors.Is(err, io.EOF) {
		return nil
	}
	if err != nil {
		return WrapError(err)
	}
	state.Memory.Write(Cell(cellByte))
	return nil
}

type JumpFuncRepo struct {
	JumpForwardCmdChar  CommandChar
	JumpBackwardCmdChar CommandChar
}

func (r JumpFuncRepo) AddJumpFunctions(cmdMap map[CommandChar]Command) error {
	if _, ok := cmdMap[r.JumpForwardCmdChar]; ok {
		return errors.Errorf(
			"cannot add JumpForward cmd on char %c. It's already occupied with another cmd",
			r.JumpForwardCmdChar,
		)
	}
	if _, ok := cmdMap[r.JumpForwardCmdChar]; ok {
		return errors.Errorf(
			"JumpBackward cmd will be add on char %c. Its place must not be occupied with another cmd",
			r.JumpBackwardCmdChar,
		)
	}
	cmdMap[r.JumpForwardCmdChar] = r.JumpForward
	return nil
}

func (r JumpFuncRepo) JumpForward(state InterpreterState) error {
	err := r.jumpForwardHelper(state, false)
	return WrapError(err)
}

func (r JumpFuncRepo) JumpForwardAgain(state InterpreterState) error {
	err := r.jumpForwardHelper(state, true)
	return WrapError(err)
}

func (r JumpFuncRepo) jumpForwardHelper(state InterpreterState, useParentStackItem bool) error {
	parentStackItem, exists := state.Stack.Top()
	if !exists {
		return errors.Errorf("stack is empty, cannot create a new stack item")
	}
	execCmdMap := parentStackItem.CopyCommands()
	cmdWrapper := SaveCmdWrapper{}
	if state.Memory.Read() == MinCellVal {
		countWrapper := CountCmdWrapper{Counter: 1}
		for ch := range execCmdMap {
			execCmdMap[ch] = DoNothing
		}
		execCmdMap[r.JumpForwardCmdChar] = countWrapper.IncreaseCounterCmd
		execCmdMap[r.JumpBackwardCmdChar] = countWrapper.DecreaseCounterCmd
		if !useParentStackItem {
			newStackItem := StdStackItem{CmdMap: execCmdMap}
			err := state.Stack.Push(newStackItem)
			return WrapError(err)
		}
	} else {
		execCmdMap[r.JumpForwardCmdChar] = r.JumpForwardAgain
		execCmdMap[r.JumpBackwardCmdChar] = cmdWrapper.ExecuteSavedCommands
	}
	saveCmdMap := make(map[CommandChar]Command, len(execCmdMap))
	var parentStackItemCopy CommandStackItem
	if useParentStackItem {
		parentStackItemCopy = parentStackItem
	}
	for ch := range execCmdMap {
		saveCmdMap[ch] = cmdWrapper.GenerateSaveCmdFunc(ch, parentStackItemCopy)
	}
	stackItem := SavedCommandsStackItem{
		Id:            state.Stack.Len(),
		ExecCommands:  StdStackItem{CmdMap: execCmdMap},
		SaveCommands:  StdStackItem{CmdMap: saveCmdMap},
		savedCommands: &cmdWrapper.SavedCommands,
	}
	err := state.Stack.Push(stackItem)
	return WrapError(err)
}

type SaveCmdWrapper struct {
	SavedCommands   []CommandChar
	isSavingStopped bool
	stopSavingMinId *int
}

func (w *SaveCmdWrapper) GenerateSaveCmdFunc(cmdChar CommandChar, parentStackItem CommandStackItem) Command {
	return func(state InterpreterState) error {
		wasSaved := w.saveCmdChar(cmdChar)
		if wasSaved && parentStackItem != nil {
			err := parentStackItem.Execute(state, cmdChar)
			return WrapError(err)
		}
		return nil
	}
}

func (w *SaveCmdWrapper) ExecuteSavedCommands(state InterpreterState) error {
	w.isSavingStopped = true
	w.SavedCommands = w.SavedCommands[:len(w.SavedCommands)-1] // ignore last JumpBackwardCmdChar
	for {
		mem := state.Memory
		cellVal := mem.Read()
		if cellVal == MinCellVal {
			break
		}
		for _, cmdChar := range w.SavedCommands {
			execStackItem, ok := state.Stack.Top()
			if !ok {
				return errors.Errorf("cannot execute postponed commands %v, because stack is empty", w.SavedCommands)
			}
			err := execStackItem.Execute(state, cmdChar)
			if err != nil {
				return WrapError(err)
			}
		}
	}
	_, err := state.Stack.Pop()
	return WrapError(err)
}

func (w *SaveCmdWrapper) saveCmdChar(char CommandChar) bool {
	if !w.isSavingStopped {
		w.SavedCommands = append(w.SavedCommands, char)
		return true
	}
	return false
}

type CountCmdWrapper struct {
	Counter int
}

func (w *CountCmdWrapper) IncreaseCounterCmd(InterpreterState) error {
	w.Counter++
	return nil
}

func (w *CountCmdWrapper) DecreaseCounterCmd(state InterpreterState) error {
	w.Counter--
	if w.Counter == 0 {
		_, err := state.Stack.Pop()
		return WrapError(err)
	}
	return nil
}

func DoNothing(InterpreterState) error {
	return nil
}

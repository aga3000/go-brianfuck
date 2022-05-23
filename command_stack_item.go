package brainfuck

import "github.com/pkg/errors"

type StdStackItem struct {
	CmdMap map[CommandChar]Command
}

func (it StdStackItem) CopyCommands() map[CommandChar]Command {
	return copyCmdMap(it.CmdMap)
}

func (it StdStackItem) Execute(state InterpreterState, ch CommandChar) error {
	cmd, ok := it.CmdMap[ch]
	if !ok {
		return NewUnknownCharErr(ch)
	}
	err := cmd(state)
	return WrapError(err)
}

func copyCmdMap(cmdMap map[CommandChar]Command) map[CommandChar]Command {
	copiedCmdMap := make(map[CommandChar]Command, len(cmdMap))
	for ch, cmd := range cmdMap {
		copiedCmdMap[ch] = cmd
	}
	return copiedCmdMap
}

type SavedCommandsStackItem struct {
	Id            int
	ExecCommands  StdStackItem
	SaveCommands  StdStackItem
	savedCommands *[]CommandChar // it's easier to debug that way, sry ^^
}

func (it SavedCommandsStackItem) CopyCommands() map[CommandChar]Command {
	return it.ExecCommands.CopyCommands()
}

func (it SavedCommandsStackItem) Execute(state InterpreterState, ch CommandChar) error {
	topIt, exists := state.Stack.Top()
	if !exists {
		return errors.Errorf("cannot execute with empty stack")
	}
	isTop := false
	switch topIt.(type) {
	case SavedCommandsStackItem:
		if topIt.(SavedCommandsStackItem).Id == it.Id {
			isTop = true
		}
	case *SavedCommandsStackItem:
		if topIt.(*SavedCommandsStackItem).Id == it.Id {
			isTop = true
		}
	}
	err := it.SaveCommands.Execute(state, ch)
	if err == nil && isTop {
		err = it.ExecCommands.Execute(state, ch)
	}
	return WrapError(err)
}

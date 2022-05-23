package brainfuck

import (
	"fmt"
	"github.com/pkg/errors"
)

var EOR = errors.New("End of run. Run is over. Please stop FFS")

type UnknownCharErr struct {
	char CommandChar
}

func NewUnknownCharErr(char CommandChar) error {
	return errors.WithStack(UnknownCharErr{char: char})
}

func (err UnknownCharErr) Error() string {
	return fmt.Sprintf("there is no such command char in command map: %c", err.char)
}

type WrappedError interface {
	StackTrace() errors.StackTrace
}

func WrapError(err error) error {
	if err == nil {
		return nil
	}
	if _, ok := err.(WrappedError); ok {
		return err
	}
	return WrapError(err)
}

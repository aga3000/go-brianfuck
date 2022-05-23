package brainfuck

import (
	"github.com/pkg/errors"
)

const fixedMemSize = 30000

type FixedSizeMemory struct {
	cells [fixedMemSize]Cell
	head  Position
}

func (mem *FixedSizeMemory) Read() Cell                { return mem.cells[mem.head] }
func (mem *FixedSizeMemory) Write(value Cell)          { mem.cells[mem.head] = value }
func (mem *FixedSizeMemory) GetHeadPosition() Position { return mem.head }

func (mem *FixedSizeMemory) SetHeadPosition(position Position) error {
	if position < 0 || position >= fixedMemSize {
		return errors.Errorf("head position %d is out of memory bounds [%d, %d]", position, 0, fixedMemSize-1)
	}
	mem.head = position
	return nil
}

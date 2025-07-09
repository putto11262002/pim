package pim

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuffer_NewBuffer(t *testing.T) {
	b := NewBuffer()
	assertBufferState(t, defaultBufferState(), *b)
}

func TestBuffer_InsertRune(t *testing.T) {
	b := NewBuffer()

	b.InsertRune('r')
	assertBufferState(t, bufferState{
		cursorLine: 0,
		cursorCol:  1,
		len:        1,
		buffer:     []rune{'r'},
		lineIdx:    []int{0},
	}, *b)
}

func TestBuffer_NewLine(t *testing.T) {
	b := NewBuffer()

	b.NewLine()
	assertBufferState(t, bufferState{
		cursorLine: 1,
		cursorCol:  0,
		len:        1,
		buffer:     []rune{'\n'},
		lineIdx:    []int{0, 1},
	}, *b)
}

func TestBuffer_DeleteRune(t *testing.T) {
	b := NewBuffer()
	b.InsertRune('r')

	b.DeleteRune()
	assertBufferState(t, bufferState{
		cursorLine: 0,
		cursorCol:  0,
		len:        0,
		buffer:     []rune{},
		lineIdx:    []int{0},
	}, *b)

	b.NewLine()
	assertBufferState(t, bufferState{
		cursorLine: 1,
		cursorCol:  0,
		len:        1,
		buffer:     []rune{'\n'},
		lineIdx:    []int{0, 1},
	}, *b)

	b.DeleteRune()
	assertBufferState(t, bufferState{
		cursorLine: 0,
		cursorCol:  0,
		len:        0,
		buffer:     []rune{},
		lineIdx:    []int{0},
	}, *b)
}

type bufferState struct {
	cursorLine int
	cursorCol  int
	len        int
	buffer     []rune
	lineIdx    []int
}

func defaultBufferState() bufferState {
	return bufferState{
		cursorLine: 0,
		cursorCol:  0,
		len:        0,
		buffer:     []rune{},
		lineIdx:    []int{0},
	}
}

func assertBufferState(t *testing.T, exp bufferState, act Buffer) {
	assert.Equal(t, exp.cursorLine, act.cursorLine, "Cursor Line mismatch")
	assert.Equal(t, exp.cursorCol, act.cursorCol, "Cursor Col mismatch")
	assert.Equal(t, exp.len, act.len, "Buffer length mismatch")
	assert.Equal(t, exp.buffer, act.buffer[:act.len], "Buffer content mismatch")
	assert.Equal(t, exp.lineIdx, act.lineIndex, "Line index mismatch")
}

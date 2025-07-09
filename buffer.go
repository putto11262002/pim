package pim

import (
	"fmt"
	"iter"
)

var (
	lineEndingRune = rune('\n')
)

const (
	initialBufferSize = 1024 // Initial size of the buffer
)

type Buffer struct {
	// buffer is a slice of runes that stores the entire text content of the editor,
	// including newline characters. Using `rune` (Go's representation for a Unicode code point)
	// ensures correct handling of multi-byte UTF-8 characters, preventing issues where
	// editing operations might truncate or corrupt characters. This means indexing into
	// 'buffer' is by Unicode character (rune) count, not raw byte count.
	//
	// This buffer is designed with an "extra space at the end" strategy. It may have
	// a capacity greater than its current 'size', allowing for efficient insertions
	// and deletions by shifting runes within the existing allocated memory,
	// minimizing the need for frequent reallocations until the capacity is exhausted.
	buffer []rune
	// lineIndex is a slice of integers where each element stores the starting index
	// (in terms of rune count from the beginning of the 'buffer') of a corresponding line.
	// Each index represents the position where a new character would be inserted
	// if typing at the very beginning of that line.
	//
	// Important Implementation Details:
	// - The index for the first line (Line 0) *is* stored at lineIndex[0], which will always be 0.
	//   This simplifies line lookups (e.g., line N's start is lineIndex[N]).
	// - The length of 'lineIndex' will be equal to the total number of lines in the buffer.
	lineIndex []int
	// cursorLine represents the current line number within the buffer.
	// It is 0-indexed.
	cursorLine int
	// cursorCol represents the current character (rune) offset from the beginning of the current line.
	// It is 0-indexed.
	cursorCol int
	// size is the current number of runes (characters) present in the 'buffer'.
	// This explicitly tracks the logical length of the text, distinct from the
	// underlying 'buffer' slice's capacity.
	len int
}

func NewBuffer() *Buffer {
	return &Buffer{
		buffer:    make([]rune, 1024),
		lineIndex: []int{0},
	}
}

func (buffer *Buffer) lookupRCIdx(row, col int) int {
	return buffer.lookupLineIdx(row) + col
}

func (b *Buffer) cursorIdx() int {
	return b.lookupRCIdx(b.cursorLine, b.cursorCol)
}

func (buffer *Buffer) numLines() int {
	return len(buffer.lineIndex)
}

func (buffer *Buffer) lookupLineIdx(n int) int {
	if n < 0 || n >= len(buffer.lineIndex) {
		panic(fmt.Sprintf("index out of range: %d", n))
	}

	return buffer.lineIndex[n]
}

// insertRune insert the character to the current cursor positions.
// All the characters to the right of the original cursor position are shifted to the right by one position.
// insertRune only modify the state of the buffer, it does not update the cursor position or the line index.
func (buffer *Buffer) insertRune(char rune) {
	// Allocate more space if needed
	if len(buffer.buffer) < buffer.len+1 {
		buffer.buffer = make([]rune, len(buffer.buffer)*2)
		copy(buffer.buffer, buffer.buffer[:buffer.len])
	}

	// 1. Shift over the characters to the right of the cursor
	// 2. Update the line index
	copy(buffer.buffer[buffer.cursorIdx()+1:], buffer.buffer[buffer.cursorIdx():])
	buffer.buffer[buffer.cursorIdx()] = char
	buffer.len++
}

func (buffer *Buffer) InsertRune(char rune) {
	buffer.insertRune(char)

	// Update the line index
	for i := buffer.cursorLine + 1; i < len(buffer.lineIndex); i++ {
		buffer.lineIndex[i]++
	}

	// Update the cursor position
	buffer.cursorCol++
}

func (b *Buffer) lineLen(n int) int {
	if n < 0 || n >= len(b.lineIndex) {
		panic(fmt.Sprintf("line index out of range: %d", n))
	}

	if n < len(b.lineIndex)-1 {
		return b.lineIndex[n+1] - b.lineIndex[n]
	}

	return b.len - b.lineIndex[n]
}

func (b *Buffer) DeleteRune() bool {

	if b.cursorIdx() < 1 {
		return false
	}

	// 1. Shift characters to the right of the cursor by 1 position
	// 2. Decrement line indexes of the lines after the cursor by one
	b.buffer[b.cursorIdx()] = 0 // Clear the character at the cursor position
	copy(b.buffer[b.cursorIdx()-1:], b.buffer[b.cursorIdx():])
	b.len--

	if b.cursorCol > 0 {
		b.cursorCol--
		for i := b.cursorLine + 1; i < len(b.lineIndex); i++ {
			b.lineIndex[i]--
		}
	} else if b.cursorLine > 0 {
		// Delete the deleted line from the line index
		copy(b.lineIndex[b.cursorLine:], b.lineIndex[b.cursorLine+1:])
		b.lineIndex = b.lineIndex[:len(b.lineIndex)-1]

		b.cursorLine--
		for i := b.cursorLine + 1; i < len(b.lineIndex); i++ {
			b.lineIndex[i]--
		}

		b.cursorCol = b.lineLen(b.cursorLine)

	}

	return true

}

// NewLine inserts a new line at the cursor position and move the cursor to the left one position.
func (b *Buffer) NewLine() {
	b.insertRune(lineEndingRune)

	b.lineIndex = append(b.lineIndex, 0)
	copy(b.lineIndex[b.cursorLine+2:], b.lineIndex[b.cursorLine:])
	b.lineIndex[b.cursorLine+1] = b.lookupRCIdx(b.cursorLine, b.cursorCol) + 1

	b.cursorLine++
	for i := b.cursorLine + 1; i < len(b.lineIndex); i++ {
		b.lineIndex[i]++
	}

	b.cursorCol = 0

}

func (buffer *Buffer) Runes() iter.Seq[rune] {
	return func(yield func(rune) bool) {
		for _, r := range buffer.buffer {
			if !yield(r) {
				return
			}
		}
	}
}

func (buffer *Buffer) CursorPosition() (row, col int) {
	return buffer.cursorLine, buffer.cursorCol
}

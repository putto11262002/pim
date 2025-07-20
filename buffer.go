package pim

import (
	"fmt"
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
	// size is the current number of runes (characters) present in the 'buffer'.
	// This explicitly tracks the logical length of the text, distinct from the
	// underlying 'buffer' slice's capacity.
	len int

	point     int
	pointLine int
	name      string
}

func NewBuffer() *Buffer {
	return &Buffer{
		lineIndex: []int{0}, // Initialize with the first line starting at index 0
	}

}

func (b Buffer) Point() int {
	return b.point
}

// line returns the line number for a given point in the buffer.
func (b Buffer) line(p int) int {
	if p < 0 || p >= len(b.buffer) {
		panic(fmt.Sprintf("point out of range: %d", p))
	}

	if p == 0 {
		return 0
	}

	// TODO: Use binary search for efficiency
	for cl, idx := range b.lineIndex[:len(b.lineIndex)-1] {
		nl := b.lineIndex[cl+1]
		if p >= idx && p < nl {
			return cl
		}
	}
	return len(b.lineIndex) - 1 // Last line if point is beyond the last line index
}

// SetPoint sets the current point in the buffer to the specified index.
// The point is placed before the character at the specified index.
func (b *Buffer) SetPoint(idx int) bool {
	if idx < 0 || idx > b.len {
		return false
	}
	b.point = idx
	b.pointLine = b.line(idx)
	return true
}

// MovePoint moves the current point by n positions.
// If n is positive, it moves forward; if negative, it moves backward.
func (b *Buffer) MovePoint(n int) bool {
	newPoint := b.point + n
	if newPoint < 0 || newPoint > b.len {
		return false
	}
	b.SetPoint(newPoint)
	return true
}

// insertRunes inserts a slice of runes starting at the specified virtual point p without updating
// the internal point.
// It returns true if the insertion was successful, along the position and line where p ended up.
// insertRunes does not update the point position.
func (b *Buffer) insertRunes(p int, char []rune) (bool, int, int) {

	if p < 0 || p > b.len {
		return false, 0, 0
	}

	nchar := len(char)

	if b.len+nchar > len(b.buffer) {
		l := len(b.buffer)
		if l == 0 {
			l = initialBufferSize
		}

		newbuf := make([]rune, l*2)
		copy(newbuf, b.buffer[:b.len])
		b.buffer = newbuf
	}

	copy(b.buffer[p+nchar:], b.buffer[p:])
	copy(b.buffer[p:], char)

	pl := b.line(p)
	for p < p+nchar {
		c := b.buffer[p]

		if c == lineEndingRune {
			b.lineIndex = append(b.lineIndex, 0)
			if pl < len(b.lineIndex)-2 {
				copy(b.lineIndex[pl+2:], b.lineIndex[pl+1:])
			}
			b.lineIndex[pl+1] = p + 1
			pl++
		}
		p++
	}

	b.len += nchar

	return true, p, pl
}

func (b *Buffer) InsertRunes(char []rune) bool {

	ok, newPoint, newPointLine := b.insertRunes(b.point, char)
	if !ok {
		return false
	}

	b.point = newPoint
	b.pointLine = newPointLine

	return true
}

// deleteRunes deletes n runes starting at the specified position p without moving the internal point.
// It returns true if the deletion was successful, along with the new point and line.
func (b *Buffer) deleteRunes(p int, n int) (bool, int, int) {
	if n <= 0 || p-n < 0 {
		return false, 0, 0
	}

	lineToDelete := 0
	for i := p - 1; i > p-n-1; i-- {
		if b.buffer[i] == lineEndingRune {
			lineToDelete++
		}
	}

	pl := b.line(p)
	if lineToDelete > 0 {
		copy(b.lineIndex[pl-lineToDelete:], b.lineIndex[pl+1:])
		b.lineIndex = b.lineIndex[:len(b.lineIndex)-lineToDelete]
		for i := pl - lineToDelete + 1; i < len(b.lineIndex); i++ {
			b.lineIndex[i] -= n
		}
	}

	copy(b.buffer[p-n:], b.buffer[p:])
	b.len -= n

	return true, p - n, pl - lineToDelete
}

func (b *Buffer) DeleteRunes(n int) bool {
	ok, newPoint, newPointLine := b.deleteRunes(b.point, n)
	if !ok {
		return false
	}

	b.point = newPoint
	b.pointLine = newPointLine

	return true
}

func (b *Buffer) Len() int {
	return b.len
}

func (b *Buffer) String() string {
	if b.len == 0 {
		return ""
	}
	return string(b.buffer[:b.len])
}

func (b *Buffer) LineCount() int {
	return len(b.lineIndex)
}

// Size returns the size of the buffer content is bytes.
func (b *Buffer) Size() int {
	return len(string(b.buffer[:b.len]))
}

package pim

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func TestBuffer_InsertChar(t *testing.T) {
// 	b := NewBuffer()
// 	assertBufferState(t, defaultBufferState(), *b)
// }

func TestBuffer_insertRunes(t *testing.T) {
	b := NewBuffer()

	// Insert a character at the end
	ok, np, npl := b.insertRunes(0, []rune{'a'})
	require.True(t, ok)
	require.Equal(t, np, 1)
	require.Equal(t, npl, 0)
	assert.Equal(t, 1, b.len)
	assert.Equal(t, []int{0}, b.lineIndex)
	assert.Equal(t, []rune{'a'}, b.buffer[:b.len])

	// Insert a new line at the end
	ok, np, npl = b.insertRunes(1, []rune{lineEndingRune})
	require.True(t, ok)
	require.Equal(t, np, 2)
	require.Equal(t, npl, 1)
	assert.Equal(t, 2, b.len)
	assert.Equal(t, []int{0, 2}, b.lineIndex)
	assert.Equal(t, []rune{'a', lineEndingRune}, b.buffer[:b.len])

	// Insert a new character beween existing characters
	ok, np, npl = b.insertRunes(1, []rune{'b'})
	require.True(t, ok)
	require.Equal(t, np, 2)
	require.Equal(t, npl, 0)
	assert.Equal(t, 3, b.len)
	assert.Equal(t, []int{0, 2}, b.lineIndex)
	assert.Equal(t, []rune{'a', 'b', lineEndingRune}, b.buffer[:b.len])
}

func TestBuffer_deleteRunes(t *testing.T) {
	tests := []struct {
		name              string
		initialBuffer     []rune
		initialLineIndex  []int
		initialP          int // This is 'p' for deleteRunes
		n                 int
		expectedContent   string
		expectedNewPoint  int // This is the newPoint returned by deleteRunes
		expectedNewLine   int // This is the newPointLine returned by deleteRunes
		expectedLineIndex []int
		expectedOk        bool
	}{
		{
			name:              "Delete 1 rune from middle",
			initialBuffer:     []rune("abcde"),
			initialLineIndex:  []int{0},
			initialP:          3, // point at 'd'
			n:                 1,
			expectedContent:   "abce",
			expectedNewPoint:  2, // point at 'c'
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete 2 runes from beginning",
			initialBuffer:     []rune("abcde"),
			initialLineIndex:  []int{0},
			initialP:          2, // point at 'c'
			n:                 2,
			expectedContent:   "cde",
			expectedNewPoint:  0, // point at 'c'
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete 3 runes from end",
			initialBuffer:     []rune("abcde"),
			initialLineIndex:  []int{0},
			initialP:          5, // point after 'e'
			n:                 3,
			expectedContent:   "ab",
			expectedNewPoint:  2, // point after 'b'
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete across line ending",
			initialBuffer:     []rune("abc\ndef"),
			initialLineIndex:  []int{0, 4},
			initialP:          4, // point at 'd'
			n:                 2, // delete '\n' and 'd'
			expectedContent:   "abcef",
			expectedNewPoint:  2, // point at 'c'
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete multiple lines",
			initialBuffer:     []rune("line1\nline2\nline3"),
			initialLineIndex:  []int{0, 6, 12},
			initialP:          17, // point after 'line3'
			n:                 12, // delete "line2\nline3"
			expectedContent:   "line1",
			expectedNewPoint:  5, // point after '1'
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete more than available",
			initialBuffer:     []rune("abc"),
			initialLineIndex:  []int{0},
			initialP:          3,
			n:                 5,
			expectedContent:   "abc", // Should not change
			expectedNewPoint:  0,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        false,
		},
		{
			name:              "Delete zero runes",
			initialBuffer:     []rune("abc"),
			initialLineIndex:  []int{0},
			initialP:          2,
			n:                 0,
			expectedContent:   "abc",
			expectedNewPoint:  0,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        false,
		},
		{
			name:              "Delete negative runes",
			initialBuffer:     []rune("abc"),
			initialLineIndex:  []int{0},
			initialP:          2,
			n:                 -1,
			expectedContent:   "abc",
			expectedNewPoint:  0,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        false,
		},
		{
			name:              "Delete from empty buffer",
			initialBuffer:     []rune(""),
			initialLineIndex:  []int{0},
			initialP:          0,
			n:                 1,
			expectedContent:   "",
			expectedNewPoint:  0,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        false,
		},
		{
			name:              "Delete first rune",
			initialBuffer:     []rune("abc"),
			initialLineIndex:  []int{0},
			initialP:          1,
			n:                 1,
			expectedContent:   "bc",
			expectedNewPoint:  0,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete last rune",
			initialBuffer:     []rune("abc"),
			initialLineIndex:  []int{0},
			initialP:          3,
			n:                 1,
			expectedContent:   "ab",
			expectedNewPoint:  2,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete all runes",
			initialBuffer:     []rune("abc"),
			initialLineIndex:  []int{0},
			initialP:          3,
			n:                 3,
			expectedContent:   "",
			expectedNewPoint:  0,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete a line ending at the end of the buffer",
			initialBuffer:     []rune("abc\n"),
			initialLineIndex:  []int{0, 4},
			initialP:          4, // point after '\n'
			n:                 1,
			expectedContent:   "abc",
			expectedNewPoint:  3,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete a line ending in the middle",
			initialBuffer:     []rune("abc\ndef"),
			initialLineIndex:  []int{0, 4},
			initialP:          4, // point at 'd'
			n:                 1, // delete '\n'
			expectedContent:   "abcdef",
			expectedNewPoint:  3,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
		{
			name:              "Delete multiple line endings",
			initialBuffer:     []rune("a\nb\nc\n"),
			initialLineIndex:  []int{0, 2, 4, 6},
			initialP:          6, // point at 'c'
			n:                 4, // delete 'b\nc\n'
			expectedContent:   "a",
			expectedNewPoint:  2,
			expectedNewLine:   0,
			expectedLineIndex: []int{0},
			expectedOk:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := NewBuffer()
			b.buffer = make([]rune, len(tt.initialBuffer))
			copy(b.buffer, tt.initialBuffer)
			b.len = len(tt.initialBuffer)
			b.lineIndex = make([]int, len(tt.initialLineIndex))
			copy(b.lineIndex, tt.initialLineIndex)

			ok, newPoint, newPointLine := b.deleteRunes(tt.initialP, tt.n)

			assert.Equal(t, tt.expectedOk, ok, "Expected ok status mismatch")
			assert.Equal(t, tt.expectedContent, b.String(), "Expected content mismatch")
			assert.Equal(t, tt.expectedNewPoint, newPoint, "Expected newPoint mismatch")
			assert.Equal(t, tt.expectedNewLine, newPointLine, "Expected newPointLine mismatch")
			assert.Equal(t, tt.expectedLineIndex, b.lineIndex, "Expected lineIndex mismatch")
		})
	}
}

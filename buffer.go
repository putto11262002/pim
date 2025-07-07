package pim

import "iter"

type Buffer struct {
	buffer [][]rune
	// cursor represents the current position in the buffer.
	// The position is represented as a row and column and it is 0-indexed.
	// In EditMode, the current position is where the next character will be inserted.
	// In CommandMode, the cursor position is where the command will be executed.
	cursor Cursor
}

func NewBuffer() *Buffer {
	return &Buffer{
		buffer: make([][]rune, 1),
		cursor: Cursor{
			Row: 0,
			Col: 0,
		},
	}
}

func (buffer *Buffer) InsertChars(chars []rune) {
	// TODO: insert at cursor position not at the end
	buffer.buffer[buffer.cursor.Row] = append(buffer.buffer[buffer.cursor.Row], chars...)
	buffer.cursor.Col += len(chars)
}

func (buffer *Buffer) DeleteChar() bool {

	if buffer.cursor.Col > 0 {
		// Delete the character before the cursor
		buffer.buffer[buffer.cursor.Row] = append(buffer.buffer[buffer.cursor.Row][:buffer.cursor.Col-1],
			buffer.buffer[buffer.cursor.Row][buffer.cursor.Col:]...)
		buffer.cursor.Col--
		return true
	} else if buffer.cursor.Row > 0 {
		// No more to delete from current line
		// Delete current line and move up to the last character of the previous line if exist
		buffer.buffer = buffer.buffer[:buffer.cursor.Row]

		buffer.cursor.Row--
		buffer.cursor.Col = len(buffer.buffer[buffer.cursor.Row])
	}

	return false

}

func (buffer *Buffer) InsertNewLineBelow() {
	buffer.buffer = append(buffer.buffer, []rune{})
	buffer.cursor.Row++
	buffer.cursor.Col = 0
}

type Cursor struct {
	Row int
	Col int
}

func (buffer *Buffer) Lines() iter.Seq[[]rune] {
	return func(yield func([]rune) bool) {
		for _, line := range buffer.buffer {
			if !yield(line) {
				return
			}
		}
	}
}

func (buffer *Buffer) CursorPosition() (row, col int) {
	return buffer.cursor.Row, buffer.cursor.Col
}

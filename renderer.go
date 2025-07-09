package pim

import (
	"bufio"
	"fmt"
	"os"

	"golang.org/x/term"
)

const (
	// clearEscape is the escape sequence to clear the terminal screen
	clearEscape = "\033[2J"

	enterAlternateScreen = "\033[?1049h"
	exitAlternateScreen  = "\033[?1049l"

	cursorHome            = "\033[H" // Move cursor to the home position (top-left corner)
	moveCursorLeft        = "\033[D"
	moveCursorPosition    = "\033[%d;%dH" // Move cursor to row and column (1-indexed)
	saveCursorPosition    = "\033[s"      // Save cursor position
	restoreCursorPosition = "\033[u"      // Restore cursor position
	clearLine             = "\033[2K"     // Clear the current line
)

func NewRenderer(stdout *os.File) *Renderer {
	return &Renderer{
		stdout: stdout,
		writer: bufio.NewWriter(stdout),
	}
}

type Renderer struct {
	stdout        *os.File
	writer        *bufio.Writer
	height        int
	width         int
	terminalState *term.State
}

// cursorToPosition moves the cursor to the specified row and column in the terminal. The row and column are 1-indexed.
func (r *Renderer) cursorToPosition(row, col int) error {
	_, err := fmt.Fprintf(r.writer, moveCursorPosition, row, col)
	return err
}

func (r *Renderer) Initialize() error {
	var err error

	if _, err := r.writer.WriteString(enterAlternateScreen); err != nil {
		return fmt.Errorf("failed to enter alternate screen: %w", err)
	}

	r.terminalState, err = term.MakeRaw(int(r.stdout.Fd()))
	if err != nil {
		return fmt.Errorf("failed to set terminal to raw mode: %w", err)
	}

	// Get terminal size
	r.width, r.height, err = term.GetSize(int(r.stdout.Fd()))
	if err != nil {
		return fmt.Errorf("failed to get terminal size: %w", err)
	}

	return nil
}

func (r *Renderer) Cleanup() error {
	if err := term.Restore(int(r.stdout.Fd()), r.terminalState); err != nil {
		return fmt.Errorf("failed to restore terminal state: %w", err)
	}
	if _, err := r.writer.WriteString(exitAlternateScreen); err != nil {
		return fmt.Errorf("failed to exit alternate screen: %w", err)
	}
	r.writer.Flush()
	return nil
}

func (r *Renderer) saveCursorPosition() error {
	_, err := fmt.Fprint(r.writer, saveCursorPosition)
	return err
}

func (r *Renderer) restoreCursorPosition() error {
	_, err := fmt.Fprint(r.writer, restoreCursorPosition)
	return err
}

func (r *Renderer) cursorHome() error {
	_, err := fmt.Fprint(r.writer, cursorHome)
	return err
}

func (r *Renderer) Render(editor *Editor) error {

	// Move the cursor to the home position
	if err := r.cursorHome(); err != nil {
		return fmt.Errorf("failed to move cursor to home position: %w", err)
	}

	// Clear the screen
	if _, err := fmt.Fprint(r.writer, clearEscape); err != nil {
		return fmt.Errorf("failed to clear screen: %w", err)
	}

	// Render the buffer
	for char := range editor.buffer.Runes() {
		fmt.Fprint(r.writer, string(char))
	}

	// Move the cursor to the current position
	row, col := editor.buffer.CursorPosition()
	if err := r.cursorToPosition(row+1, col+1); err != nil {
		return fmt.Errorf("failed to move cursor to position (%d, %d): %w", row+1, col+1, err)
	}

	// Render status line
	if err := r.saveCursorPosition(); err != nil {
		return fmt.Errorf("failed to save cursor position: %w", err)
	}

	if err := r.cursorToPosition(r.height, 1); err != nil {
		return fmt.Errorf("failed to move cursor to status line: %w", err)
	}

	fmt.Fprintf(r.writer, clearLine)

	fmt.Fprintf(r.writer, "Mode: %s | Cursor: (%d, %d)", editor.mode, row+1, col+1)
	if err := r.restoreCursorPosition(); err != nil {
		return fmt.Errorf("failed to restore cursor position: %w", err)
	}

	// IMPORTANT: Flush the buffered writer to ensure immediate display
	if err := r.writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush output: %w", err)
	}

	return nil
}

package pim

import (
	"atomicgo.dev/keyboard/keys"
)

type EditorMode uint8

const (
	EditMode EditorMode = iota
	CommandMode
)

func (em EditorMode) String() string {
	switch em {
	case EditMode:
		return "Edit"
	case CommandMode:
		return "Command"
	default:
		return "?"
	}
}

type Editor struct {
	mode   EditorMode
	buffer *Buffer
}

func NewEditor() *Editor {
	return &Editor{
		mode:   EditMode,
		buffer: NewBuffer(),
	}
}

func (e *Editor) SetMode(mode EditorMode) {
	e.mode = mode
}

func (e *Editor) Mode() EditorMode {
	return e.mode
}

func (e *Editor) Execute(key keys.Key) error {
	if e.mode == CommandMode {
		switch key.Code {
		case keys.RuneKey:
			switch key.String() {
			case "i":
				e.SetMode(EditMode)
			}

		}
	} else if e.mode == EditMode {
		switch key.Code {
		case keys.Esc:
			e.SetMode(CommandMode)
		case keys.Backspace:
			e.buffer.DeleteRune()
		case keys.Enter:
			e.buffer.NewLine()
		default:
			for _, c := range key.Runes {
				e.buffer.InsertRune(c)
			}
		}
	}

	return nil

}

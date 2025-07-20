## Feature Set Document: Simple TUI Text Editor

**Project Name:** Pim
**Version:** 0.1.0 (Initial Release)
**Target Audience:** Developers and users who prefer a lightweight, terminal-based text editing experience.

### 1. Introduction

This document outlines the core features and functionalities for the initial release of the Pim TUI text editor. The goal is to provide a simple yet functional editor with essential text manipulation and file management capabilities within a terminal environment, leveraging a **mode-based interaction model**.

### 2. Scope

The initial scope focuses on fundamental text editing, basic file operations, and a clear, intuitive terminal user interface. The editor operates on a single buffer. Advanced features like syntax highlighting, complex search/replace, or multi-buffer management are explicitly out of scope for this version but may be considered for future iterations.

### 3. Core Features

#### 3.1. Mode-Based Interaction

The editor operates primarily in two distinct modes:

- **Edit Mode:** For direct text input and basic cursor navigation.
- **Command Mode:** For executing editor commands and advanced navigation.

#### 3.2. Text Editing

- **Text Insertion:** Users can type characters, and they will be inserted at the current cursor position. (Already implemented)
- **Text Deletion:**
  - **Backspace:** Deletes the character to the left of the cursor. (Already implemented)
  - **Delete:** Deletes the character under the cursor. (To be implemented)
- **New Lines:** Pressing `Enter` inserts a new line and moves the cursor to the beginning of the next line. (Already implemented)
- **Undo:** Reverts the last text modification. (To be implemented)

#### 3.3. Cursor Movement & Navigation

- **Edit Mode Navigation:**
  - **Arrow Keys:** Move cursor left, right, up, and down. (To be implemented for full functionality)
- \*\*Command Mode Navigation (Vim-like):
  - `h`: Move cursor left.
  - `j`: Move cursor down.
  - `k`: Move cursor up.
  - `l`: Move cursor right.
  - **Line-wise:** Move to the beginning (`0`) and end (`$`) of the current line. (To be implemented)
  - **Word-wise:** Move cursor to the beginning/end of the next/previous word (`w`, `b`, `e`). (To be implemented)

#### 3.4. User Interface & Experience

- **Terminal Integration:**
  - Utilize alternate screen buffer for a clean editing environment. (Already implemented)
  - Proper handling of terminal resizing (re-render content).
- **Status Bar:** A persistent line at the bottom of the terminal displaying:
  - Current editor mode (e.g., `EDIT`, `COMMAND`). (Already implemented)
  - Current cursor position (Line, Column). (Already implemented)
  - Current file name (e.g., `filename.txt`).
  - Modified status (e.g., `[Modified]` if changes are unsaved).

#### 3.5. Editor Commands (Accessible via Command Mode)

- **Mode Switching:**
  - `i`: Enter Edit Mode (insert at current cursor position).
  - `a`: Enter Edit Mode (insert after current cursor position).
  - `o`: Insert new line below current line and enter Edit Mode.
  - `O`: Insert new line above current line and enter Edit Mode.
- **File Management:**
  - `:q`: Quit the editor. If there are unsaved changes, prompt the user to save or discard.
  - `:q!`: Quit the editor, discarding unsaved changes.
  - `:w`: Save the current buffer content to the currently opened file.
  - `:w <filename>`: Save the current buffer content to a new file (Save As).
  - `:wq`: Save the current file and then exit.
  - `:o <filename>`: Open a specified file, loading its content into the editor buffer.
  - `:new`: Create a new empty buffer.

### 4. Underlying Text Representation

The editor will support or consider the following text representation types for efficient text manipulation:

- **Buffer with Extra Space:** A pre-allocated buffer that allows for efficient insertions and deletions by shifting runes within the existing allocated memory, minimizing frequent reallocations. (Currently implemented)
- **Gap Buffer:** A specialized data structure designed for efficient insertions and deletions near the cursor by maintaining a "gap" in the buffer. (Considered for future optimization or alternative implementation)

### 5. Future Considerations (Out of Scope for v0.1.0)

- Copy, Cut, and Paste operations.
- Redo functionality.
- Search and Replace functionality.
- Syntax highlighting.
- Line numbers display.
- Multi-file/buffer management.
- Integration with system clipboard.

#### 5.1. Current Limitations and Unimplemented Cases

- **Large File Handling:** Currently, there is no specific strategy for handling large files; they are loaded entirely into memory.
- **Binary File Support:** The editor does not fully handle binary files. Binary files are currently read-only, and invalid Unicode characters are replaced with a replacement Unicode character.

### 6. Text Encoding

Currently, the Pim TUI text editor exclusively supports UTF-8 text encoding. While this simplifies initial development, future versions may include auto-detection of text encodings and support for a wider range of character encodings to enhance compatibility with various file types.


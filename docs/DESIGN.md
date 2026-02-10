Components:

## Buffer

- This store the state of the text being edited along with some metadata and
  provides an interface for manipulating the text.

- buffer is used to store content of the text being edtied in this case we use
  `[]rune`. As is elivates the complexity of working with unicode charcaters.

- A buffer can exist by itself, or it can be associated with at most one file.
  When associated with a file, the buffer is a copy of the contents of the file
  at a specific time.

- Open Buffer from a file read the contents of the file into the buffer.

- Write buffer replace the contents of the file with the contents of the buffer.

- The concept of line is presented in the buffer as a new line character.

- Basic unit of operations in the buffer component is rune or uni code cahrcater

- The **point** is where the operations to the buffer are applied. The poiot is
  represnted with (`line_idx`, `offset`) `line_idx` is the index of the line and
  the `offset_idx` is the position between caharacers in the line the position
  before the first cahracter of the line is 0 and the position after the last
  character of the line is .... Transalting this to cahcartacter index the edit
  happens the the character at .... The point can only exist between two
  cahracters. So tehg point is between the character at `offset` and the
  character at `offset + 1` at line `line_idx`. delete

- We provide serveral implemntation of the underlying strorage data structure.

1.  `BasicBuffer` Buffer with extra space at the end. Deltetion never require
    reallocation. Insertion only requires reallocation when the buffer is full. But
    require a lot of shuffling of the data.

2.  `GapBuffer`

    - continue block fo mem with gap with unused space between them.
    - maintain a pointer to the gap_start ane gap_end.
    - Inserting and deleting characters within the gap involve moving the
      gap_start and gap_end pointers.
    - When cursor is moved outside the gap, the gap relocation is delay until
      the next edit operation is performed to avoid unnecessary shifting of
      characters.
    - Extremely efficent local edits.
    - Contiguous memory allocation.
    - Still quite simple to implement.
    - Not great for long distance edits.
    - Not great for large files as it requires reallocation requires copying the
      entire buffer.
    - Assumes that the cose of resizeing and non-frequient loing distant edits
      can be amotized over the more common local edits.
    - Considerations: tradeoff between memory usage and reallocation frequency.

## Dispaly

The display component is responsible for rendering the buffer on the display
device. In this case, it is a terminal. It takse the state of the buffer and
and sends appropriate escape sequences to the TTY to render the text on the
screen.

Requirements: It is essstaily that the repaly compoenent correct dispaly the buffer and does
so swiftly. Any changes made to the buffer should be reflected on the screen
instantly, or at least as soon as possible, to provide a smooth user.

Constraints:
It takes time to render the buffer on the screen, especially for large buffers.
There are a few constraints to consdier, the dispaly communication channel speed
(in this case TTY), the rendering device speed (the terminal), CPU speed, and
available memory. Since CPU and memory are usually not the bottleneck, we shift
most of the burden to the CPU to perofrm some logic on dertermiing what part to
replay this reduce the load on the rest as it only has to render incrementally.

To achieve

we take the advantage of the fact that changes to the buffer are applied
incrementally so we can update the display incrementally as well - only
rediplaying small parts of the creencrrent. So we use an algorithm to determine
regions of the screen that need to be updated and only update those regions.

- Do not support line break instead use horizontal scrolling.

## Basic Redisplay Algorithm

## Assumptiops

- The buffer line occupies the exactly one line on the screen.

## Frammer

The frist component is dispaly handles determining which part of the buffer to
display on the screen. We maintain two pointers to the buffer: `start_index`:
The index of the first character to display and the `end_index`: the index of
the last character to display.

As long as the point is in visible (between `start_index` and `end_index`),
the the `start_index` remains unchanged, only adjust `end_index` as needed. In
the case that the point is not visible, we adjust the `start_index` in a way
that the point is as great as or equal to the `point_pct`. Where `point_pct`
is the percentage of the lines that are visible on the screen.

## Handling Overflow

We are using horizonal scroll for now as line wrap introduce a bunch more
complexity

## Algorithm

We are using a really basic display algorith.

- store screen contents
- loop through every position in the window cacahe comparing the current buffer
  content with what's in the buffer if they don't match update the cacahe and
  draw the graphmeem cluster. This design uses O(window_size memoery. Not the
  most efficent design as event if only some line is change we need to check all
  cahracters.

- Probabvry store some flag to indicate if the buffer content has changed if not
  we can just skip the valdiation and only render the cursor, statusline

## Text Encoding

# Buffer

## Data structure

An slice of runes is used to store the text.

- Simple
- Fairly efficient for reading and jumping around. As characters are stored in contiguous memory, it is fast to access any character.
- Not very efficient for inserting or deleting characters in the middle, as it requires shifting characters around.

## Line index

A line index is used to keep track of the start of each line in the buffer. This allows for efficient access to lines and line numbers.

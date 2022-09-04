package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const GRID_SIZE = 9

type Config struct {
	OpenCellsPerc int
}

type point struct {
	x, y int
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Example board (numbers are not real).
//
// -------------------------
// | 1 2 3 | 4 5 6 | 7 8 9 |
// | 4 5 6 | 7 8 9 | 1 2 3 |
// | 7 8 9 | 1 2 3 | 5 6 7 |
// -------------------------
// | 1 2 3 | 4 5 6 | 7 8 9 |
// | 4 5 6 | 7 8 9 | 1 2 3 |
// | 7 8 9 | 1 2 3 | 5 6 7 |
// -------------------------
// | 1 2 3 | 4 5 6 | 7 8 9 |
// | 4 5 6 | 7 8 9 | 1 2 3 |
// | 7 8 9 | 1 2 3 | 5 6 7 |
// -------------------------

type Sudoku struct {
	board     [][]int
	solution  [][]int
	cursor    point
	conflicts [][]bool
}

func (s *Sudoku) Validate() {
	// Reset conflicts.
	for y := 0; y < GRID_SIZE; y++ {
		for x := 0; x < GRID_SIZE; x++ {
			s.conflicts[y][x] = false
		}
	}
}

func (s *Sudoku) IsValidCell(x, y int) bool {
	return s.conflicts[y][x]
}

func (s *Sudoku) Randomize() {
	s.clearBoard()
}

func (s *Sudoku) clearBoard() {
	s.board = make([][]int, GRID_SIZE)
	s.solution = make([][]int, GRID_SIZE)
	for i := 0; i < GRID_SIZE; i++ {
		s.board[i] = make([]int, GRID_SIZE)
		s.solution[i] = make([]int, GRID_SIZE)
	}
	s.cursor = point{0, 0}
}

func (s *Sudoku) SetValue(x, y, value int) {
	s.solution[y][x] = value
}

func (s *Sudoku) IsValueValidAt(x, y int) bool {
	return s.board[y][x] == s.solution[y][x]
}

func (s *Sudoku) IsCursorAt(x, y int) bool {
	return s.cursor.x == x && s.cursor.y == y
}

func (s *Sudoku) MoveCursor(dx, dy int) {
	s.cursor.x += dx
	s.cursor.y += dy

	if dx > 0 {
		s.cursor.x = min(s.cursor.x, GRID_SIZE-1)
	} else {
		s.cursor.x = max(s.cursor.x, 0)
	}

	if dy > 0 {
		s.cursor.y = min(s.cursor.y, GRID_SIZE-1)
	} else {
		s.cursor.y = max(s.cursor.y, 0)

	}
}

func (s *Sudoku) WriteTo(w io.Writer) (int64, error) {
	vertical := func() {
		fmt.Fprintf(w, "%s\n", strings.Repeat("-", GRID_SIZE*3+4))
	}

	for y := 0; y < GRID_SIZE; y++ {
		if y%3 == 0 {
			vertical()
		}
		fmt.Fprint(w, "|")
		for x := 0; x < GRID_SIZE; x++ {
			if s.IsCursorAt(x, y) {
				fmt.Fprintf(w, "[")
			} else {
				fmt.Fprintf(w, " ")
			}
			if s.solution[y][x] == 0 {
				fmt.Fprint(w, "*")
			} else {
				fmt.Fprintf(w, "%d", s.solution[y][x])
			}
			if s.IsCursorAt(x, y) {
				fmt.Fprintf(w, "]")
			} else {
				fmt.Fprintf(w, " ")
			}
			if (x+1)%3 == 0 {
				fmt.Fprintf(w, "|")
			}
		}
		fmt.Fprintln(w)
	}
	vertical()
	return 0, nil
}

func NewSudoku() *Sudoku {
	return &Sudoku{}
}

func start() {
	sudoku := NewSudoku()
	sudoku.Randomize()
	sudoku.WriteTo(os.Stdout)

	quit := true
	for !quit {
		// TODO: Readline and make some action.

		sudoku.Validate()
	}
}

func main() {
	start()
}

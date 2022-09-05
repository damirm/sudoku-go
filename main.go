package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	REGION_SIZE    = 3
	REGIONS_IN_ROW = 3
	GRID_SIZE      = REGION_SIZE * REGIONS_IN_ROW
	// Arithmetic progression from 1 to 9
	EXPECTED_SUM = GRID_SIZE * (1 + GRID_SIZE) / 2
)

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
	// board is valid combination.
	board [][]int
	// solution is user values.
	solution [][]int
	// conflicts is invalid user values.
	conflicts [][]bool
	cursor    point
}

// Validate checks if "solution" has valid combinations,
// otherwise "conflicts" slice will contain invalid cells.
func (s *Sudoku) Validate() {
	// Reset conflicts.
	for y := 0; y < GRID_SIZE; y++ {
		for x := 0; x < GRID_SIZE; x++ {
			s.conflicts[y][x] = false
		}
	}

	for y := 0; y < GRID_SIZE; y++ {
		for x := 0; x < GRID_SIZE; x++ {
			if s.solution[y][x] == 0 {
				continue
			}

			// rx := int(math.Ceil(float64(x) / 3.0))
			// ry := int(math.Ceil(float64(y) / 3.0))
			valid := s.validateRow(y) ||
				s.validateCol(x) ||
				s.validateRegion(x/REGION_SIZE, y/REGION_SIZE)

			if !valid {
				s.conflicts[y][x] = true
			}
		}
	}
}

func (s *Sudoku) validateCol(x int) bool {
	sum := 0
	seen := make(map[int]bool, GRID_SIZE)
	for i := 0; i < GRID_SIZE; i++ {
		n := s.solution[i][x]
		if exists, _ := seen[n]; exists {
			return false
		}
		seen[n] = true
		sum += n
	}
	return sum <= EXPECTED_SUM
}

func (s *Sudoku) validateRow(y int) bool {
	sum := 0
	seen := make(map[int]bool, GRID_SIZE)
	for i := 0; i < GRID_SIZE; i++ {
		n := s.solution[y][i]
		if exists, _ := seen[n]; exists {
			return false
		}
		seen[n] = true
		sum += n
	}
	return sum <= EXPECTED_SUM
}

// validateRegion takes region coords (e.g. x=0, y=0 is the first region, x=1, y=0 is the second region).
func (s *Sudoku) validateRegion(x, y int) bool {
	sum := 0
	seen := make(map[int]bool, GRID_SIZE)
	for iy := 0; iy < REGION_SIZE; iy++ {
		for ix := 0; ix < REGION_SIZE; ix++ {
			px := x*REGION_SIZE + ix
			py := y*REGION_SIZE + iy
			n := s.solution[py][px]
			if exists, _ := seen[n]; exists {
				return false
			}
			seen[n] = true
			sum += n
		}
	}
	return sum <= EXPECTED_SUM
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

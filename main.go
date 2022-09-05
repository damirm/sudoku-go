package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"golang.org/x/term"
)

const (
	REGION_SIZE    = 3
	REGIONS_IN_ROW = 3
	GRID_SIZE      = REGION_SIZE * REGIONS_IN_ROW
	VERT_BORDERS   = REGIONS_IN_ROW + 1
	// Arithmetic progression from 1 to 9
	EXPECTED_SUM = GRID_SIZE * (1 + GRID_SIZE) / 2

	COLOR_RED   = "\x1B[31m"
	COLOR_RESET = "\x1B[0m"
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

func random(start, end int) int {
	return rand.Intn(end-start) + start
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
	for iy := 0; iy < GRID_SIZE; iy++ {
		n := s.solution[iy][x]
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
	for ix := 0; ix < GRID_SIZE; ix++ {
		n := s.solution[y][ix]
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
	return !s.conflicts[y][x]
}

// HasCompletedValidArea returns true if region or row or col
// is filled and it is valid.
func (s *Sudoku) HasCompletedValidArea(x, y int) bool {
	return (s.isRowFilled(y) && s.validateRow(y)) ||
		(s.isColFilled(x) && s.validateCol(x)) ||
		(s.isRegionFilled(x/REGION_SIZE, y/REGION_SIZE) &&
			s.validateRegion(x/REGION_SIZE, y/REGION_SIZE))
}

func (s *Sudoku) isRowFilled(y int) bool {
	for ix := 0; ix < GRID_SIZE; ix++ {
		if s.solution[y][ix] == 0 {
			return false
		}
	}
	return true
}

func (s *Sudoku) isColFilled(x int) bool {
	for iy := 0; iy < GRID_SIZE; iy++ {
		if s.solution[iy][x] == 0 {
			return false
		}
	}
	return true
}

func (s *Sudoku) isRegionFilled(x, y int) bool {
	for iy := 0; iy < REGION_SIZE; iy++ {
		for ix := 0; ix < REGION_SIZE; ix++ {
			px := x*REGION_SIZE + ix
			py := y*REGION_SIZE + iy
			if s.solution[py][px] == 0 {
				return false
			}
		}
	}
	return true
}

func (s *Sudoku) Randomize() {
	s.clearBoard()

	for y := 0; y < GRID_SIZE; y++ {
		for x := 0; x < GRID_SIZE; x++ {
			v := random(1, GRID_SIZE+1)

			if random(0, 2) == 0 {
				s.solution[y][x] = v
			}
		}
	}
}

func (s *Sudoku) clearBoard() {
	s.board = make([][]int, GRID_SIZE)
	s.solution = make([][]int, GRID_SIZE)
	s.conflicts = make([][]bool, GRID_SIZE)
	for i := 0; i < GRID_SIZE; i++ {
		s.board[i] = make([]int, GRID_SIZE)
		s.solution[i] = make([]int, GRID_SIZE)
		s.conflicts[i] = make([]bool, GRID_SIZE)
	}
	s.cursor = point{0, 0}
}

func (s *Sudoku) setValue(x, y, value int) {
	s.solution[y][x] = value
}

func (s *Sudoku) SetValueUnderCursor(value int) {
	s.setValue(s.cursor.x, s.cursor.y, value)
}

func (s *Sudoku) ClearValueUnderCursor() {
	s.SetValueUnderCursor(0)
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
	hr := func() {
		fmt.Fprintf(w, "%s\x1B[G\n", strings.Repeat("-", REGION_SIZE*REGIONS_IN_ROW*3+VERT_BORDERS))
	}

	for y := 0; y < GRID_SIZE; y++ {
		if y%3 == 0 {
			hr()
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
				val := s.solution[y][x]
				if s.IsValidCell(x, y) {
					fmt.Fprintf(w, "%d", val)
				} else {
					fmt.Fprintf(w, "%s%d%s", COLOR_RED, val, COLOR_RESET)
				}
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
		fmt.Fprint(w, "\x1B[G\n")
	}
	hr()

	fmt.Printf("\x1B[%dD\x1B[%dA", REGION_SIZE*REGIONS_IN_ROW+VERT_BORDERS, REGION_SIZE*REGIONS_IN_ROW+VERT_BORDERS)

	return 0, nil
}

func NewSudoku() *Sudoku {
	return &Sudoku{}
}

func isDigit(c rune) bool {
	return c >= '1' && c <= '9'
}

func start() {
	rand.Seed(time.Now().UnixNano())

	sudoku := NewSudoku()
	sudoku.Randomize()
	sudoku.Validate()
	sudoku.WriteTo(os.Stdout)

	reader := bufio.NewReaderSize(os.Stdin, 1)

	quit := false
	for !quit {
		input, _, _ := reader.ReadRune()
		switch input {
		case 'h':
			sudoku.MoveCursor(-1, 0)
		case 'j':
			sudoku.MoveCursor(0, 1)
		case 'k':
			sudoku.MoveCursor(0, -1)
		case 'l':
			sudoku.MoveCursor(1, 0)
		case 'q':
			quit = true
		case ' ':
			sudoku.ClearValueUnderCursor()
		default:
			if isDigit(input) {
				sudoku.SetValueUnderCursor(int(input - '0'))
				sudoku.Validate()
			}
		}

		sudoku.WriteTo(os.Stdout)
	}
}

func main() {
	state, err := term.MakeRaw(0)
	if err != nil {
		log.Fatal(err)
	}
	defer term.Restore(0, state)

	start()
}

package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
)

var (
	UpVec    = Vector{Y: -1}
	RightVec = Vector{X: 1}
	DownVec  = Vector{Y: 1}
	LeftVec  = Vector{X: -1}

	AllDirections = []Vector{
		UpVec,
		RightVec,
		DownVec,
		LeftVec,
	}
)

type ConnectionMap map[Vector]map[uint8]bool

var connectionMap = ConnectionMap{
	UpVec: {
		'.': false,
		'-': false,
		'7': true,
		'|': true,
		'F': true,
		'J': false,
		'L': false,
		'S': true,
	},
	RightVec: {
		'.': false,
		'-': true,
		'7': true,
		'F': false,
		'|': false,
		'J': true,
		'L': false,
		'S': true,
	},
	DownVec: {
		'.': false,
		'-': false,
		'7': false,
		'F': false,
		'|': true,
		'J': true,
		'L': true,
		'S': true,
	},
	LeftVec: {
		'.': false,
		'-': true,
		'7': false,
		'F': true,
		'|': false,
		'J': false,
		'L': true,
		'S': true,
	},
}

func (m *ConnectionMap) Get(dir Vector, c uint8) bool {
	cMap, ok := (*m)[dir]
	if !ok {
		panic(fmt.Errorf("dir %v not found in connection map", dir))
	}
	connected, ok := cMap[c]
	if !ok {
		panic(fmt.Errorf("c %c not found in connection map", c))
	}
	return connected
}

func (v Vector) Add(o Vector) Vector {
	return Vector{
		X: v.X + o.X,
		Y: v.Y + o.Y,
	}
}

type Vector struct {
	X, Y int
}

type Grid [][]uint8

func getChar(lines Grid, pos Vector) uint8 {
	x := pos.X
	y := pos.Y
	if y < 0 || y >= len(lines) {
		return '.'
	}
	line := lines[y]
	if x < 0 || x >= len(line) {
		return '.'
	}
	c := line[x]
	return c
}

func ConnectedVectors(lines Grid, pos Vector) []Vector {
	vectors := []Vector{}
	dirsToCheck := []Vector{}
	origC := getChar(lines, pos)

	switch origC {
	case '.':
	case '-':
		dirsToCheck = []Vector{
			LeftVec, RightVec,
		}
	case '7':
		dirsToCheck = []Vector{
			LeftVec, DownVec,
		}
	case 'F':
		dirsToCheck = []Vector{
			DownVec, RightVec,
		}
	case '|':
		dirsToCheck = []Vector{
			UpVec, DownVec,
		}
	case 'J':
		dirsToCheck = []Vector{
			LeftVec, UpVec,
		}
	case 'L':
		dirsToCheck = []Vector{
			UpVec, RightVec,
		}
	case 'S':
		dirsToCheck = AllDirections
	default:
		panic(fmt.Errorf("c %c not in mapping", origC))
	}

	for _, dir := range dirsToCheck {
		newP := pos.Add(dir)
		newC := getChar(lines, newP)
		if connected := connectionMap.Get(dir, newC); connected {
			vectors = append(vectors, dir)
		}
	}
	return vectors
}

func Traverse(lines Grid, pos Vector, distance int, distances map[Vector]int) {

	existingDist, ok := distances[pos]
	if !ok {
		distances[pos] = distance
	} else if distance < existingDist {
		distances[pos] = distance
	} else {
		return
	}

	connectedVectors := ConnectedVectors(lines, pos)
	for _, vec := range connectedVectors {
		nextPos := pos.Add(vec)
		Traverse(lines, nextPos, distance+1, distances)
	}
}

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	var grid Grid
	var start Vector
	lineNo := 0
	for fileScanner.Scan() {
		lineBuf := fileScanner.Bytes()
		line := make([]byte, len(lineBuf))
		copy(line, lineBuf)
		grid = append(grid, line)
		if ind := bytes.Index(line, []byte{'S'}); ind != -1 {
			start = Vector{
				X: ind,
				Y: lineNo,
			}
		}
		lineNo += 1
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	distances := map[Vector]int{}
	Traverse(grid, start, 0, distances)
	insideCount := 0
	for y, line := range grid {
		for x := range line {
			if _, ok := distances[Vector{X: x, Y: y}]; ok {
				continue
			}
			// Count the number of times we fully traverse the loop on our way out of the field.
			// If we cross an odd number of times, we are inside the loop. If even or zero, we are out.
			// "Fully traverse" meaning pass over a character where we actually go from one side to another.
			// "Corners" like "L" and "7" do not cause us to cross over, if we move the cursor diagonally down and right.
			// We remain on the same side as when we started passing over those characters.
			xi, yi := x, y
			crosses := 0
			for xi < len(line) && yi < len(grid) {
				pos := Vector{X: xi, Y: yi}
				c := getChar(grid, pos)
				if _, ok := distances[pos]; ok && c != 'L' && c != '7' {
					crosses += 1
				}
				xi += 1
				yi += 1
			}
			if crosses%2 == 1 {
				insideCount += 1
				grid[y][x] = 'I'
			}
		}
	}

	fmt.Printf("InsideCount: %d\n", insideCount)
}

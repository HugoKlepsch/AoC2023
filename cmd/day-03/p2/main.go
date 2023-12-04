package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

func isNumber(c uint8) bool {
	return c >= '0' && c <= '9'
}

func getChar(lines []string, x int, y int) uint8 {
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

func isSymbol(c uint8) bool {
	if isNumber(c) {
		return false
	}
	return c != '.' && c != '\n'
}

type StringPosition struct {
	x, y, length int
}

// parseNumber will return the number containing x, y, startX, and its length.
func parseNumber(lines []string, x, y int) (int, int, int) {
	var (
		start = x
		end   = x
	)
	line := lines[y]
	// find start
	for i := start; i >= 0; i-- {
		c := line[i]
		if !isNumber(c) {
			break
		}
		start = i
	}
	// find end
	for i := start; i < len(line); i++ {
		c := line[i]
		if !isNumber(c) {
			break
		}
		end = i
	}
	numStr := line[start : end+1]
	num, err := strconv.Atoi(numStr)
	if err != nil {
		panic(fmt.Errorf("Could not Atoi number: %w", err))
	}
	return num, start, end + 1 - start
}

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	var lines []string

	for fileScanner.Scan() {
		line := fileScanner.Text()
		lines = append(lines, line)
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	score := 0

	for y, line := range lines {
		for x := 0; x < len(line); x++ {
			c := getChar(lines, x, y)
			if c == '*' {
				adjacentParts := map[StringPosition]int{} // value is part number

				for yi := y - 1; yi <= y+1; yi++ {
					for xi := x - 1; xi <= x+1; xi++ {
						ci := getChar(lines, xi, yi)
						if isNumber(ci) {
							number, start, length := parseNumber(lines, xi, yi)
							fmt.Printf("Found adjacent number for %c (%d, %d): %d[%d] (%d, %d)\n", c, x, y, number, length, xi, yi)
							pos := StringPosition{
								x:      start,
								y:      yi,
								length: length,
							}
							adjacentParts[pos] = number
						}
					}
				}

				fmt.Printf("Adjacent numbers for %c (%d, %d):\n", c, x, y)
				for pos, number := range adjacentParts {
					fmt.Printf("\t%d[%d] (%d, %d)\n", number, pos.length, pos.x, pos.y)
				}
				fmt.Printf("---------")
				gearRatio := 1
				if len(adjacentParts) == 2 {
					for _, number := range adjacentParts {
						gearRatio *= number
					}
					score += gearRatio
				}
			}
		}
	}

	fmt.Printf("Score: %d\n", score)
}

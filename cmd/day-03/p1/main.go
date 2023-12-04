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

// parseNumber will return the number starting at x, y, and its length.
func parseNumber(lines []string, x, y int) (int, int) {
	var (
		start = x
		end   = x
	)
	line := lines[y]
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
	return num, end + 1 - start
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
	symbols := map[uint8]struct{}{}

	for y, line := range lines {
		for x := 0; x < len(line); x++ {
			c := getChar(lines, x, y)
			if isNumber(c) {
				// Number start!
				number, length := parseNumber(lines, x, y)
				hasSymbol := false
				// Is there a symbol touching this number on the previous line?
				for yi := y - 1; yi <= y+1 && !hasSymbol; yi++ {
					for xi := x - 1; xi <= x+length && !hasSymbol; xi++ {
						ci := getChar(lines, xi, yi)
						if isSymbol(ci) {
							fmt.Printf("Found symbol for %d (%d, %d): %c (%d, %d)\n", number, x, y, ci, xi, yi)
							symbols[ci] = struct{}{}
							hasSymbol = true
							break
						}
					}
				}
				if hasSymbol {
					score += number
				} else {
					fmt.Printf("Did not find symbol for %d (%d, %d)\n", number, x, y)
				}
				x += length - 1
			}
		}
	}

	fmt.Printf("Symbols: \n")
	for c, _ := range symbols {
		fmt.Printf("'%c'\n", c)
	}
	fmt.Printf("Score: %d\n", score)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

type DstTuple struct {
	L, R string
}

var lineRegex = regexp.MustCompile(`^([A-Z]+) = \(([A-Z]+), ([A-Z]+)\)`)

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan()
	directionLine := fileScanner.Text()
	directions := []uint8{}
	for i := 0; i < len(directionLine); i++ {
		directions = append(directions, directionLine[i])
	}

	fwd := map[string]DstTuple{}

	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line == "" {
			continue
		}

		if matches := lineRegex.FindStringSubmatch(line); matches != nil && len(matches) == 4 {
			src := matches[1]
			dstl := matches[2]
			dstr := matches[3]
			fwd[src] = DstTuple{
				L: dstl,
				R: dstr,
			}
		}

	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	hops := 0
	start := "AAA"
	end := "ZZZ"
	current := start
	for current != end {
		for _, direction := range directions {
			if current == end {
				break
			}
			next, ok := fwd[current]
			if !ok {
				panic(fmt.Errorf("could not find mapping for current node %s", current))
			}

			if direction == 'L' {
				current = next.L
			} else if direction == 'R' {
				current = next.R
			}
			hops += 1
		}
	}
	fmt.Printf("Hops: %d. Score: %d\n", hops, hops/len(directions))
}

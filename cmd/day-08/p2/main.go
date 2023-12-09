package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type DstTuple struct {
	L, R string
}

type Route struct {
	start, end string
	length     int
}

var lineRegex = regexp.MustCompile(`^([A-Z]+) = \(([A-Z]+), ([A-Z]+)\)`)

func In(element string, set map[string]struct{}) bool {
	_, ok := set[element]
	return ok
}

// Taken from least common divisor on Wikipedia
func lcm(in []int64) int64 {
	tmp := make([]int64, len(in))
	copy(tmp, in)

	allSame := false
	allElement := int64(0)

	for !allSame {
		lowest := tmp[0]
		lowestInd := 0

		allSame = true
		allElement = lowest

		for i, e := range tmp {
			if e != allElement {
				allSame = false
			}
			if e < lowest {
				lowest = e
				lowestInd = i
			}
		}
		tmp[lowestInd] += in[lowestInd]
	}

	return allElement
}

// Traverse traverses the chain until it reaches an end node
func Traverse(start string, fwd map[string]DstTuple, directionInd int, directions []uint8, ends map[string]struct{}) Route {
	hops := 0
	current := start
	for ; !In(current, ends); directionInd++ {
		direction := directions[directionInd%len(directions)]
		if In(current, ends) {
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
	return Route{
		start:  start,
		end:    current,
		length: hops,
	}
}

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
	starts := map[string]struct{}{}
	ends := map[string]struct{}{}

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
			if strings.HasSuffix(src, "A") {
				starts[src] = struct{}{}
			}
			if strings.HasSuffix(src, "Z") {
				ends[src] = struct{}{}
			}
		}

	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	loopLengths := []int64{}

	for start := range starts {
		directionInd := 0

		// Through empirical analysis, I have determined that the first ghost path is the longest, then
		// further ghost paths are a subset of the first path.
		// Ex:
		// 	A -> Z (len 20)
		// 	Z -> M (len: 1)
		// 	M -> Z (len: 10)
		// 	Z -> M (len: 1)
		// 	M -> Z (len: 10)
		// We can therefore describe the start->end route as a first length then a recurring loop length.
		// -----
		// Using Least Common Multiple (lcm) definition, algorithm from Wikipedia
		route := Traverse(start, fwd, directionInd, directions, ends)
		fmt.Printf("Route: %s -> %s [%d]\n", route.start, route.end, route.length)
		loopLengths = append(loopLengths, int64(route.length))
	}
	fmt.Printf("-----\n")
	score := lcm(loopLengths)
	fmt.Printf("LCM: %d\n", score)
}

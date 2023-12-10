package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func intListToString(in []int) string {
	strs := []string{}
	for _, e := range in {
		strs = append(strs, fmt.Sprintf("%d", e))
	}
	return fmt.Sprintf("[%s]", strings.Join(strs, ","))
}

func recursiveDerivativeNext(in []int) int {
	allZero := true
	for i := 0; i < len(in)-1; i++ {
		diff := in[i+1] - in[i]
		if diff != 0 {
			allZero = false
		}
		in[i] = diff
	}
	newIn := in[:len(in)-1]

	if allZero {
		predictVal := in[len(in)-1]
		return predictVal
	}
	predictVal := in[len(in)-1] + recursiveDerivativeNext(newIn)
	return predictVal
}

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	score := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		nums := []int{}
		numStrs := strings.Fields(line)
		for _, numStr := range numStrs {
			val, err := strconv.ParseInt(numStr, 10, 64)
			if err != nil {
				panic(err)
			}
			nums = append(nums, int(val))
		}
		next := recursiveDerivativeNext(nums)
		fmt.Printf("next: %s %d\n", line, next)
		score += next
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("Score: %d\n", score)
}

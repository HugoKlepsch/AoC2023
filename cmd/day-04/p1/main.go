package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	inputRegex = regexp.MustCompile(`^.+: ([^|]+) \| ([^|]+)$`)
)

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	score := 0
	for fileScanner.Scan() {
		line := fileScanner.Text()
		cardMatches := inputRegex.FindStringSubmatch(line)
		if cardMatches == nil || len(cardMatches) != 3 {
			panic(fmt.Errorf("could not parse line: %s", line))
		}
		winStr := cardMatches[1]
		haveStr := cardMatches[2]

		winners := map[int]struct{}{}
		havers := map[int]struct{}{}

		winSplit := strings.Fields(winStr)
		for _, winner := range winSplit {
			num, err := strconv.Atoi(winner)
			if err != nil {
				panic(fmt.Errorf("could not parse number: %w", err))
			}
			winners[num] = struct{}{}
		}

		haveSplit := strings.Fields(haveStr)
		for _, have := range haveSplit {
			num, err := strconv.Atoi(have)
			if err != nil {
				panic(fmt.Errorf("could not parse number: %w", err))
			}
			havers[num] = struct{}{}
		}

		cardScore := 0
		for have := range havers {
			if _, ok := winners[have]; ok {
				fmt.Printf("winner: %d\n", have)
				if cardScore == 0 {
					cardScore = 1
				} else {
					cardScore = cardScore << 1
				}
			}
		}
		score += cardScore
		fmt.Printf("-----------\n")
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("Score: %d\n", score)
}

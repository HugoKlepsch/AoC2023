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

type DefaultOneMap map[int]int

func (m DefaultOneMap) Get(i int) int {
	if val, ok := m[i]; ok {
		return val
	} else {
		m[i] = 1
		return m[i]
	}
}

func (m *DefaultOneMap) Inc(i int) {
	val := m.Get(i)
	(*m)[i] = val + 1
}

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	var copyCounts DefaultOneMap = map[int]int{}

	score := 0
	lineI := 0
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

		winnerCount := 0
		for have := range havers {
			if _, ok := winners[have]; ok {
				winnerCount++
			}
		}

		copyCount := copyCounts.Get(lineI)
		// Run scoring for this card N times, where N is the number of copies of this card we have
		for i := 0; i < copyCount; i++ {

			// When we score this card, we add additional copies of later cards
			for winnerI := lineI + 1; winnerI < (winnerCount + lineI + 1); winnerI++ {
				copyCounts.Inc(winnerI)
			}
		}
		lineI++
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	for cardI, copies := range copyCounts {
		fmt.Printf("Card %d had %d total copies\n", cardI, copies)
		score += copies // Score one for every instance of the card
	}

	fmt.Printf("Score: %d\n", score)
}

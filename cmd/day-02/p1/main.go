package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Cube int

const (
	ColorRed Cube = iota
	ColorGreen
	ColorBlue
)

type Game struct {
	ID     int
	Rounds []Round
}

type Round struct {
	CubeCount map[Cube]int
}

var (
	GameIDRegex = regexp.MustCompile(`Game ([0-9]+):`)
	RedRegex    = regexp.MustCompile(` ([0-9]+) red`)
	GreenRegex  = regexp.MustCompile(` ([0-9]+) green`)
	BlueRegex   = regexp.MustCompile(` ([0-9]+) blue`)
)

var elfGame = map[Cube]int{
	ColorRed:   12,
	ColorGreen: 13,
	ColorBlue:  14,
}

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	score := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		game := Game{}
		idResults := GameIDRegex.FindStringSubmatch(line)
		possible := true
		if idResults != nil && len(idResults) == 2 {
			game.ID, err = strconv.Atoi(idResults[1])
			if err != nil {
				panic(fmt.Errorf("failed to atoi game ID %s: %w", line, err))
			}
		} else {
			panic(fmt.Errorf("failed to parse game ID %s", line))
		}

		rounds := strings.Split(line, ";")
		for _, round := range rounds {
			redCount, err := parseCount(round, RedRegex)
			if err != nil {
				panic(fmt.Errorf("red count: %w", err))
			}
			blueCount, err := parseCount(round, BlueRegex)
			if err != nil {
				panic(fmt.Errorf("blue count: %w", err))
			}
			greenCount, err := parseCount(round, GreenRegex)
			if err != nil {
				panic(fmt.Errorf("green count: %w", err))
			}
			round := Round{CubeCount: map[Cube]int{
				ColorRed:   redCount,
				ColorGreen: greenCount,
				ColorBlue:  blueCount,
			}}

			if round.CubeCount[ColorRed] > elfGame[ColorRed] ||
				round.CubeCount[ColorGreen] > elfGame[ColorGreen] ||
				round.CubeCount[ColorBlue] > elfGame[ColorBlue] {
				fmt.Printf("Game %d was not possible: %v\n", game.ID, round)
				possible = false
			}

			game.Rounds = append(game.Rounds, round)
		}

		if possible {
			score += game.ID
		}

		fmt.Printf("game: %v\n", game)
	}

	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("score: %d\n", score)
}

func parseCount(round string, regex *regexp.Regexp) (int, error) {
	countStr := regex.FindStringSubmatch(round)
	var count int
	var err error
	if countStr != nil && len(countStr) == 2 {
		count, err = strconv.Atoi(countStr[1])
		if err != nil {
			err = fmt.Errorf("failed to atoi cube count %s: %w", round, err)
			return 0, err
		}
	} else {
		return 0, nil
	}
	return count, err
}

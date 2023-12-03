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

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	score := 0

	for fileScanner.Scan() {
		line := fileScanner.Text()
		game := Game{}
		idResults := GameIDRegex.FindStringSubmatch(line)
		gameMin := map[Cube]int{
			ColorRed:   0,
			ColorGreen: 0,
			ColorBlue:  0,
		}

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

			for cube, curMin := range gameMin {
				if round.CubeCount[cube] > curMin {
					gameMin[cube] = round.CubeCount[cube]
				}
			}

			game.Rounds = append(game.Rounds, round)
		}

		power := 1
		for _, curMin := range gameMin {
			power *= curMin
		}
		score += power

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

package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Race struct {
	time, distance int
}

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	times := []int{}
	distances := []int{}

	fileScanner.Scan()
	timeLine := fileScanner.Text()
	fileScanner.Scan()
	distanceLine := fileScanner.Text()

	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	timeStrings := strings.Fields(strings.Split(timeLine, ": ")[1])
	distanceStrings := strings.Fields(strings.Split(distanceLine, ": ")[1])

	for _, timeStr := range timeStrings {
		val, err := strconv.ParseInt(timeStr, 10, 64)
		if err != nil {
			panic(err)
		}

		times = append(times, int(val))
	}

	for _, distanceStr := range distanceStrings {
		val, err := strconv.ParseInt(distanceStr, 10, 64)
		if err != nil {
			panic(err)
		}

		distances = append(distances, int(val))
	}

	races := []Race{}
	for i := range times {
		races = append(races, Race{
			time:     times[i],
			distance: distances[i],
		})
	}

	score := 1

	// all distances in millimeters
	// all times in milliseconds
	// all speeds in millimeters per millisecond
	a := 1 // mm/ms/ms

	// Iterative solution: Just try all the possibilities in order
	for raceNum, race := range races {
		raceSoltions := 0
		for holdTime := 1; holdTime < race.time; holdTime++ {
			v1 := a * holdTime
			timeLeft := race.time - holdTime
			distance := v1 * timeLeft
			if distance > race.distance {
				fmt.Printf("Found a solution for race %d. HoldTime: %d\n", raceNum+1, holdTime)
				raceSoltions++
			}
		}
		score *= raceSoltions
	}
	fmt.Printf("Score: %d\n", score)
}

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

	timeStr := strings.Join(timeStrings, "")
	distanceStr := strings.Join(distanceStrings, "")

	time, err := strconv.ParseInt(timeStr, 10, 64)
	if err != nil {
		panic(err)
	}

	distance, err := strconv.ParseInt(distanceStr, 10, 64)
	if err != nil {
		panic(err)
	}

	race := Race{
		time:     int(time),
		distance: int(distance),
	}

	// all distances in millimeters
	// all times in milliseconds
	// all speeds in millimeters per millisecond
	a := 1 // mm/ms/ms

	// Iterative solution: Just try all the possibilities in order
	raceSoltions := 0
	for holdTime := 1; holdTime < race.time; holdTime++ {
		v1 := a * holdTime
		timeLeft := race.time - holdTime
		distance := v1 * timeLeft
		if distance > race.distance {
			raceSoltions++
		}
	}
	fmt.Printf("Score: %d\n", raceSoltions)
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
)

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	total := 0
	regexes := map[*regexp.Regexp]int{
		regexp.MustCompile(`^one`):   1,
		regexp.MustCompile(`^two`):   2,
		regexp.MustCompile(`^three`): 3,
		regexp.MustCompile(`^four`):  4,
		regexp.MustCompile(`^five`):  5,
		regexp.MustCompile(`^six`):   6,
		regexp.MustCompile(`^seven`): 7,
		regexp.MustCompile(`^eight`): 8,
		regexp.MustCompile(`^nine`):  9,
	}

	for fileScanner.Scan() {
		line := fileScanner.Text()
		digits := []int{}

		for {
			if len(line) == 0 {
				break
			}
			matched := false
			for r, dig := range regexes {
				if r.MatchString(line) {
					line = line[1:]

					digits = append(digits, dig)
					matched = true
					break
				}
			}
			if !matched {
				c := line[0]
				if c >= '0' && c <= '9' {
					digits = append(digits, int(c-48))
					line = line[1:]
					matched = true
				}
			}
			if !matched {
				line = line[1:]
			}
		}

		fmt.Printf("Digits: %v\n", digits)
		code := 10*digits[0] + digits[len(digits)-1]
		fmt.Printf("Code: %d\n", code)
		total += code
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	fmt.Printf("Sum: %d\n", total)
}

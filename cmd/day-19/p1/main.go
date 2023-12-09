package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		line := fileScanner.Text()
		_ = line
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

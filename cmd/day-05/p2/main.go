package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
)

func AtoI(s string) int64 {
	val, _ := strconv.ParseInt(s, 10, 64)
	return val
}

type MapRule struct {
	start, end, diff int64
}

func NewMapRule(dst, src, len int64) MapRule {
	return MapRule{
		start: src,
		end:   src + len,
		diff:  dst - src,
	}
}

func (r MapRule) Map(in int64) (int64, bool) {

	if in >= r.start && in < r.end {
		return in + r.diff, true
	}
	return 0, false
}

type Mapper struct {
	ruleSet []MapRule
}

func (m *Mapper) Map(in int64) int64 {
	for _, rule := range m.ruleSet {
		if val, ok := rule.Map(in); ok {
			return val
		}
	}
	return in
}

var mapLineReg = regexp.MustCompile(`^([0-9]+) ([0-9]+) ([0-9]+)$`)

type ParseState int

const (
	parseStateEmpty ParseState = iota
	parseStateSeed2Soil
	parseStateSoil2Fert
	parseStateFert2Water
	parseStateWater2Light
	parseStateLight2Temp
	parseStateTemp2Humid
	parseStateHumid2Location
)

var parseStateMapperMap = map[ParseState]*Mapper{
	parseStateEmpty:          {},
	parseStateSeed2Soil:      {},
	parseStateSoil2Fert:      {},
	parseStateFert2Water:     {},
	parseStateWater2Light:    {},
	parseStateLight2Temp:     {},
	parseStateTemp2Humid:     {},
	parseStateHumid2Location: {},
}

type Range struct {
	start, end int64
}

type Work struct {
	r           Range
	resultsChan chan int64
}

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan()
	seedLine := fileScanner.Text()
	seedStrings := strings.Split(strings.Split(seedLine, ": ")[1], " ")
	ranges := []Range{}
	for i := 0; i < len(seedStrings); i += 2 {
		start, err := strconv.ParseInt(seedStrings[i], 10, 64)
		length, err := strconv.ParseInt(seedStrings[i+1], 10, 64)
		if err != nil {
			panic(err)
		}
		ranges = append(ranges, Range{start: start, end: start + length})
	}

	parseState := parseStateEmpty
	for fileScanner.Scan() {
		line := fileScanner.Text()
		if line == "" {
			parseState = parseStateEmpty
			continue
		} else if strings.HasPrefix(line, "seed-to-soil") {
			parseState = parseStateSeed2Soil
		} else if strings.HasPrefix(line, "soil-to-fert") {
			parseState = parseStateSoil2Fert
		} else if strings.HasPrefix(line, "fert") {
			parseState = parseStateFert2Water
		} else if strings.HasPrefix(line, "water") {
			parseState = parseStateWater2Light
		} else if strings.HasPrefix(line, "light") {
			parseState = parseStateLight2Temp
		} else if strings.HasPrefix(line, "temperature") {
			parseState = parseStateTemp2Humid
		} else if strings.HasPrefix(line, "humidity") {
			parseState = parseStateHumid2Location
		} else if mapLineMatch := mapLineReg.FindStringSubmatch(line); mapLineMatch != nil && len(mapLineMatch) == 4 {
			mapRule := NewMapRule(
				AtoI(mapLineMatch[1]),
				AtoI(mapLineMatch[2]),
				AtoI(mapLineMatch[3]),
			)
			mapper := parseStateMapperMap[parseState]
			mapper.ruleSet = append(mapper.ruleSet, mapRule)
		}
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	stages := []ParseState{
		parseStateEmpty,
		parseStateSeed2Soil,
		parseStateSoil2Fert,
		parseStateFert2Water,
		parseStateWater2Light,
		parseStateLight2Temp,
		parseStateTemp2Humid,
		parseStateHumid2Location,
	}
	resultsChans := []chan int64{}
	const batchSize = 10000
	for _, r := range ranges {
		for i := r.start; i < r.end; i += batchSize {
			length := min(batchSize, r.end-i)

			resultsChan := make(chan int64)
			resultsChans = append(resultsChans, resultsChan)
			work := Work{
				r: Range{
					start: i,
					end:   i + length,
				},
				resultsChan: resultsChan,
			}
			go func(w Work) {
				r := w.r
				localLowest := int64(math.MaxInt64)
				defer func() {
					w.resultsChan <- localLowest
					close(w.resultsChan)
				}()
				for i := r.start; i < r.end; i++ {
					x := i
					for _, stage := range stages {
						mapper := parseStateMapperMap[stage]
						x = mapper.Map(x)
					}
					if x < localLowest {
						localLowest = x
					}
				}
			}(work)
		}
	}

	lowests := []int64{}
	for _, resultChan := range resultsChans {
		lowests = append(lowests, <-resultChan)
	}

	slices.Sort(lowests)

	fmt.Printf("Lowest: %d\n", lowests[0])
}

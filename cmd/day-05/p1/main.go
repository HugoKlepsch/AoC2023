package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"regexp"
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

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	fileScanner.Scan()
	seedLine := fileScanner.Text()
	seedStrings := strings.Split(strings.Split(seedLine, ": ")[1], " ")
	seeds := []int64{}
	for _, seedString := range seedStrings {
		seed, err := strconv.ParseInt(seedString, 10, 64)
		if err != nil {
			panic(err)
		}
		seeds = append(seeds, seed)
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

	lowestLocation := int64(math.MaxInt64)
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
	for _, seed := range seeds {
		x := seed
		for _, stage := range stages {
			mapper := parseStateMapperMap[stage]
			x = mapper.Map(x)
		}
		if x < lowestLocation {
			lowestLocation = x
		}
	}
	fmt.Printf("Lowest: %d\n", lowestLocation)
}

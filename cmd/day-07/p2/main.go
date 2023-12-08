package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Card uint8

func (c Card) ToString() string {
	return fmt.Sprintf("%c", c)
}

var cardRanks = map[Card]int{
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'J': 1, // Joker
	'Q': 12,
	'K': 13,
	'A': 14,
}

type Hand struct {
	bid   int
	cards []Card
	// histogram counts the occurrences of each type of card
	histogram        map[Card]int
	reverseHistogram map[int][]Card
	highCard         Card
}

func (h Hand) ToString() string {
	strs := []string{}
	for _, card := range h.cards {
		strs = append(strs, card.ToString())
	}
	capabilities := []string{}
	if h.IsFiveOfAKind() {
		capabilities = append(capabilities, "5")
	}
	if h.IsFourOfAKind() {
		capabilities = append(capabilities, "4")
	}
	if h.IsFullHouse() {
		capabilities = append(capabilities, "F")
	}
	if h.IsThreeOfAKind() {
		capabilities = append(capabilities, "3")
	}
	if h.IsTwoPair() {
		capabilities = append(capabilities, "2")
	}
	if h.IsPair() {
		capabilities = append(capabilities, "P")
	}
	return fmt.Sprintf("[%s] [%s] %d", strings.Join(strs, ""), strings.Join(capabilities, ","), h.bid)
}

// Compare returns 0 when equal power, +ve when this hand is stronger than the other, and
// -ve when weaker than the other.
func (h Hand) Compare(other Hand) int {
	var mineIsType, otherIsType bool

	mineIsType = h.IsFiveOfAKind()
	otherIsType = other.IsFiveOfAKind()
	if mineIsType && !otherIsType {
		return 1
	} else if otherIsType && !mineIsType {
		return -1
	} else if mineIsType && otherIsType {
		return h.compareCardsInOrder(other)
	}

	mineIsType = h.IsFourOfAKind()
	otherIsType = other.IsFourOfAKind()
	if mineIsType && !otherIsType {
		return 1
	} else if otherIsType && !mineIsType {
		return -1
	} else if mineIsType && otherIsType {
		return h.compareCardsInOrder(other)
	}

	mineIsType = h.IsFullHouse()
	otherIsType = other.IsFullHouse()
	if mineIsType && !otherIsType {
		return 1
	} else if otherIsType && !mineIsType {
		return -1
	} else if mineIsType && otherIsType {
		return h.compareCardsInOrder(other)
	}

	mineIsType = h.IsThreeOfAKind()
	otherIsType = other.IsThreeOfAKind()
	if mineIsType && !otherIsType {
		return 1
	} else if otherIsType && !mineIsType {
		return -1
	} else if mineIsType && otherIsType {
		return h.compareCardsInOrder(other)
	}

	mineIsType = h.IsTwoPair()
	otherIsType = other.IsTwoPair()
	if mineIsType && !otherIsType {
		return 1
	} else if otherIsType && !mineIsType {
		return -1
	} else if mineIsType && otherIsType {
		return h.compareCardsInOrder(other)
	}

	mineIsType = h.IsPair()
	otherIsType = other.IsPair()
	if mineIsType && !otherIsType {
		return 1
	} else if otherIsType && !mineIsType {
		return -1
	} else if mineIsType && otherIsType {
		return h.compareCardsInOrder(other)
	}

	return h.compareCardsInOrder(other)
}

func (h Hand) compareCardsInOrder(other Hand) int {
	for i := 0; i < len(h.cards); i++ {

		var myPower, otherPower int

		myPower = cardRanks[h.cards[i]]
		otherPower = cardRanks[other.cards[i]]
		diff := myPower - otherPower
		if diff != 0 {
			return myPower - otherPower
		}
	}
	return 0
}

func (h Hand) IsFiveOfAKind() bool {
	highestCount := 0
	for card, count := range h.histogram {
		if card != 'J' && count > highestCount {
			highestCount = count
		}
	}
	jokerCount, ok := h.histogram['J']
	if !ok {
		jokerCount = 0
	}
	return highestCount+jokerCount >= 5
}

func (h Hand) IsFourOfAKind() bool {
	highestCount := 0
	for card, count := range h.histogram {
		if card != 'J' && count > highestCount {
			highestCount = count
		}
	}
	jokerCount, ok := h.histogram['J']
	if !ok {
		jokerCount = 0
	}
	return highestCount+jokerCount >= 4
}

func (h Hand) IsFullHouse() bool {
	highestCountA := 0
	highestCardA := Card('0')
	for card, count := range h.histogram {
		if card != 'J' && count > highestCountA {
			highestCountA = count
			highestCardA = card
		}
	}
	highestCountB := 0
	for card, count := range h.histogram {
		if card != 'J' && card != highestCardA && count > highestCountB {
			highestCountB = count
		}
	}
	jokerCount, ok := h.histogram['J']
	if !ok {
		jokerCount = 0
	}

	if highestCountA < 3 {
		diff := 3 - highestCountA
		highestCountA += diff
		jokerCount -= diff
	}

	if highestCountB < 2 {
		diff := 2 - highestCountB
		highestCountB += diff
		jokerCount -= diff
	}

	// This logic seems really gross but it's late
	return highestCountA >= 3 && highestCountB >= 2 && jokerCount >= 0
}

func (h Hand) IsThreeOfAKind() bool {
	highestCount := 0
	for card, count := range h.histogram {
		if card != 'J' && count > highestCount {
			highestCount = count
		}
	}
	jokerCount, ok := h.histogram['J']
	if !ok {
		jokerCount = 0
	}
	return highestCount+jokerCount >= 3
}

func (h Hand) IsTwoPair() bool {
	highestCountA := 0
	highestCardA := Card('0')
	for card, count := range h.histogram {
		if card != 'J' && count > highestCountA {
			highestCountA = count
			highestCardA = card
		}
	}
	highestCountB := 0
	for card, count := range h.histogram {
		if card != 'J' && card != highestCardA && count > highestCountB {
			highestCountB = count
		}
	}
	jokerCount, ok := h.histogram['J']
	if !ok {
		jokerCount = 0
	}

	if highestCountA < 2 {
		diff := 2 - highestCountA
		highestCountA += diff
		jokerCount -= diff
	}

	if highestCountB < 2 {
		diff := 2 - highestCountB
		highestCountB += diff
		jokerCount -= diff
	}

	// This logic seems really gross but it's late
	return highestCountA >= 2 && highestCountB >= 2 && jokerCount >= 0
}

func (h Hand) IsPair() bool {
	highestCount := 0
	for card, count := range h.histogram {
		if card != 'J' && count > highestCount {
			highestCount = count
		}
	}
	jokerCount, ok := h.histogram['J']
	if !ok {
		jokerCount = 0
	}
	return highestCount+jokerCount >= 2
}

func ParseHand(line string) Hand {
	handFields := strings.Fields(line)

	bid, err := strconv.ParseInt(handFields[1], 10, 64)
	if err != nil {
		panic(err)
	}

	cards := []Card{}
	histogram := map[Card]int{}
	reverseHistogram := map[int][]Card{}
	highestPower := 0
	highCard := Card('2')

	for i := 0; i < len(handFields[0]); i++ {
		card := Card(handFields[0][i])
		cards = append(cards, card)

		power := cardRanks[card]
		if power > highestPower {
			highCard = card
		}

		prev, ok := histogram[card]
		if !ok {
			prev = 1
		} else {
			prev += 1
		}
		histogram[card] = prev
	}

	for card, count := range histogram {
		val, ok := reverseHistogram[count]
		if !ok {
			val = []Card{}
		}
		val = append(val, card)
		reverseHistogram[count] = val
	}

	h := Hand{
		bid:              int(bid),
		cards:            cards,
		histogram:        histogram,
		reverseHistogram: reverseHistogram,
		highCard:         highCard,
	}
	return h
}

func main() {

	var err error
	fileScanner := bufio.NewScanner(os.Stdin)
	fileScanner.Split(bufio.ScanLines)

	hands := []Hand{}
	for fileScanner.Scan() {
		line := fileScanner.Text()
		hands = append(hands, ParseHand(line))
	}
	if err = fileScanner.Err(); err != nil {
		fmt.Println(err)
		panic(err)
	}

	sort.Slice(hands, func(i, j int) bool { return hands[i].Compare(hands[j]) < 0 })

	score := 0
	for rank, hand := range hands {
		handScore := (rank + 1) * hand.bid
		score += handScore
		fmt.Printf("Hand: %s: %d\n", hand.ToString(), handScore)
	}

	fmt.Printf("Score: %d\n", score)
}

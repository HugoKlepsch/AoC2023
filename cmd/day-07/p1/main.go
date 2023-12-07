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
	'J': 11,
	'Q': 12,
	'K': 13,
	'A': 14,
}

type Hand struct {
	bid   int
	cards []Card
	// histogram counts the occurrences of each type of card
	histogram map[Card]int
	highCard  Card
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
	return len(h.histogram) == 1
}

func (h Hand) IsFourOfAKind() bool {
	hasFour := false
	for _, count := range h.histogram {
		if count == 4 {
			hasFour = true
			break
		}
	}
	return hasFour
}

func (h Hand) IsFullHouse() bool {
	hasThree := false
	hasTwo := false
	for _, count := range h.histogram {
		if count == 3 {
			hasThree = true
		} else if count == 2 {
			hasTwo = true
		}
		if hasThree && hasTwo {
			break
		}
	}
	return hasTwo && hasThree
}

func (h Hand) IsThreeOfAKind() bool {
	hasThree := false
	for _, count := range h.histogram {
		if count == 3 {
			hasThree = true
			break
		}
	}
	return hasThree
}

func (h Hand) IsTwoPair() bool {
	hasTwoA := false
	hasTwoB := false
	for _, count := range h.histogram {
		if count == 2 && !hasTwoA {
			hasTwoA = true
		} else if count == 2 && hasTwoA {
			hasTwoB = true
		}
		if hasTwoA && hasTwoB {
			break
		}
	}
	return hasTwoA && hasTwoB
}

func (h Hand) IsPair() bool {
	hasTwo := false
	for _, count := range h.histogram {
		if count == 2 {
			hasTwo = true
			break
		}
	}
	return hasTwo
}

func ParseHand(line string) Hand {
	handFields := strings.Fields(line)

	bid, err := strconv.ParseInt(handFields[1], 10, 64)
	if err != nil {
		panic(err)
	}

	cards := []Card{}
	histogram := map[Card]int{}
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

	h := Hand{
		bid:       int(bid),
		cards:     cards,
		histogram: histogram,
		highCard:  highCard,
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

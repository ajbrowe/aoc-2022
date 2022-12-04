package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Section interface {
	Contains(SectionRange)
	Overlaps(SectionRange)
	Label()
	Print()
}

type SectionRange struct {
	start int
	end int
}

func (s SectionRange) Contains(partner SectionRange) bool {
	return s.start <= partner.start && s.end >= partner.end
}
func (s SectionRange) Overlaps(partner SectionRange) bool {
	if s.start > partner.end || s.end < partner.start {
		return false
	}
	return true
}


func (s SectionRange) Label() string {
	return fmt.Sprintf("Start: %d, End: %d", s.start, s.end)
}

func (s SectionRange) Print() {
	fmt.Println(s.Label())
}

func main() {
	file, err := os.Open("sample")
	//file, err := os.Open("input")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	consumedPairs := 0
	overlappingPairs := 0
	for scanner.Scan() {
		value := scanner.Text()
		var sections []SectionRange
		pair := strings.Split(value, ",")
		for _, rangeString := range pair {
			sections = append(sections, parseRange(rangeString))
		}
		if sections[0].Contains(sections[1]) {
			consumedPairs++
		} else if sections[1].Contains(sections[0]) {
			consumedPairs++
		}
		if sections[0].Overlaps(sections[1]) {
			overlappingPairs++
		}

	}

	fmt.Printf("Consumed Pairs - Part 1: %d\n", consumedPairs)
	fmt.Printf("Overlapping Pairs - Part 2: %d\n", overlappingPairs)

}

func parseRange(rangeString string) SectionRange {
	var rangeLimit []int
	rangeLimitStrings := strings.Split(rangeString, "-")
	for _, limit := range rangeLimitStrings {
		limitValue, err := strconv.Atoi(limit)
		if err != nil {
			log.Fatal(err)
		}
		rangeLimit = append(rangeLimit, limitValue)
	}
	section :=SectionRange{rangeLimit[0], rangeLimit[1]}
	return section
}

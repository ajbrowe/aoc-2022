package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

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

	sum_of_priorities := 0

	var rucksacks []string
	for scanner.Scan() {
		value := scanner.Text()
		rucksacks = append(rucksacks,value);
		number_of_items := len(value);
		part1 := value[0:(number_of_items/2)]
		part2 := value[(number_of_items/2):(number_of_items)]
		duplicate := findDuplicate(part1,part2)
		priority := getValue(duplicate)
		sum_of_priorities += priority
	}

	badgeScore := 0
	for i := 0; i < len(rucksacks); i+=3 {
		dup_string := findDuplicate(rucksacks[i],rucksacks[i+1])
		badge := findDuplicate(dup_string,rucksacks[i+2])
		if len(badge) == 0 {
			log.Fatalf("No duplicates found in %d, %s",i, strings.Join(rucksacks[i:i+2],", "))
		}
		if len(badge) > 1 {
			log.Fatalf("Found multiple badges %d, %s\n", i, badge)
		}
		badgeScore += getValue(badge)

	}
	fmt.Printf("part1 sum: %d\n",sum_of_priorities)
	fmt.Printf("part2 sum: %d\n", badgeScore)


}


func findDuplicate(part1 string, part2 string) string {
	matches := make(map[string]bool)
	dupes := make(map[string]bool)
	for _, character := range part1 {
		letter := string(character)
		if _, exists := matches[letter]; !exists {
			matches[letter] = true
		}
	}
	var duplicates []string
	for _, character := range part2 {
		letter := string(character)
		if _, present := matches[letter]; present {

			// Don't duplicate dupes !
			if _, exists := dupes[letter]; !exists {
				duplicates = append(duplicates, letter)
				dupes[letter] = true
			}
		}
	}

	return strings.Join(duplicates,"")
}

func getValue(letter string) int {
	characters := []rune(letter)
	asciiValue := int(characters[0]);
	if (asciiValue >= 97) {
		return asciiValue - 96;
	}
	return asciiValue - 38;
}


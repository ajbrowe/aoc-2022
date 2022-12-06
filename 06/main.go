package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
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

	for scanner.Scan() {
		value := scanner.Text()
		processSignal(value)
		processSignal2(value)
	}
}

func processSignal(signal string) {
	signalLength := len(signal)
	for i := 4; i < signalLength; i++ {
		if !hasRepeat(signal[i-4:i]) {
			fmt.Printf("Part 1 Found signal at %d\n",i)
			break
		}
	}
}

func processSignal2(signal string) {
	signalLength := len(signal)
	for i := 14; i < signalLength; i++ {
		if !hasRepeat(signal[i-14:i]) {
			fmt.Printf("Part 2 Found signal at %d\n",i)
			break
		}
	}
}


func hasRepeat(signal string) bool {
	counts := make(map[string]int)
	repeatFound := false;
	for _, character := range signal {
		letter := string(character)
		if _, exists := counts[letter]; !exists {
			counts[letter] = 1
		} else {
			counts[letter]++
			repeatFound = true
		}
	}
	return repeatFound
}



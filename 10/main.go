package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
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

	part2Crt := make([]bool,240)
	part1Sum := 0
	register := 1
	cycle := 1
	for scanner.Scan() {
		value := scanner.Text()

		increment := 0

		if value != "noop" {
			fields := strings.Fields(value)
			regIncrement, err := strconv.Atoi(fields[1])
			if err != nil {
				log.Fatal(err)
			}
			increment = regIncrement
			part1Sum += getCycleSum(cycle, register, value)
			renderSprite(part2Crt, cycle, register)

			// First cycle complete
			cycle++
		}
		part1Sum += getCycleSum(cycle, register, value)
		renderSprite(part2Crt, cycle, register)

		//Increment register AFTER second cycle
		register += increment

		cycle++
	}

	fmt.Printf("Part1 sum: %d\n", part1Sum)
	fmt.Println("Part2 CRT")
	for y := 0; y < 6; y++ {
		var line strings.Builder
		startPosition := y * 40
		endPosition := startPosition + 40
		for _, p := range part2Crt[startPosition:endPosition] {
			pixelChar := '.'
			if p {
				pixelChar = '#'
			}
			line.WriteRune(pixelChar)
		}
		fmt.Println(line.String())
	}
}

func getCycleSum(cycle int, register int, value string) int {
	sum := 0
	if cycle > 0 && cycle <= 220 && ((cycle - 20 ) % 40 == 0) {
		fmt.Printf("Cycle %d: X = %d, Signal Strength: %d, [%s]\n", cycle, register, (cycle * register), value)
		sum = (cycle * register)
	}
	return sum
}


func renderSprite(crt []bool, cycle int, register int) {
	crtPosition := (cycle - 1) % 240
	spritePosition := register % 40
	crtHorizontal := crtPosition % 40
	pixel := false
	if (crtHorizontal >= spritePosition -1 && crtHorizontal <= spritePosition + 1) {
		pixel = true
	}
	crt[crtPosition] = pixel
}


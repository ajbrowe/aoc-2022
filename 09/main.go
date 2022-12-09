package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Pos struct {
	x int
	y int
}

func (p *Pos) asString() string {
	return fmt.Sprintf("%d:%d", p.x, p.y)
}

func (tail *Pos) Follow(head *Pos) {
	tailX := tail.x
	tailY := tail.y

	xDistance := tailX - head.x
	yDistance := tailY - head.y

	if (Abs(xDistance) > 1) || (Abs(yDistance) > 1) {
		xMove := 0
		if xDistance > 0 {
			xMove = 1
		}
		if xDistance < 0 {
			xMove = -1
		}
		yMove := 0
		if yDistance > 0 {
			yMove = 1
		}
		if yDistance < 0 {
			yMove = -1
		}
		tailX -= xMove
		tailY -= yMove
	}
	tail.x = tailX
	tail.y = tailY
}

func main() {
	file, err := os.Open("sample")
	//file, err := os.Open("sample2")
	//file, err := os.Open("input")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	tailMoves := make(map[string]int)
	tailMoves2 := make(map[string]int)

	part1Knots := []*Pos{&Pos{500, 500},&Pos{500, 500}}
	var part2Knots []*Pos
	for i := 0; i < 10; i++ {
		part2Knots = append(part2Knots, &Pos{500, 500})
	}

	tailMoves[part1Knots[1].asString()] = 1
	tailMoves2[part2Knots[9].asString()] = 1

	for scanner.Scan() {
		value := scanner.Text()
		direction, distance := parseInstruction(value)
		runMoves(direction, distance, part1Knots, tailMoves)
		runMoves(direction, distance, part2Knots, tailMoves2)
	}
	fmt.Printf("Part 1: Tail moved %d times\n", len(tailMoves))
	fmt.Printf("Part 2: Tail moved %d times\n", len(tailMoves2))
}

func parseInstruction(instruction string) (string, int) {
	fields := strings.Fields(instruction)
	distance, err := strconv.Atoi(fields[1])
	if err != nil {
		log.Fatal(err)
	}
	return fields[0], distance
}

func runMoves(direction string, distance int, knots []*Pos, tailMoves map[string]int) {
	numberOfKnots := len(knots)

	head := knots[0]
	for d := 0; d < distance; d++ {
		headX := head.x
		headY := head.y
		if direction == "U" {
			headY--
		}
		if direction == "R" {
			headX++
		}
		if direction == "L" {
			headX--
		}
		if direction == "D" {
			headY++
		}
		head.x = headX
		head.y = headY
		for n := 1; n < numberOfKnots; n++ {
			knot := knots[n]
			fmt.Println(n, knot)
			if (n == numberOfKnots-1) {
				tailX := knot.x
				tailY := knot.y

				not.Follow(knots[n-1])

				if (tailX != knot.x) || (tailY != knot.y) {
					location := knot.asString()
					if _, exists := tailMoves[location]; exists {
						tailMoves[location]++
					} else {
						tailMoves[location] = 1
					}
				}
			} else {
				knot.Follow(knots[n-1])
			}
		}
	}

}

func Abs(val int) int {
	if val < 0 {
		return -val
	}
	return val
}

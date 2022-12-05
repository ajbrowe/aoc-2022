package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"
)


func main() {
	//file, err := os.Open("sample")
	file, err := os.Open("input")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var stackDescriptions []string
	var moves []string
	finishedStack := false
	for scanner.Scan() {
		value := scanner.Text()

		if len(value) == 0 {
			finishedStack = true
			continue
		}

		if finishedStack {
			moves = append(moves, value)
		} else {
			stackDescriptions = append(stackDescriptions, value)
		}
	}

	stacks := buildStacks(stackDescriptions)
	stacks2 := buildStacks(stackDescriptions)
	stacks = runMoves(stacks, moves)
	stacks2 = runMoves2(stacks2, moves)
	fmt.Println(stacks)
	fmt.Println(stacks2)
}

func buildStacks(stackDescriptions []string) map[string][]string {
	stacks := make(map[string][]string)
	for _, line := range stackDescriptions {
		for i := 0; i < len(line); i+=4 {
			entry := strings.TrimSpace(line[i:i+3])
			stackLabel := fmt.Sprintf("%d", (i / 4) + 1)
			if len(entry) == 0 {
				continue
			}
			re := regexp.MustCompile(`\d+`)
			if re.MatchString(entry) {
				fmt.Printf("Skipping stack labels %s\n", entry)
				continue
			}
			stacks[stackLabel] = append(stacks[stackLabel], entry)
		}
	}
	return stacks
}

func runMoves(stacks map[string][]string, moves []string) map[string][]string {
	for _, moveString := range moves {
		move := parseMoveString(moveString)
		fmt.Println(move)
		stacks = moveStack(stacks, move)
	}

	return stacks
}
func runMoves2(stacks map[string][]string, moves []string) map[string][]string {
	for _, moveString := range moves {
		move := parseMoveString(moveString)
		fmt.Println(move)
		stacks = moveStack2(stacks, move)
	}

	return stacks
}
type Move struct {
	number int
	from string
	to string
}

func parseMoveString(move string) Move {
	re := regexp.MustCompile(`move (\d+) from (\d+) to (\d+)`)
	moveData := re.FindAllStringSubmatch(move, -1)
	if moveData == nil {
		log.Fatalf("Move data '%s' doesn't match regexp", move)
	}
	moveNumber, err := strconv.Atoi(moveData[0][1])
	if err != nil {
		log.Fatal(err)
	}
	return Move{moveNumber, moveData[0][2], moveData[0][3]}
}

func moveStack(stacks map[string][]string, move Move) map[string][]string {
	for i := 0; i < move.number; i++ {
		entry := stacks[move.from][0]
		stacks[move.from] = stacks[move.from][1:]
	    stacks[move.to] = append([]string{entry}, stacks[move.to]...)
	}
	return stacks
}

func moveStack2(stacks map[string][]string, move Move) map[string][]string {
	fromStack := stacks[move.from]
	toStack := stacks[move.to]
	count := move.number

	// get first number of entries on FROM stack
	entries := fromStack[0:count]

	fromStackCopy := make([]string, (len(fromStack) - count))
	copy(fromStackCopy, fromStack[count:])
	// Add first number of entries from FROM stack onto start of TO stack
	toStack = append(entries, toStack...)

	stacks[move.from] = fromStackCopy
	stacks[move.to] = toStack
	return stacks
}

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

	strategy_one_total_score := 0
	strategy_two_total_score := 0
	for scanner.Scan() {
		choices := strings.Fields(scanner.Text())
		strategy_one_total_score += score1(choices)
		strategy_two_total_score += score2(choices)
	}

	fmt.Printf("Total score: strategy1 %d, strategy2 %d\n", strategy_one_total_score, strategy_two_total_score)
}

//A Rock beats Z Sissors
//B Paper beats X Rock
//C Sissors beats Y Paper
//X Rock
//Y Paper
//Z Sissors
func score1(choices[] string) int {
	outcome_map := map[string]map[string]int{
		"A": map[string]int{
			"X": 3, //rock rock draw
			"Y": 6, //rock paper win
			"Z": 0, //rock sissors loss
		},
		"B": map[string]int{
			"X": 0, //paper:rock loss
			"Y": 3, //paper:paper draw
			"Z": 6, //paper:sissors win
		},
		"C": map[string]int{
			"X": 6, //sissors:rock win
			"Y": 0, //sissors:paper loss
			"Z": 3, //sissors:sissors draw
		},
	}

	shape_score_map := map[string]int{
		"X": 1,
		"Y": 2,
		"Z": 3,
	}

	shape_score := shape_score_map[choices[1]]
	outcome_score := outcome_map[choices[0]][choices[1]]
	return shape_score + outcome_score
}

//A Rock beats Z Sissors
//B Paper beats X Rock
//C Sissors beats Y Paper
//X Lose
//Y Draw
//Z Win
func score2(choices[] string) int {
	strategy_choice_map := map[string]map[string]string{
		"A": map[string]string{
			"X": "C",
			"Y": "A",
			"Z": "B",
		},
		"B": map[string]string{
			"X": "A",
			"Y": "B",
			"Z": "C",
		},
		"C": map[string]string{
			"X": "B",
			"Y": "C",
			"Z": "A",
		},
	}

	outcome_score_map := map[string]int{
		"X": 0,
		"Y": 3,
		"Z": 6,
	}
	shape_score_map := map[string]int{
		"A": 1,
		"B": 2,
		"C": 3,
	}
	strategy_choice := strategy_choice_map[choices[0]][choices[1]]
	shape_score := shape_score_map[strategy_choice]
	outcome_score := outcome_score_map[choices[1]]
	return shape_score + outcome_score
}

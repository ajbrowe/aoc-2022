package main

import (
	"bufio"
	"fmt"
	"os"
	"log"
	"strconv"
	"sort"
)

func main() {
	//file, err := os.Open("sample.txt")
	file, err := os.Open("input")

	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	var calories []int
	var current_sum int

	i := 0
	for scanner.Scan() {
		value_string := scanner.Text()
		//fmt.Println(value_string)
		if len(value_string) == 0 {
			calories = append(calories, current_sum)
			fmt.Printf("Elf %d has %d\n", i, current_sum)
			i++
			current_sum = 0
			continue
		}
		value, err := strconv.Atoi(value_string)
		if err != nil {
			log.Fatal(err)
		}

		current_sum += value

	}
	calories = append(calories, current_sum)
	fmt.Printf("Elf %d has %d\n", i, current_sum)

	caloriesSlice := calories [:]
	sort.Sort(sort.Reverse(sort.IntSlice(caloriesSlice)))
	fmt.Printf("Most Calories %d\n",calories[0])

	sum := calories[0] + calories[1] + calories[2];
	fmt.Printf("Sum of top three Elves %d\n", sum)
}

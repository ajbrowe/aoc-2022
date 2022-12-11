package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type worryOperation struct{
	operator string
	operand string
}
type Monkey struct {
	items []int
	operation *worryOperation
	test int
	testConditions map[string]int
	inspectionCount int
}

func (m *Monkey) hasItem() bool {
	return len(m.items) > 0
}

func (m *Monkey) Inspect(divide bool, lcm int) (int, int) {
	worry := m.items[0];
	m.items = m.items[1:]

	//fmt.Printf("  Monkey inspects an item with a worry level of %d\n", worry)
	worry = m.Operation(worry)
	if divide {
		worry = worry / 3
		//fmt.Printf("    Monkey gets bored with item. Worry level is divided by 3 to %d\n", worry)
	} else {
		// Reduce worry to the lowest common modulo to prevent integer overflow
		worry = worry % lcm
	}
	var destination int
	if (worry % m.test == 0) {
		//fmt.Printf("    Current worry level is divisible by %d\n", m.test)
		destination = m.testConditions["true"]
	} else {
		//fmt.Printf("    Current worry level is not divisible by %d\n", m.test)
		destination = m.testConditions["false"]
	}

	m.inspectionCount++
	return worry, destination
}

func (m *Monkey) Operation(worry int) int{
	op := m.operation
	var operand int;
	if (op.operand == "old") {
		operand = worry
	} else {
		value, err := strconv.Atoi(op.operand)
		if err != nil {
			log.Fatal(err)
		}
		operand = value
	}
	if (op.operator == "*") {
		worry = worry * operand
		//fmt.Printf("    Worry level is multimplied by %d to %d\n", operand, worry)
	}
	if (op.operator == "+") {
		worry = worry + operand
		//fmt.Printf("    Worry level increases by %d to %d\n", operand, worry)
	}
	return worry 

}

func (m *Monkey) receiveItem(worry int) {
	m.items = append(m.items, worry)
}

func (m *Monkey) showWorryLevels(id int) {
	var itemStrings []string
	for _, v := range m.items { itemStrings = append(itemStrings, fmt.Sprintf("%d",v)) }
	fmt.Printf("Monkey %d: %s\n",id, strings.Join(itemStrings, ", "))
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

	var monkeys []*Monkey
	var monkeys2 []*Monkey

	var gatheredInfo []string
	for scanner.Scan() {
		value := scanner.Text()
		if value == "" {
			monkeys = append(monkeys, ParseMonkeyInfo(gatheredInfo))
			monkeys2 = append(monkeys2, ParseMonkeyInfo(gatheredInfo))
			gatheredInfo = []string{}
		} else {
			gatheredInfo = append(gatheredInfo, strings.TrimSpace(value))
		}
	}
	monkeys = append(monkeys, ParseMonkeyInfo(gatheredInfo))
	monkeys2 = append(monkeys2, ParseMonkeyInfo(gatheredInfo))
	gatheredInfo = []string{}

	for i := 0; i < 20; i++ {
		monkeyRun(monkeys, true, 1)
	}

	var counts []int
	for i, m := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n", i, m.inspectionCount)
		counts = append(counts, m.inspectionCount)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	fmt.Printf("Part1 Monkey Business = %d\n", (counts[0] * counts[1]))

	lcm := 1
	for _, m := range monkeys2 {
		lcm = lcm * m.test
	}

	for i := 0; i < 10000; i++ {
		round := i + 1;
		monkeyRun(monkeys2, false, lcm)
		if round == 1 || round == 20 || round % 1000 == 0 {
			//fmt.Printf("\n== After round %d ==\n", round)
			//reportWorry(monkeys2)
			//reportInspections(monkeys2)
		}
	}

	// reset counts
	counts = []int{}
	for i, m := range monkeys2 {
		fmt.Printf("Monkey %d inspected items %d times.\n", i, m.inspectionCount)
		counts = append(counts, m.inspectionCount)
	}
	sort.Sort(sort.Reverse(sort.IntSlice(counts)))
	fmt.Printf("Part2 Monkey Business = %d\n", (counts[0] * counts[1]))

}

func ParseMonkeyInfo (gatheredInfo []string) *Monkey {
	var items []int
	itemRe := regexp.MustCompile(`Starting items: ([\d ,]+)`)

	var operation *worryOperation
	operationRe := regexp.MustCompile(`Operation: new = old (.*)`)

	var test int
	testRe := regexp.MustCompile(`Test: divisible by (\d+)`)

	conditions := make(map[string]int)
	conditionsRe := regexp.MustCompile(`If (true|false): throw to monkey (\d+)`)

	for _, info := range gatheredInfo {
		if matches := itemRe.FindAllStringSubmatch(info,-1); matches != nil {
			itemList := matches[0][1]
			itemsStr := strings.Split(itemList, ", ")
			for _, v := range itemsStr {
				value, err := strconv.Atoi(v)
				if err != nil {
					log.Fatal(err)
				}
				items = append(items, value)
			}
		}

		if matches := operationRe.FindAllStringSubmatch(info,-1); matches != nil {
			fields := strings.Fields(matches[0][1])
			operation = &worryOperation{
				operator: fields[0],
				operand: fields[1],
			}
		}

		if matches := testRe.FindAllStringSubmatch(info,-1); matches != nil {
			value, err := strconv.Atoi(matches[0][1])
			if err != nil {
				log.Fatal(err)
			}
			test = value
		}

		if matches := conditionsRe.FindAllStringSubmatch(info,-1); matches != nil {
			value, err := strconv.Atoi(matches[0][2])
			if err != nil {
				log.Fatal(err)
			}
			conditions[matches[0][1]] = value
		}
	}
	return &Monkey{
		items: items,
		operation: operation,
		test: test,
		testConditions: conditions,
		inspectionCount: 0,
	}
}

func monkeyRun(monkeys []*Monkey, divide bool, lcm int) {
	for _, monkey := range monkeys {
		//fmt.Printf("Monkey %d:\n", n)
		for monkey.hasItem() {
			worry, destination := monkey.Inspect(divide, lcm)
			monkeys[destination].receiveItem(worry)
			//fmt.Printf("    Item with worry level %d is thrown to monkey %d.\n", worry, destination)
		}
	}
}

func reportInspections( monkeys []*Monkey) {
	for i, m := range monkeys {
		fmt.Printf("Monkey %d inspected items %d times.\n", i, m.inspectionCount)
	}
}

func reportWorry( monkeys []*Monkey) {
	for i, m := range monkeys {	m.showWorryLevels(i) }
}

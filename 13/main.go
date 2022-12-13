package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)
type ListNode struct {
	children []*ListNode
	value int
	valueSet bool
	level int
}

func (l *ListNode) hasChildren() bool {
	return len(l.children) > 0 || !l.valueSet
}

func (l *ListNode) isValue() bool {
	return l.valueSet
}

func (l *ListNode) isEmpty() bool {
	return !(l.hasChildren() || l.isValue())
}

func (l *ListNode) Label() string {
	if l.isValue() {
		return fmt.Sprintf("%d", l.value)
	}
	var sb strings.Builder
	sb.WriteString("[")
	var subLabels []string
	for _, c := range l.children {
		subLabels = append(subLabels, c.Label())
	}
	if len(subLabels) > 0 {
		sb.WriteString(strings.Join(subLabels, ","))
	}
	sb.WriteString("]")
	return sb.String()
}

func (l *ListNode) Prefix() string {
	
	var spaces []rune
	for i:= 0; i < l.level; i++ {
		spaces = append(spaces, ' ',' ')
	}

	return fmt.Sprintf("%s-",string(spaces))
}

func (l *ListNode) Compare(r *ListNode) int {
	if (l.isValue() && r.isValue()) {
		//fmt.Printf("%s Compare %d vs %d\n",l.Prefix(), l.value, r.value)
		return l.value - r.value
	}
	if (l.hasChildren() && r.hasChildren()) {
		rLen := len(r.children)
		lLen := len(l.children)
		//fmt.Printf("%s Compare %s vs %s\n", l.Prefix(), l.Label(), r.Label())
		if lLen == 0 || rLen == 0 {
			return lLen - rLen
		}
		
		for i, c := range l.children {
			if rLen < (i + 1) {
				//fmt.Printf("%s Right side ran out of items so inputs are not in the right order\n", l.Prefix())
				return lLen - rLen
			}
			result := c.Compare(r.children[i])
			if result != 0 {
				return result
			}
		}
		if rLen > lLen {
			//fmt.Printf("%s Left side ran out of items so inputs are in the right order\n", l.Prefix())
			return lLen - rLen
		}
		if rLen < lLen {
			//fmt.Printf("%s Right side ran out of items so inputs are not in the right order\n", l.Prefix())
			return lLen - rLen
		}
		return 0
	}
	if (l.isValue() && r.hasChildren()) {
		var children []*ListNode
		children = append(children, l)
		newL := &ListNode{
			children: children,
			level: l.level+1,
		}
		return newL.Compare(r)
	}
	if (l.hasChildren() && r.isValue()) {
		var children []*ListNode
		children = append(children, r)
		newR := &ListNode{
			children: children,
			level: r.level+1,
		}
		return l.Compare(newR)
	}
	return 0
}

type Pair struct {
	left *ListNode
	right *ListNode
}

type SignalPackets []*ListNode

func (s SignalPackets) Len() int {
	return len(s)
}
func (s SignalPackets) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SignalPackets) Less(i, j int) bool {
	return (s[i].Compare(s[j]) < 0)
}


func (p *Pair) Compare() bool {
	l := p.left
	r := p.right
	result:= l.Compare(r)
	return result <= 0
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

	var rawPairs [][]string
	var pair []string

	for scanner.Scan() {
		value := scanner.Text()
		if (value == "") {
			rawPairs = append(rawPairs, pair)
			pair = []string{}
		} else {
			pair = append(pair, value)
		}
	}
	rawPairs = append(rawPairs, pair)


	var Pairs []*Pair
	for _, p := range rawPairs {
		Pairs = append(Pairs, &Pair{
			left: parseListString(p[0]),
			right: parseListString(p[1]),
		})
	}

	part1Sum := 0
	part2Set := SignalPackets{parseListString("[[2]]"), parseListString("[[6]]")}
	for i, p := range Pairs {
		//fmt.Printf("- Compare %s vs %s\n", p.left.Label(), p.right.Label())
		if p.Compare() {
			pairNumber := i + 1
			//fmt.Printf("Pair %d are in order\n", pairNumber)
			part1Sum += pairNumber
		} else {
			//fmt.Printf("Pair %d are not in order\n", i + 1)
		}

		part2Set = append(part2Set, p.left, p.right)
	}

	fmt.Printf("Part1 Sum: %d\n", part1Sum)

	sort.Sort(part2Set)

	part2Product := 1
	for i, n := range part2Set {
		if n.Label() == "[[2]]" {
			part2Product = part2Product * (i + 1)
		}
		if n.Label() == "[[6]]" {
			part2Product = part2Product * (i + 1)
		}
	}

	fmt.Printf("Part2 Product: %d\n", part2Product)

}


func parseListString(listString string) *ListNode {
	var stack []*ListNode
	var current *ListNode
	var nodeValue []rune
	for _, c := range listString {
		if c == '[' {
			if current != nil {
				stack = append(stack, current)
			}
			current = &ListNode{level: len(stack)}
		}
		if c >= '0' && c <= '9' {
			nodeValue = append(nodeValue, c)
		}
		if (c == ','|| c == ']') && len(nodeValue) > 0 {
			value, err := strconv.Atoi(string(nodeValue))
			if err != nil {
				log.Fatal(err)
			}
			nodeValue = []rune{}
			current.children = append(current.children, &ListNode{value: value, valueSet: true, level: len(stack) + 1})
		}
		if c == ']' {
			if len(stack) > 0 {
				endOfStackIndex := len(stack) - 1
				parent := stack[endOfStackIndex]
				stack = stack[:endOfStackIndex]
				parent.children = append(parent.children, current)
				current = parent
			}
		}
	}
	return current
}

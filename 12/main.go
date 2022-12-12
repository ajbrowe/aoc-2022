package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
)

type Pos struct {
	x int
	y int
}

type Node struct {
	pos *Pos
	height int
	priority int
	parent *Pos
}

type Queue []*Node

func (p *Pos) Key() string {
	return fmt.Sprintf("%d-%d", p.x, p.y)
}

func (p *Pos) Neighbours(widthLimit int, heightLimit int) []*Pos {
	var neighbours []*Pos
	if p.x > 0 {
		neighbours = append(neighbours, &Pos{p.x - 1, p.y})
	}
	if p.x < (widthLimit - 1) {
		neighbours = append(neighbours, &Pos{p.x + 1, p.y})
	}
	if p.y > 0 {
		neighbours = append(neighbours, &Pos{p.x, p.y - 1})
	}
	if p.y < (heightLimit - 1) {
		neighbours = append(neighbours, &Pos{p.x, p.y + 1})
	}
	return neighbours
}


func (n *Node) PossibleNeighbourNodes(terrain [][]int) []*Node {
	heightLimit := len(terrain)
	widthLimit := len(terrain[0])
	neighbours := n.pos.Neighbours(widthLimit, heightLimit)

	var possibleNeighbours []*Node
	for _, p := range neighbours {
		neighbourHeight := terrain[p.y][p.x]
		if neighbourHeight - n.height <= 1 {
			node := &Node{
				pos: p,
				height: neighbourHeight,
				parent: n.pos,
			}
			possibleNeighbours = append(possibleNeighbours, node)
		}
	}
	return possibleNeighbours

}

func (p *Pos) ManhattanDistance( goal *Pos ) int {
	return Abs(p.x - goal.x) + Abs(p.y - goal.y)
}

func (p *Pos) HeightDifference (destination *Pos, terrain [][]int) int {
	return terrain[destination.y][destination.x] - terrain[p.y][p.x]
}
func (p *Pos) Equals(destination *Pos) bool {
	return p.y == destination.y && p.x == destination.x
}

func (n *Node) Key() string {
	return n.pos.Key()
}

func (q Queue) Len() int {
	return len(q)
}

func (q Queue) Swap (i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q Queue) Less(i, j int) bool {
	return q[i].priority < q[j].priority
}

func (q Queue) ContainsLowerH(pos *Pos, h int) bool {
	for _, n := range q {
		if n.pos.Equals(pos) && n.priority < h {
			return true
		}
	}
	return false
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

	var terrain [][]int
	var start *Pos
	var end *Pos
	y := 0
	for scanner.Scan() {
		var row []int
		value := scanner.Text()
		for x, l := range value {
			if (l == 'S') {
				start = &Pos{x,y}
				l = 'a'
			}
			if (l == 'E') {
				end = &Pos{x,y}

				l = 'z'
			}
			row = append(row, getValue(l))
		}
		terrain = append(terrain, row)
		y++
	}

	path := path(start, end, terrain)
	fmt.Printf("Part1 method call had %d steps\n", path)

	part2Count := path2(end, terrain)
	fmt.Printf("Part2 method call had %d steps\n", part2Count)


}

func path (start, end *Pos, terrain [][]int) int {
	closedList := make(map[string]*Node)
	startNode := &Node{
		pos: start,
		height : 1,
		priority: 0,
	}
	var endNode *Node
	openList := []*Node{startNode}
	cost := map[string]int{ start.Key() : 0}
	cameFrom := map[string]*Node{ start.Key() : startNode}

	for len(openList) > 0 {
		sort.Sort(Queue(openList))
		current := openList[0]
		openList = openList[1:]
		if current.pos.Equals(end) {
			endNode = current
			break
		}

		nextNeighbours := current.PossibleNeighbourNodes(terrain)
		for _, neighbour := range nextNeighbours {
			newCost := cost[current.Key()] + 1
			if _, exists := cost[neighbour.Key()]; !exists || newCost < cost[neighbour.Key()] {
				cost[neighbour.Key()] = newCost
				neighbour.priority = newCost + neighbour.pos.HeightDifference(end, terrain)
				cameFrom[neighbour.Key()] = current
				openList = append(openList, neighbour)
			}
		}
		closedList[current.Key()] = current
	}
	steps := 0;
	parent := endNode.parent
	for parent != nil {
		steps++
		parentNode := closedList[parent.Key()]
		parent = parentNode.parent
	}

	return steps 
}


func path2 (end *Pos, terrain [][]int) int {
	closedList := make(map[string]*Node)
	var endNode *Node
	openList := findAll(terrain, 1)
	cost := make(map[string]int)
	cameFrom := make(map[string]*Node)
	for _, n := range openList {
		cost[n.Key()] = 0
		cameFrom[n.Key()] = n
	}

	for len(openList) > 0 {
		sort.Sort(Queue(openList))
		current := openList[0]
		openList = openList[1:]
		if current.pos.Equals(end) {
			endNode = current
			break
		}

		nextNeighbours := current.PossibleNeighbourNodes(terrain)
		for _, neighbour := range nextNeighbours {
			newCost := cost[current.Key()] + 1
			if _, exists := cost[neighbour.Key()]; !exists || newCost < cost[neighbour.Key()] {
				cost[neighbour.Key()] = newCost
				neighbour.priority = newCost + neighbour.pos.HeightDifference(end, terrain)
				cameFrom[neighbour.Key()] = current
				openList = append(openList, neighbour)
			}
		}
		closedList[current.Key()] = current
	}
	steps := 0;
	parent := endNode.parent
	for parent != nil {
		steps++
		parentNode := closedList[parent.Key()]
		parent = parentNode.parent
	}

	return steps 
}

func findAll(terrain [][]int, value int) []*Node {
	var nodes []*Node
	for y, row := range terrain {
		for x, v := range row {
			if (value == v) {
				nodes = append(nodes, &Node{
					pos: &Pos{x,y},
					height: value,
					priority: 0,
				})
			}
		}
	}
	return nodes
}


func getValue(letter rune) int {
    asciiValue := int(letter);
    if (asciiValue >= 97) {
        return asciiValue - 96;
    }
    return asciiValue - 38;
}


func Abs( value int) int {
	if value < 0 {
		return 0 - value
	}
	return value
}

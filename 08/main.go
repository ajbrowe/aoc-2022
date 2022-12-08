package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Tree struct {
	height int
	visible bool
}

var treeMap [][]*Tree

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
		var row []*Tree
		for _, v := range value {
			height, err := strconv.Atoi(string(v))
			if err != nil {
				log.Fatal(err)
			}

			newTree := &Tree{height, true}
			row = append(row, newTree)
		}
		treeMap = append(treeMap, row)
	}
	hiddenTrees := findHidden(treeMap)
	fmt.Printf("Part1 Hidden Trees: %d\n", hiddenTrees)

	maxScenicScore := findMaxScenicScore(treeMap)
	fmt.Printf("Part2 Scenic Score: %d\n", maxScenicScore)
}


func findHidden(treeMap [][]*Tree) int {
	count := 0
	for y, row := range treeMap {
		for x, tree := range row {
			tree.visible = isVisible(treeMap, x, y)
			if tree.visible {
				count++;
			}
		}
	}
	return count
}

func isVisible(treeMap [][]*Tree, x int, y int) bool {
	mapHeight := len(treeMap)
	mapWidth := len(treeMap[0])
	if y == 0 || y == (mapHeight - 1) {
		return true
	}
	if x == 0 || x == (mapWidth - 1) {
		return true
	}

	treeHeight := treeMap[y][x].height
	// check north
	visibleNorth := true
	for scanY := y - 1; scanY >=0; scanY-- {
		if treeMap[scanY][x].height >= treeHeight {
			visibleNorth = false
			break
		}
	}

	// check East
	visibleEast := true
	for scanX := x + 1; scanX < mapWidth; scanX++ {
		if treeMap[y][scanX].height >= treeHeight {
			visibleEast = false
			break
		}
	}

	// check South
	visibleSouth := true
	for scanY := y + 1; scanY < mapHeight; scanY++ {
		if treeMap[scanY][x].height >= treeHeight {
			visibleSouth = false
			break
		}
	}

	// check West
	visibleWest := true
	for scanX := x - 1; scanX >= 0; scanX-- {
		if treeMap[y][scanX].height >= treeHeight {
			visibleWest = false
			break
		}
	}

	return visibleNorth || visibleEast || visibleSouth || visibleWest
}

func findMaxScenicScore(treeMap [][]*Tree) int {
	max := 0
	for y, row := range treeMap {
		for x, _ := range row {
			treeScore := getScenicScore(treeMap, x, y)
			if treeScore > max {
				max = treeScore
			}
		}
	}
	return max
}


func getScenicScore(treeMap [][]*Tree, x int, y int) int {
	mapHeight := len(treeMap)
	mapWidth := len(treeMap[0])

	treeHeight := treeMap[y][x].height
	// check north
	northScore := 0
	for scanY := y - 1; scanY >=0; scanY-- {
		northScore++
		if treeMap[scanY][x].height >= treeHeight {
			break;
		}
	}
	if northScore == 0 && y > 0 {
		northScore = 1
	}

	// check East
	eastScore := 0
	for scanX := x + 1; scanX < mapWidth; scanX++ {
		eastScore++
		if treeMap[y][scanX].height >= treeHeight {
			break
		}

	}
	if eastScore == 0 && x < (mapWidth-1) {
		eastScore = 1
	}

	// check South
	southScore := 0
	for scanY := y + 1; scanY < mapHeight; scanY++ {
		southScore++
		if treeMap[scanY][x].height >= treeHeight {
			break
		}
	}
	if southScore == 0 && y < (mapHeight-1){
		southScore = 1
	}

	// check West
	westScore := 0
	for scanX := x - 1; scanX >= 0; scanX-- {
		westScore++
		if treeMap[y][scanX].height >= treeHeight {
			break
		}
	}
	if westScore == 0 && x > 0 {
		westScore = 1
	}

	return northScore * eastScore * southScore * westScore
}


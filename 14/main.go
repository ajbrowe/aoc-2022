package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Waypoint struct {
	x int
	y int
}

type OffGrid struct {
	count int
	height int
	width int
	floorY int
}

func (og *OffGrid) AddSand() {
	og.count = og.count + 1

	og.width = int(math.Sqrt(float64(og.count * 2) + 0.5))
	if (og.width * (og.width + 1)) / 2 == og.count {
		og.height = og.width
	} else {
		og.height = og.width - 1
	}
	//fmt.Printf("After adding 1 sand {count: %d, height : %d width %d\n", og.count, og.height, og.width)
}

func (og *OffGrid) IsOpen(y int) bool {
	minY := og.floorY - og.height + 1
	//fmt.Printf("offGrid height %d, min Y %d, for %d\n", og.height, minY, y)
	return y < minY
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

	minX, minY, maxX, maxY := 500, 0, 500, 0

	var rockWaypoints [][]*Waypoint
	for scanner.Scan() {
		value := scanner.Text()
		waypointStrings := strings.Split(value, " -> ")

		var waypoints []*Waypoint
		for _, c := range waypointStrings {
			coords := strings.Split(c, ",")
			x, err := strconv.Atoi(coords[0]) 
			if err != nil {
				log.Fatal(err)
			}

			y, err := strconv.Atoi(coords[1])
			if err != nil {
				log.Fatal(err)
			}
			if x > maxX { maxX = x }
			if x < minX { minX = x }
			if y > maxY { maxY = y }
			if y < minY { minY = y }

			nextWaypoint := &Waypoint{
				x: x,
				y: y,
			}
			waypoints = append(waypoints, nextWaypoint)
		}
		rockWaypoints = append(rockWaypoints, waypoints)
	}

	width := (maxX - minX) + 3
	height := (maxY - minY) + 2

	offsetX := (0 - minX) + 1

	//fmt.Printf("width: %d, minX: %d, offsetX: %d\n", width, minX, offsetX)
	//fmt.Printf("height: %d, minY: %d, \n", height, minY,)
	var grid [][]byte
	for y := 0; y < height; y++ {
		var row []byte
		for x := 0; x < width; x++ {
			row = append(row, 0)
		}
		grid = append(grid, row)
	}
	fillGrid(grid, offsetX, rockWaypoints)
	for !dropSand(grid, 500 + offsetX, 100) { }
	drawGrid(grid)
	sandDropped := countSand(grid)
	fmt.Printf("Part1 : Sand dropped = %d\n", sandDropped)

	//grid = addFloor(grid, 0)
	grid = addFloor(grid, 255)

	sandDropped += dropSand2(grid, 500 + offsetX)
	fmt.Printf("Part2 : Sand dropped = %d\n", sandDropped)
	drawGrid(grid)
}


func fillGrid(grid [][]byte, offsetX int,  waypoints [][]*Waypoint) {
	for _, wpline := range waypoints {
		lenLine := len(wpline)
		for i, wp := range wpline[:lenLine-1] {
			grid[wp.y][wp.x + offsetX] = 255

			wpnext := wpline[i+1]
			if (wpnext.x == wp.x) {
				x := wp.x + offsetX
				startY, endY := wp.y, wpnext.y
				if startY > endY {
					startY, endY = endY, startY
				}
				//fmt.Printf("Drawing line from %d,%d to %d,%d\n", x, startY, x, endY)
				for y := startY; y <= endY; y++ {
					grid[y][x] = 255
				}
			} else if wpnext.y == wp.y {
				y := wp.y
				startX, endX := wp.x + offsetX, wpnext.x + offsetX
				if startX > endX {
					startX, endX = endX, startX
				}
				//fmt.Printf("Drawing line from %d,%d to %d,%d\n", startX, y, endX, y)
				for x := startX; x <= endX; x++ {
					grid[y][x] = 255
				}
			}
		}
	}
}


func dropSand(grid [][]byte, startX int, count int) bool {
	for i := 0; i < count; i++ {
		sandX := startX
		for sandY := 1; sandY < len(grid); sandY++ {
			if sandY == len(grid) - 1 {
				return true
			}
			if grid[sandY+1][sandX] == 0 {
				continue
			}
			if grid[sandY+1][sandX] != 0 {
				if grid[sandY+1][sandX-1] == 0 {
					sandX -= 1
				} else if grid[sandY+1][sandX+1] == 0 {
					sandX += 1
				} else {
					grid[sandY][sandX] = 1
					break
				}
			}
		}
	}
	return false
}

func addFloor(grid [][]byte, value byte) [][]byte {
	width := len(grid[0])
	var floor []byte
	for i := 0; i < width; i++ {
		floor = append(floor, value)
	}
	grid = append(grid, floor)
	return grid
}


func dropSand2(grid [][]byte, startX int) int {
	minX := 0
	maxX := len(grid[0]) - 1;
	maxY := len(grid) - 2;
	sandFixed := 0
	offGridLeft := &OffGrid{count: 0, width: 0, height: 0, floorY: maxY }
	offGridRight := &OffGrid{count: 0, width: 0, height: 0, floorY: maxY}
	for {
		sandX := startX
		if grid[0][startX] != 0 {
			break
		}
		for sandY := 0; sandY < len(grid) - 1; sandY++ {
			if sandY  == maxY && grid[sandY][sandX] == 0 {
				grid[sandY][sandX] = 1
				sandFixed++
				break
			}
			if grid[sandY+1][sandX] == 0 {
				continue
			}

			if (sandX - 1) < minX {
				//fmt.Printf("%d - 1 < %d\n", sandX, minX)
				if offGridLeft.IsOpen(sandY+1) {
					offGridLeft.AddSand()
					break
				} else {
					if grid[sandY+1][sandX+1] == 0 { 
						sandX += 1
					} else {
						grid[sandY][sandX] = 1
						sandFixed++
						break
					}
				}
			} else if grid[sandY+1][sandX-1] == 0 {
				sandX -= 1
				//fmt.Printf("Moving left 1 column to %d\n", sandX)
			} else if (sandX + 1) > maxX {
				//fmt.Printf("%d + 1 > %d\n", sandX, maxX)
				if offGridRight.IsOpen(sandY+1) {
					offGridRight.AddSand()
					break
				} else {
					if grid[sandY+1][sandX-1] == 0 {
						sandX -=1
					} else {
						grid[sandY][sandX] = 1
						sandFixed++
						break
					}
				}
			} else if grid[sandY+1][sandX+1] == 0 {
				sandX += 1
				//fmt.Printf("Moving right 1 column to %d\n", sandX)
			} else{
				grid[sandY][sandX] = 1
				sandFixed++
				break
			}
		}
	}

	fmt.Printf("Sand in grid: %d, sand off grid left: %d, sand off grid right: %d\n",sandFixed, offGridLeft.count, offGridRight.count)
	return sandFixed + offGridLeft.count + offGridRight.count
}



func drawGrid(grid [][]byte) {
	for _, row := range grid {
		var sb strings.Builder
		for _, v := range row {
			if v == 255 {
				sb.WriteString("#")
			}
			if v == 0 {
				sb.WriteString(".")
			}
			if v == 127 {
				sb.WriteString("+")
			}
			if v == 1 {
				sb.WriteString("O")
			}
		}
		fmt.Println(sb.String())
	}
}


func countSand(grid [][]byte) int {
	count := 0;
	for _, row := range grid {
		for _, v := range row {
			if v == 1 {
				count += 1
			}
		}
	}
	return count
}


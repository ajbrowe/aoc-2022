package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

type RockShape struct {
	height int
	width int
	layout [][]bool
}


type Rock struct {
	x int
	y int
	shape RockShape
}

func (r *Rock) Top() int {
	return r.y + r.shape.height - 1
}

func (r *Rock) Left() int {
	return r.x
}

func (r *Rock) Right() int {
	return r.x + r.shape.width - 1 
}

func (r *Rock) Bottom() int {
	return r.y
}


func (r *Rock) Collision(chamber [][]byte, dx,dy int) bool {


	if r.Bottom() + dy <= 0 || r.Left() + dx < 0 || r.Right() + dx > 6 {
		return true

	}

	lt := r.shape.layout
	for i, ltRow := range lt {
		for j, solid := range ltRow {
			scanX := r.x + j + dx
			scanY := r.y + i + dy
			if scanY < 0 || scanY >= len(chamber) || scanX < 0 || scanX > 6 {
				continue
			}
			if solid && chamber[scanY][scanX] != 0 {
				return true
			}
		}
	}
	return false
}

func (r *Rock) CanPushRight(chamber [][]byte) bool {
	return !r.Collision(chamber, 1, 0)
}

func (r *Rock) PushRight(chamber [][]byte) {
	if r.CanPushRight(chamber) {
		r.x++
	}
}
func (r *Rock) CanPushLeft(chamber [][]byte) bool {
	return !r.Collision(chamber, -1, 0)
}

func (r *Rock) PushLeft(chamber [][]byte) {
	if r.CanPushLeft(chamber) {
		r.x--
	}
}

func (r *Rock) CanFall(chamber [][]byte) bool {
	return !r.Collision(chamber, 0, -1)
}

func (r *Rock) Fall() {
	r.y--
}

func (r *Rock) DisplayPos() {
	fmt.Printf("Rock at x: %d, y: %d\n", r.x, r.y)
}

func (r *Rock) Rest(chamber [][]byte) {
	//fmt.Printf("Rock resting at %d,%d\n", r.x, r.y)	
	for i := (r.shape.height - 1); i >= 0; i-- {
		for j := 0; j < r.shape.width; j++ {
			if r.shape.layout[i][j] {
				chamber[r.y+i][r.x+j] = 255
			}
		}
	}
}

func initRockTypes() []RockShape {

	rockTypes := []RockShape{}

	rockTypes = append(rockTypes, RockShape{1, 4, [][]bool{0:{true,true,true,true}}})
	rockTypes = append(rockTypes, RockShape{3, 3, [][]bool{0:{false,true,false},1:{true,true,true},2:{false,true,false}}})
	rockTypes = append(rockTypes, RockShape{3, 3, [][]bool{0:{true,true,true},1:{false,false,true},2:{false,false,true}}})
	rockTypes = append(rockTypes, RockShape{4, 1, [][]bool{0:{true},1:{true},2:{true},3:{true}}})
	rockTypes = append(rockTypes, RockShape{2, 2, [][]bool{0:{true,true},1:{true,true}}})

	return rockTypes
}

type SignatureData struct{
	count int
	height int
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

	var jetPattern string
	for scanner.Scan() {jetPattern = scanner.Text() }
	jetPatternCount := len(jetPattern)

	var chamber [][]byte
	floor := []byte{127,127,127,127,127,127,127}
	chamber = append(chamber, floor)

	rockTypes := initRockTypes()

	rockCounter := 0
	jetCounter := 0

	highestRock := 0

	
	patternRepeats := make(map[string][]SignatureData)

	var loopIncrement SignatureData
	for {
		if (len(chamber) -1) < highestRock + 4 {
			toAdd := (highestRock + 4) - (len(chamber) - 1)
			chamber = GrowChamber(chamber, toAdd)
		}
		chamberHeight := len(chamber) - 1
		if rockCounter < 12 || rockCounter % 200 == 0 {fmt.Println(rockCounter);DrawChamber(chamber, 20) }
		rock := NewRock(rockTypes, rockCounter, chamberHeight)
		signature := fmt.Sprintf("%d|%d|%s", rockCounter % 5, jetCounter % jetPatternCount, ChamberProfile(chamber))
		if _, exists := patternRepeats[signature]; !exists {
			patternRepeats[signature] = []SignatureData{SignatureData{rockCounter,highestRock}}
		} else {
			patternRepeats[signature] = append(patternRepeats[signature], SignatureData{rockCounter,highestRock})
			if len(patternRepeats[signature]) == 3 {
				// Found enough data about repeats to make some calculations
				loopIncrement.count = patternRepeats[signature][1].count - patternRepeats[signature][0].count
				loopIncrement.height = patternRepeats[signature][1].height - patternRepeats[signature][0].height
				break
			}
		}

		if rockCounter == 2020 {
			fmt.Printf("Part1: HighestRock %d\n", highestRock)
		}
		rockCounter++
		for {
			jetMove := jetPattern[ jetCounter % jetPatternCount ]
			jetCounter++
			if jetMove == '>' {
				rock.PushRight(chamber)
			}
			if jetMove == '<' {
				rock.PushLeft(chamber)
			}
			if !rock.CanFall(chamber) {
				rock.Rest(chamber)
				if rock.Top() > highestRock {
					highestRock = rock.Top()
				}
				break
			}
			rock.Fall()
		}
	}


	fmt.Println("Part2....")
	fmt.Printf("Number of patterns %d\n", len(patternRepeats))
	repeatingPatterns := 0
	var startsRepeating SignatureData
	differences := make(map[int]int)
	for _, v := range patternRepeats {
		if len(v) > 1 {
			repeatingPatterns++;
			if v[0].count < startsRepeating.count {
				startsRepeating.count = v[0].count
				startsRepeating.height = v[0].height
			}
			for i, _ := range v[:len(v)-1] {
				diff := v[i+1].height - v[i].height
				differences[diff]++
			}
		}
	}
	fmt.Printf("Number of Repeating patterns %d\n", repeatingPatterns)
	fmt.Printf("Starts Repeating from %d\n", startsRepeating.count)
	fmt.Println(differences)

	heightIncrements := make([]int, repeatingPatterns * 2)
	for _, v := range patternRepeats {
		if len(v) > 0 {
			idx := v[0].count - startsRepeating.count
			hi := v[0].height - startsRepeating.height
			heightIncrements[idx] = hi
		}
	}

	iterations := 1000000000000
	repeats := (iterations - startsRepeating.count) / loopIncrement.count
	remainder := (iterations - startsRepeating.count) % loopIncrement.count

	trilliHeight := repeats * loopIncrement.height + startsRepeating.height + heightIncrements[remainder]

	fmt.Printf("Part 2 Height %d\n", trilliHeight)


}

func NewRock(rockTypes []RockShape, counter int, chamberHeight int ) *Rock {
	newRockShape := rockTypes[ counter % 5 ]
	return &Rock{
		x: 2,
		y: chamberHeight, // Rock starts at highest point in chamber
		shape: newRockShape,
	}
}

func GrowChamber(chamber [][]byte, rows int) [][]byte {
	for r := 0; r < rows; r++ {
		space := []byte{0, 0, 0, 0, 0, 0, 0}
		chamber = append(chamber, space)
	}
	return chamber
}


func ChamberProfile(chamber [][]byte) string {
	var profile []string
	for x := 0; x < 7; x++ {
		profileDepth := 0
		for y := len(chamber) - 1; y >= 0; y-- {
			if chamber[y][x] == 0 {
				profileDepth++
			} else {
				break
			}
		}
		profile = append(profile, fmt.Sprintf("%d",profileDepth))
	}

	return strings.Join(profile, ":")
}


func DrawChamber(chamber [][]byte, head int) {
	height := len(chamber)
	for y:= height-1; y >= 0 && y >= (height - head); y-- {
		var sb strings.Builder
		sb.WriteString("|")
		for _, b := range chamber[y] {
			if b == 0 {
				sb.WriteString(" ")
			} else if b == 127 {
				sb.WriteString("-")
			} else {
				sb.WriteString("#")
			}
		}
		sb.WriteString("|")
		fmt.Println(sb.String())
	}
	fmt.Println()
}

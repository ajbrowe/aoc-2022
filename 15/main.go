package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Pos struct{
	x int
	y int
}

type Range struct{
	start int
	end int
}

type Sensor struct{
	pos *Pos
	beacon *Pos
}

func (s *Sensor) BeaconDistance() int {
	return Abs(s.pos.x - s.beacon.x) + Abs(s.pos.y - s.beacon.y)
}


func (s *Sensor) CoverageAtY(y int) *Range {
	bDistance := s.BeaconDistance()
	yDistance := Abs(s.pos.y - y)
	if yDistance > bDistance {
		return nil
	}
	start := s.pos.x - (bDistance - yDistance)
	end := s.pos.x + (bDistance - yDistance)
	return &Range{start, end,}
}
func (s *Sensor) BeaconAtY(y int) bool {
	return s.beacon.y == y
}


func (r *Range) Overlaps(o *Range) bool {
	if r.start > o.end || r.end < o.start {
		return false
	}
	return true
}

func (r *Range) Contains(o *Range) bool {
	return r.start <= o.start && r.end >= o.end
}

func (r *Range) Outside(o *Range) []*Range {
	var outsideRanges []*Range
	if r.start < o.start {
		outsideRanges = append(outsideRanges, &Range{r.start, o.start-1})
	}
	if r.end > o.end {
		outsideRanges = append(outsideRanges, &Range{o.end+1,r.end})
	}
	return outsideRanges
}

func (r *Range) Merge(s *Range) *Range {
	minStart := r.start
	maxEnd := r.end
	if s.start < minStart {
		minStart = s.start
	}
	if s.end > maxEnd {
		maxEnd = s.end
	}
	return &Range{minStart, maxEnd}
}

func (r *Range) Length() int {
	return (r.end - r.start) + 1
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

	var sensors []*Sensor

	for scanner.Scan() {
		value := scanner.Text()
		sensorRe := regexp.MustCompile(`Sensor at x=(-?\d+), y=(-?\d+): closest beacon is at x=(-?\d+), y=(-?\d+)`)
		matches := sensorRe.FindAllStringSubmatch(value, -1)

		sx, err := strconv.Atoi(matches[0][1])
		if err != nil {	log.Fatal(err) 	}
		sy, err := strconv.Atoi(matches[0][2])
		if err != nil {	log.Fatal(err) 	}
		bx, err := strconv.Atoi(matches[0][3])
		if err != nil {	log.Fatal(err) 	}
		by, err := strconv.Atoi(matches[0][4])
		if err != nil {	log.Fatal(err) 	}

		sensors = append(sensors, &Sensor{
			pos: &Pos{sx, sy},
			beacon: &Pos{bx, by},
		})
	}
	mergedRanges, beaconsAtY := findRangesAt(10, sensors)
	if len(sensors) > 14 {
		mergedRanges, beaconsAtY = findRangesAt(2000000, sensors)
	}


	part1Sum := 0 - len(beaconsAtY)
	for _, m := range mergedRanges {
		part1Sum += m.Length()
	}

	fmt.Printf("Part1 Sum: %d\n", part1Sum)

	part2MaxY := 20
	if len(sensors) > 14 {
		part2MaxY = 4000000
	}
	captureRange := &Range{0, part2MaxY}
	var foundBeacon Pos
	for y := 0; y <= part2MaxY; y++ {
		mRanges, _ := findRangesAt(y, sensors)
		var outsideRanges []*Range
		captured := false
		for _, m := range mRanges {
			if m.Contains(captureRange) {
				captured = true
				break
			}
			outside := captureRange.Outside(m)
			if len(outside) > 0 {
				outsideRanges = append(outsideRanges, outside...)
			}
		}
		if captured {
			continue
		}

		for len(outsideRanges) > 0 {
			// shift
			or := outsideRanges[0]
			outsideRanges = outsideRanges[1:]
			
			overlaps := false
			for _, m := range mRanges {
				if (m.Overlaps(or)) {
					overlaps = true
					outsideRanges = append(outsideRanges, or.Outside(m)...)
				}
			}

			if !overlaps && len(outsideRanges) == 1 {
				foundSpace := outsideRanges[0]
				if foundSpace.Length() == 1 {
					foundBeacon = Pos{foundSpace.start, y}
					break
				}
			}
		}
	}

	part2Sum := (foundBeacon.x * 4000000) + foundBeacon.y
	fmt.Printf("Part2 Sum: %d\n", part2Sum)
}


func findRangesAt(y int, sensors []*Sensor) ([]*Range, map[string]bool){
	var matchingRanges []*Range
	beaconsAtY := make(map[string]bool)
	for _, s := range sensors {
		r := s.CoverageAtY(y) // sample
		//r := s.CoverageAtY(2000000) // input
		if r != nil {
			matchingRanges = append(matchingRanges, r)
		}
		if s.BeaconAtY(y) { // sample
		//if s.BeaconAtY(2000000) { // input
			beaconsAtY[fmt.Sprintf("%d,%d", s.beacon.x, s.beacon.y)] = true
		}
	}

	var mergedRanges []*Range
	merged := make(map[int]bool) 
	for i := 0; i < len(matchingRanges) - 1; i++ {
		testRange := matchingRanges[i]
		//fmt.Printf("Inspecting range %d, start: %d end: %d\n", i +1, testRange.start, testRange.end)
		for j, r := range matchingRanges[i+1:] {
			if _, exists := merged[i + j + 1]; exists {
				continue
			}
			if testRange.Overlaps(r) {
				//fmt.Printf("  Range %d, overlaps with range %d, start: %d end: %d\n", i +1, (i + j + 2), r.start, r.end)
				testRange = testRange.Merge(r)
				merged[i] = true
				merged[i + j + 1] = true
			}
		}
		independentRange := true
		for i, r := range mergedRanges {
			if r.Contains(testRange) {
				independentRange = false
			}
			if r.Overlaps(testRange) {
				mergedRanges[i] = mergedRanges[i].Merge(testRange)
				independentRange = false
			}
		}
		if independentRange {
			mergedRanges = append(mergedRanges, testRange)
		}

	}
	return mergedRanges, beaconsAtY
}
	

func Abs(value int) int {
	if value < 0 {
		value = 0 - value
	}
	return value
}



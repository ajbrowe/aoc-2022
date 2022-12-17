package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type graph struct {
	to string
	weight float64
}

type ValveCave struct {
	rate int
	links []string
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

	caves := make(map[string]*ValveCave)

	for scanner.Scan() {
		value := scanner.Text()
		valveRe := regexp.MustCompile(`Valve ([A-Z]{2}) has flow rate=(\d+); tunnels? leads? to valves? ([A-Z, ]+)`)
		matches := valveRe.FindAllStringSubmatch(value, -1)

		valveName := matches[0][1]
		valveRate, err := strconv.Atoi(matches[0][2])
		if err != nil {
			log.Fatal(err)
		}

		links := strings.Split(matches[0][3], ", ")
		caves[valveName] = &ValveCave{
			rate: valveRate,
			links: links,
		}
	}

	distances := buildDistances(buildGraph(caves))

	removeZeroRateCaves(distances, caves, "AA")


	fmt.Println(distances)

	part1Sum := search("AA", caves, distances, 0, 0, 0, nonZeroCaveNames(caves), 30)

	fmt.Printf("Part1 Pressure Released %d\n", part1Sum)
	me := Player{0, "AA"}
	elephant := Player{0, "AA"}
	part2Sum := search2(me, elephant, caves, distances, 0, 0, nonZeroCaveNames(caves),26)
	
	fmt.Printf("Part2 Pressure Released %d\n", part2Sum)


}

func buildGraph(caves map[string]*ValveCave) map[string][]graph {
	g := make(map[string][]graph)
	for k, v := range caves {
		var graphLinks []graph
		for _, l := range v.links {
			graphLinks = append(graphLinks, graph{l, 1})
		}
		g[k] = graphLinks
	}
	return g
}


func buildDistances( g map[string][]graph) map[string]map[string]float64 {
	// Floyd-Warshall algorithm - yeah I had to duckduckgo it too.

	var caveNames []string
	for k := range g {
		caveNames = append(caveNames, k)
	}

	dist := make(map[string]map[string]float64, len(g))
	for _, iname:= range caveNames {
		di := make(map[string]float64, len(g))

		for _, jname := range caveNames {
			di[jname] = math.Inf(1)
		}
		di[iname] = 0
		dist[iname] = di
	}

	for _, name := range  caveNames {
		graphs := g[name]
		for _, v := range graphs {
			dist[name][v.to] = v.weight
		}
	}

	// for all the cave names to all the cave names find the total distance
	// to each
	for _, primeCaveName := range caveNames {
		for _, secondCaveName := range caveNames {
			for jname, distj := range dist[secondCaveName] {
				if d := dist[secondCaveName][primeCaveName] + dist[primeCaveName][jname]; distj > d {
					dist[secondCaveName][jname] = d
				}
			}
		}
	}
	return dist
}


func removeZeroRateCaves(distances map[string]map[string]float64, caves map[string]*ValveCave, preserve string) {
	var zeroCaves []string
	for caveName, cave := range caves {
		if cave.rate == 0 	{
			zeroCaves = append(zeroCaves, caveName)
			if caveName != preserve {
				delete(distances, caveName)
			}
		}
	}

	for _, caveDist := range distances {
		for _, caveName := range zeroCaves {
			delete(caveDist, caveName)
		}
	}
}

func nonZeroCaveNames(caves map[string]*ValveCave) []string {
	var CaveNames []string
	for caveName, cave := range caves {
		if cave.rate != 0 	{
			CaveNames = append(CaveNames, caveName)
		}
	}
	return CaveNames
}


func search(caveName string, caves map[string]*ValveCave, distances map[string]map[string]float64, minute int, currentPressure int, currentFlow int, remaining []string, limit int) int {
	nowScore := currentPressure + (limit - minute) * currentFlow
	max := nowScore

	for _, cName := range remaining {
		timeTaken := int(distances[caveName][cName]) + 1
		if minute + timeTaken < limit {
			newMinute:= minute + timeTaken
			newPressure := currentPressure + timeTaken * currentFlow
			newFlow := currentFlow + caves[cName].rate
			newRemaining := removeFromList(remaining, cName)
			possibleScore := search(cName, caves, distances, newMinute, newPressure, newFlow, newRemaining, limit)
			if (possibleScore > max) {
				max = possibleScore
			}
		}
	}

	return max
}

func removeFromList( in []string, v string) []string {
	out := []string{}

	for _, s := range in {
		if s != v {
			out = append(out, s)
		}
	}
	return out
}

type Player struct {
	minute int
	caveName string
}

func search2(player Player, other Player, caves map[string]*ValveCave, distances map[string]map[string]float64, currentPressure int, currentFlow int, remaining []string, limit int) int {
	nowScore := currentPressure + (limit - player.minute) * currentFlow

	max := nowScore

	for _, cName := range remaining {
		timeTaken := int(distances[player.caveName][cName]) + 1
		if player.minute + timeTaken < limit {
			newPlayer := Player{
				minute: player.minute + timeTaken,
				caveName: cName,
			}
			newPressure := currentPressure + timeTaken * currentFlow
			newFlow := currentFlow + caves[cName].rate
			newRemaining := removeFromList(remaining, cName)
			possibleScoreP := search2(newPlayer, other, caves, distances, newPressure, newFlow, newRemaining, limit,)
			possibleScoreE := search2(other, newPlayer, caves, distances, newPressure, newFlow, newRemaining, limit,)
			if (possibleScoreP > max) {
				max = possibleScoreP
			}
			if (possibleScoreE > max) {
				max = possibleScoreE
			}
}
	}
	return max
}

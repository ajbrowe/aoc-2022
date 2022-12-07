package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	class    string
	name     string
	size     int
	children []*Node
}

func (n *Node) Size() int {
	size := n.size
	for _, c := range n.children {
		size += c.Size()
	}
	return size
}

func (n *Node) AddChild(c *Node) {
	n.children = append(n.children, c)
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
	directories := map[string]*Node{
		"/": &Node{
			class: "directory",
			name:  "/",
		},
	}
	currentDir := "/"
	listMode := false

	pathStack := []string{}

	for scanner.Scan() {
		value := scanner.Text()
		fields := strings.Fields(value)
		if fields[0] == "$" {
			listMode = false
			if fields[1] == "cd" {
				name := fields[2]
				if name == ".." {
					// pop pathStack
					currentDir, pathStack = pathStack[len(pathStack)-1], pathStack[:len(pathStack)-1]
				} else {
					// push pathStach
					pathStack = append(pathStack, currentDir)
					currentDir = changeDirectory(directories, Path(pathStack), name)
				}
			}
			if fields[1] == "ls" {
				listMode = true
			}
		} else if listMode {
			if fields[0] == "dir" {
				name := fields[1]
				newDir := &Node{
					class: "directory",
					name: name,
					size: 0,
				}
				path := Path(append(pathStack, currentDir, name));
				if _,exists := directories[path]; !exists {
					directories[path] = newDir
					directories[Path(append(pathStack, currentDir))].AddChild(newDir)
				}
			} else {
				fileSize, err := strconv.Atoi(fields[0])
				if (err != nil) {
					log.Fatal(err)
				}
				newFile := &Node{
					class: "file",
					name: fields[1],
					size: fileSize,
				}
				directories[Path(append(pathStack, currentDir))].AddChild(newFile)
			}
		}
	}
	part1Sum := 0
	unusedSpace := 70000000 - directories["/"].Size()
	requiredSpace := 30000000 - unusedSpace
	var bigEnough []int
	for _, d := range directories {
		dSize := d.Size()
		if dSize <= 100000 {
			part1Sum += dSize
		}
		if dSize >= requiredSpace {
			bigEnough = append(bigEnough, dSize)
		}
	}

	fmt.Printf("Part1 sum: %d\n", part1Sum)
	fmt.Printf("Part2 unused space: %d\n", unusedSpace)
	fmt.Printf("Part2 required space: %d\n", requiredSpace)
	fmt.Printf("Part2 smallest Dir to delete: %d\n", Min(bigEnough))

}

func changeDirectory(directories map[string]*Node, path string, name string) string {
	if name == "/" {
		return name
	}
	fullName := fmt.Sprintf("%s/%s", path, name)
	//return findChildDir(directories, currentDir, name)
	if _, exists := directories[fullName]; exists {
		return name
	}
	log.Fatalf("Directory %s does not exist\n", name)
	return "/"
}

func findChildDir(directories map[string]*Node, currentDir, name string) string {
	dir := directories[currentDir]
	for _, n := range dir.children {
		if n.class == "directory" {
			if n.name == name {
				return name
			}
		}
	}
	log.Printf("Directory %s not found\n", name)
	return "/"
}

func Path(pathStack []string) string {
	return strings.Join(pathStack[1:], "/")
}

func Min(data []int) int {
	min := data[0];
	for _, v := range data[1:] {
		if v < min {
			min = v
		}
	}
	return min
}

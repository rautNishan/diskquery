package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

type indexValue struct {
	offset int
	length int
}

func main() {
	index := make(map[int]indexValue)
	PopulateIndex(&index, "./data.txt")
	data := findByIndex(5000, &index, "./data.txt")
	fmt.Println("This is data: ", string(data))
}

func PopulateIndex(index *map[int]indexValue, path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Erro while opening file: ", err)
	}
	fmt.Println("This is file: ", file)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	startOffset := 0

	for scanner.Scan() {
		line := scanner.Text()
		length := len(line)

		idx := 0
		toAdd := 0
		for i, b := range line {
			if b == ',' {
				idx = i
				toAdd = i + 1
				break
			}
		}
		tempId, err := strconv.Atoi(line[:idx])

		if err != nil {
			log.Fatal("Error while extracting index value: ", err)
		}

		(*index)[tempId] = indexValue{
			offset: startOffset + toAdd,
			length: length - toAdd,
		}
		startOffset += length + 1
	}
}

func findByIndex(i int, index *map[int]indexValue, path string) []byte {
	file, err := os.Open(path)
	if err != nil {
		log.Fatal("Error while Opening file in find: ", err)
	}
	defer file.Close()

	buffer := make([]byte, (*index)[i].length)
	n, err := file.ReadAt(buffer, int64((*index)[i].offset))

	if err != nil {
		log.Fatal("Error while reading file in find: ", err)
	}
	fmt.Println("This is n: ", n)
	return buffer
}

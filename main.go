package main

import (
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
	fmt.Println("This is index: ", index)

}

func PopulateIndex(index *map[int]indexValue, path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Erro while reading file: ", err)
	}
	getIndex := string(data[0])
	tempI, err := strconv.Atoi(getIndex)
	if err != nil {
		log.Fatal("Error: ", err)
	}
	startOffset := 0

	for i, b := range data {
		if b == 10 {
			(*index)[tempI] = indexValue{
				offset: startOffset,
				length: i - startOffset,
			}
			tempI++
			startOffset = i + 1
		}
	}
}

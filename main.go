package main

import "fmt"

type indexValue struct {
	offset int
	length int
}

func main() {
	fmt.Println("This is disk read")
	var index map[int]indexValue
	fmt.Println("This is index: ", index)
	PopulateIndex(&index)
}

func PopulateIndex(index *map[int]indexValue) {

}

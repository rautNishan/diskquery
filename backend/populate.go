package main

// import (
// 	"fmt"
// 	"log"
// 	"os"
// )

// func main() {
// 	const n = 5000000
// 	path := "data.txt"

// 	file, err := os.Create(path)
// 	if err != nil {
// 		log.Fatal("Error creating file:", err)
// 	}
// 	defer file.Close()

// 	for i := 1; i <= n; i++ {
// 		line := fmt.Sprintf("%d,This is data %d\n", i, i)
// 		_, err := file.WriteString(line)
// 		if err != nil {
// 			log.Fatalf("Error writing line %d: %v", i, err)
// 		}
// 	}

// 	// Get file size
// 	info, err := os.Stat(path)
// 	if err != nil {
// 		log.Fatal("Error getting file info:", err)
// 	}

// 	fmt.Printf("File %s populated with %d lines.\n", path, n)
// 	fmt.Printf("File size: %d bytes (%.2f MB)\n", info.Size(), float64(info.Size())/(1024*1024))
// }

package main

import (
	"flag"
	"fmt"
	"bufio"
	"log"
	"os"
	"sync"

)
var wg = sync.WaitGroup{}


func publish(message string) {
	fmt.Println(message)
	wg.Done()
}

func main() {

	var fileName string

	flag.StringVar(&fileName, "f", "somefile.txt", "Filename")

	flag.Parse()

	file, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		wg.Add(1)
		publish(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	wg.Wait()
}
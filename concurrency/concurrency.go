package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	files := []string{"file1.txt", "file2.txt", "file3.txt"}
	flag.Parse()
	switch flag.Arg(0) {
	case "conc":
		concurrently(files)
	case "seq":
		sequentially(files)
	default:
		log.Fatal("need conc or seq arg")
	}
}

func concurrently(files []string) {
	errc := make(chan error)
	for _, file := range files {
		go func(f string) {
			errc <- processFile(f)
		}(file)
	}
	for i := 0; i < len(files); i++ {
		err := <-errc
		if err != nil {
			fmt.Println(err)
		}
	}
}

func sequentially(files []string) {
	for _, file := range files {
		if err := processFile(file); err != nil {
			fmt.Println(err)
		}
	}
}

func processFile(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	fmt.Printf("processing %s\n", filename)
	time.Sleep(50 * time.Millisecond)
	defer f.Close()
	return nil
}

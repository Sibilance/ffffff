package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

var help = flag.Bool("help", false, "print usage")
var fileName = flag.String("file", "", "YAML file to parse")

func main() {
	flag.Parse()
	if *fileName == "" {
		*help = true
		fmt.Println("Missing required flag 'file'.")
	}
	if *help {
		fmt.Printf("Usage: %s [options]\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalf("Error opening file: %s\n", err)
	}

	decoder := yaml.NewDecoder(file)
	for {
		decoded := make(map[interface{}]interface{})
		if err := decoder.Decode(&decoded); err != nil {
			if err != io.EOF {
				log.Fatalln(err)
			}
			break
		}
		fmt.Println(decoded)
	}
}

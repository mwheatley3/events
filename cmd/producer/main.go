package main

import (
	"curbflow/producer"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	fmt.Println("Running producer")
	if len(os.Args) < 2 {
		log.Fatal("must include a millisecond argument to start producer")
	}
	ms, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("error converting command line agr: [%s] to int. err: %s", os.Args[1], err)
	}

	producer.Start(time.Duration(ms))
}

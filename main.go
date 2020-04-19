package main

import (
	// "errors"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	f, err := logSetup()
	if err != nil {
		fmt.Printf("%v", err)
		return 1
	}
	defer f.Close()

	// var opt option
	opt := option{
		date: time.Now(),
	}

	cli := cli{
		option: opt,
	}

	if err := cli.run(); err != nil {
		return 1
	}

	return 0
}

func logSetup() (*os.File, error) {
	f, err := os.OpenFile("system.log", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}
	log.SetOutput(f)
	return f, nil
}

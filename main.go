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

const logFileName = "system.log"

func run(args []string) int {
	f, err := logSetup()
	if err != nil {
		fmt.Printf("%v\n", err)
		return 1
	}
	defer f.Close()

	opt := option{
		date: time.Now(),
	}

	cli := cli{
		option: opt,
	}

	log.Printf("[INFO] start pid: %d", os.Getpid())
	if err := cli.run(); err != nil {
		fmt.Printf("%v\n", err)
		log.Fatalf(fmt.Sprintf("[ERROR] %v\n", err))
	}
	fmt.Println("入力を保存しました☕")
	log.Printf("[INFO] end pid: %d", os.Getpid())

	return 0
}

func logSetup() (*os.File, error) {
	f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}
	log.SetOutput(f)
	return f, nil
}

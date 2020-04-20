package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	os.Exit(run(os.Args[1:]))
}

const (
	logFileName = "system.log"
	dsrDir = ".dsr"
)

var (
	dsrPath = filepath.Join(os.Getenv("HOME"), dsrDir)
	logFilePath = filepath.Join(dsrPath, logFileName)
)

func run(args []string) int {
	if _, err := os.Stat(dsrPath); os.IsNotExist(err) {
    os.Mkdir(dsrPath, 0777)
	}

	f, err := logSetup()
	if err != nil {
		fmt.Printf("%v\n", err)
		return 1
	}
	defer f.Close()


	var date time.Time
	var dateArg string
	flag.StringVar(&dateArg, "d", "", "specify date to record. e.g. 2006-1-2")
	flag.Parse()
	if len(dateArg) == 0 {
		date = time.Now()
	} else {
		date, err = time.Parse("2006-1-2", dateArg)
		if err != nil {
			fmt.Printf("%v\n", err)
			log.Fatalf(fmt.Sprintf("[ERROR] %v\n", err))
		}
	}

	cli := cli{
		option: option{
			date: date,
		}, 
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
	f, err := os.OpenFile(logFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)

	if err != nil {
		return nil, err
	}
	log.SetOutput(f)
	return f, nil
}

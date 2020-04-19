package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/manifoldco/promptui"
)

type option struct {
	date time.Time
}

type cli struct {
	option option
}

func (c cli) run() error {
	file, err := os.OpenFile("record.csv", os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return errors.New("ファイルの読み込みに失敗しました")
	}

	w := csv.NewWriter(file)

	if err := w.Write([]string{"hoge1", "hoge2"}); err != nil {
		return err
	}

	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}

	prompt := promptui.Select{
		Label: "Select Day",
		Items: []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday",
			"Saturday", "Sunday"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return errors.New("エラー")
	}

	fmt.Println(result)

	return nil
}

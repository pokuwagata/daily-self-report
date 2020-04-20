package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/manifoldco/promptui"
)

type option struct {
	date time.Time
}

type cli struct {
	option option
}

type question struct {
	order    int
	question string
	validate promptui.ValidateFunc
}

const csvFilename = "record.csv"

var (
	numValidate = func(input string) (float64, error) {
		n, err := strconv.ParseFloat(input, 64)
		if err != nil {
			return 0, errors.New("Invalid number")
		}
		return n, nil
	}

	questions = []question{
		{
			order:    1,
			question: "睡眠時間は？",
			validate: func(input string) error {
				num, err := numValidate(input)
				if err != nil {
					return err
				}
				if num < 0 {
					return errors.New("0以上の値を入力してください")
				}
				return nil
			},
		},
		{
			order:    2,
			question: "食事の回数は？",
			validate: func(input string) error {
				num, err := numValidate(input)
				if err != nil {
					return err
				}
				if num < 0 {
					return errors.New("0以上の値を入力してください")
				}
				return nil
			},
		},
	}
)

func buildHeader() []string {
	header := []string{"date"}
	for i := range questions {
		header = append(header, strconv.Itoa(i+1))
	}

	return header
}

func parseDate(d time.Time) string {
	return fmt.Sprintf("%d-%d-%d", d.Year(), d.Month(), d.Day())
}

func (c cli) run() error {
	log.Println("[INFO] start main process")

	date := parseDate(time.Now())

	log.Println("[INFO] open csv file")
	file, err := os.OpenFile(csvFilename, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func() {
		log.Println("[INFO] close csv file")
		file.Close()
	}()

	// CSV ファイルの前処理
	r := csv.NewReader(file)
	w := csv.NewWriter(file)

	records, err := r.ReadAll();
	
	if err != nil {
		if err == io.EOF {
			// ファイルが空の場合はヘッダーを作成する
			if err := w.Write(buildHeader()); err != nil {
				return err
			}
		} 
		return err
	}

	for _, row := range records {
		if row[0] == date{
			return fmt.Errorf("%s は記録済みです", date)
		}
	}

	answers := []string{date}

	// 質問を行う
	for i, q := range questions {
		prompt := promptui.Prompt{
			Label:    q.question,
			Validate: q.validate,
		}

		log.Printf("[INFO] start q %d\n", i+1)
		r, err := prompt.Run()

		if err != nil {
			return err
		}

		answers = append(answers, r)
		log.Printf("[INFO] answer %s\n", r)
	}

	// 結果を書き込む

	if err := w.Write(answers); err != nil {
		return err
	}

	log.Println("[INFO] start to save csv file")
	w.Flush()

	if err := w.Error(); err != nil {
		return err
	}
	log.Println("[INFO] end main process")

	return nil
}

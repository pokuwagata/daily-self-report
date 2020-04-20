package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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
	Order     int     `json:"order"`
	Question  string  `json:"question"`
	InputType string  `json:"inputType"`
	Min       float64 `json:"min"`
	Max       float64 `json:"max"`
	validate  promptui.ValidateFunc
}

const (
	csvFilename = "record.csv"
	qFileName   = "question.json"
)

var(
	csvFilePath = filepath.Join(dsrPath, csvFilename)
	qFilePath = filepath.Join(dsrPath, qFileName)
)

func numParseValidate(i string) (float64, error) {
	n, err := strconv.ParseFloat(i, 64)
	if err != nil {
		return 0, err
	}
	return n, nil
}

func numValidateFunc(min float64, max float64) func(string) error {
	return func(i string) error {
		num, err := numParseValidate(i)
		if err != nil {
			return err
		}
		if num < min {
			return fmt.Errorf("%f 以上の値を入力してください", min)
		}
		if num > max {
			return fmt.Errorf("%f 以下の値を入力してください", max)
		}
		return nil
	}
}

func strValidateFunc(min int, max int) func(string) error {
	return func(i string) error {
		if len(i) < min {
			return fmt.Errorf("%d 文字以上を入力してください", min)
		}
		if len(i) > max {
			return fmt.Errorf("%d 文字以下の値を入力してください", max)
		}
		return nil
	}
}

func buildHeader(questions []question) []string {
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
	var questions = []question{}
	date := parseDate(c.option.date)
	log.Printf("[INFO] start main process: %s", date)

	// 質問テンプレートを開く
	log.Println("[INFO] open template json")
	if _, err := os.Stat(qFilePath); os.IsNotExist(err) {
		// json ファイルが存在しない場合はデフォルトの定義を利用する
		qFilePath = qFileName
	} 
	qfile, err := os.OpenFile(qFilePath, os.O_RDONLY, 0444)
	if err != nil {
		return err
	}
	defer func() {
		log.Println("[INFO] close template json")
		qfile.Close()
	}()

	qbyte, err := ioutil.ReadAll(qfile)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(qbyte, &questions); err != nil {
		return err
	}

	log.Println("[INFO] open csv file")
	file, err := os.OpenFile(csvFilePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
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

	records, err := r.ReadAll()

	if err != nil {
		return err
	}

	if records == nil {
		// ファイルが空の場合はヘッダーを作成する
		if err := w.Write(buildHeader(questions)); err != nil {
			return err
		}
	}

	for _, row := range records {
		if row[0] == date {
			return fmt.Errorf("%s は記録済みです💡", date)
		}
	}

	answers := []string{date}

	// 質問を行う
	fmt.Printf("%s について質問します🐦\n", date)
	for i, q := range questions {
		var validateFunc promptui.ValidateFunc
		switch t := q.InputType; t {
		case "number":
			validateFunc = numValidateFunc(q.Min, q.Max)
		case "string":
			validateFunc = strValidateFunc(int(q.Min), int(q.Max))
		}

		prompt := promptui.Prompt{
			Label:    q.Question,
			Validate: validateFunc,
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

# Daily self report

何か書く

## Installation

```
go get pokuwagata/daily-self-report
```

## Usage

```bash
$ dsr
昨日について質問します。
1. 睡眠時間は？（単位：時）
7

2. 食事の回数は？
3

3. ストレス量は？（1~10)
3

4. 感動の回数は？
0

5. 幸福度は？(1~10)
3

6. カフェインの量は？(1~5)
3
```

`dsr` saves to csv file.

```
date, 1, 2, 3, 4, 5, 6
2020-4-13, 5, 3, 2, 1, 3, 4
2020-4-14, 7, 3, 3, 0, 3, 3
```

`dsr` has `-d` option to specify date.

```bash
$ dsr -d 2020-4-10
2020/4/10 について質問します。
```

`dsr` reads json template for question.

```json
[
  {
    "order": 1,
    "question": "睡眠時間は？",
    "inputType": "number",
    "min": 0,
    "max": 24
  },
  {
    "order": 2,
    "question": "今日の気分は？",
    "inputType": "string",
    "min": 0,
    "max": 10
  },
  ...
]

```

## Liense

MIT

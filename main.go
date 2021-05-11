package main

import (
	"koboBookExport/models/bookList"
	"koboBookExport/models/error"
	"koboBookExport/models/highlight"

	"github.com/AlecAivazis/survey/v2"
)

type Answer struct {
	Title string `survey:"BookList"`
}

func main() {
	var qs = []*survey.Question{
		{
			Name: "BookList",
			Prompt: &survey.Select{
				Message: "選擇你想使用的功能:",
				Options: []string{"匯出書本清單", "匯出書本劃線"},
				Default: "匯出書本清單",
			},
		},
	}
	Answers := Answer{}
	err := survey.Ask(qs, &Answers)
	error.CheckErr("Answer Fail: ", err)
	if Answers.Title == "匯出書本清單" {
		bookList.ExportAction()
	} else {
		highlight.ExportAction()
	}

}

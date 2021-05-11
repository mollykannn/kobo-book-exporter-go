package highlight

import (
	"database/sql"
	"fmt"
	"koboBookExport/models/error"
	"os"

	"github.com/AlecAivazis/survey/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Answer struct {
	Title string `survey:"BookList"`
}

type BookTitle struct {
	ContentID string
	Title     string
}

func ExportAction() {
	contentID := ""
	title := ""
	db, err := sql.Open("sqlite3", "KoboReader.sqlite")
	error.CheckErr("Failed to open database: ", err)
	defer db.Close()

	var selectBookList []string
	rowBookTitle, err := db.Query(`
	SELECT 
		ContentID, 
		Title
	FROM content 
	WHERE ContentType = 6 AND ___UserId IS NOT NULL AND ___UserId != '' AND ___UserId != 'removed'
	ORDER BY DateLastRead DESC`)
	error.CheckErr("Failed to run SQL: ", err)
	defer rowBookTitle.Close()
	Books := []BookTitle{}
	for rowBookTitle.Next() {
		Book := BookTitle{}
		err := rowBookTitle.Scan(&Book.ContentID, &Book.Title)
		error.CheckErr("Failed to insert Data: ", err)
		selectBookList = append(selectBookList, Book.Title)
		Books = append(Books, Book)
	}

	var qs = []*survey.Question{
		{
			Name: "BookList",
			Prompt: &survey.Select{
				Message: "選擇你想匯出筆記的書本:",
				Options: selectBookList,
				Default: selectBookList[0],
			},
		},
	}
	Answers := Answer{}
	err = survey.Ask(qs, &Answers)
	error.CheckErr("Answer Fail: ", err)
	for i := range Books {
		if Books[i].Title == Answers.Title {
			contentID = Books[i].ContentID
			title = Books[i].Title
			break
		}
	}

	if len(contentID) > 0 && len(title) > 0 {
		highlight(contentID, title)
	}
}

func highlight(contentID, title string) {
	db, err := sql.Open("sqlite3", "KoboReader.sqlite")
	error.CheckErr("Failed to open database: ", err)
	defer db.Close()

	row, err := db.Query(fmt.Sprintf(`
	SELECT 
		TRIM(REPLACE(REPLACE(T.Text,CHAR(10),''),CHAR(9),''))
	FROM content AS B, bookmark AS T
	WHERE B.ContentID = T.VolumeID AND T.Text != '' AND T.Hidden = 'false' AND B.ContentID = '%s'
	ORDER BY T.ContentID, T.ChapterProgress;`, contentID))
	error.CheckErr("Failed to run SQL: ", err)
	defer row.Close()

	exportMarkdown(title, row)
}

func exportMarkdown(title string, row *sql.Rows) {
	f, err := os.Create("./output/" + title + ".md")
	error.CheckErr("Failed to create json: ", err)

	_, _ = f.WriteString("# " + title + "  \n\n")
	for row.Next() {
		highLight := ""
		err := row.Scan(&highLight)
		error.CheckErr("Failed to insert Data: ", err)
		_, _ = f.WriteString(highLight + "  \n\n")
	}
}

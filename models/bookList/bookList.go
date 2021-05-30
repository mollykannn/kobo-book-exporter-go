package bookList

import (
	"database/sql"
	"encoding/json"
	"koboBookExport/models/error"
	"os"
	"strconv"

	"github.com/AlecAivazis/survey/v2"
	_ "github.com/mattn/go-sqlite3"
)

type Book struct {
	BookTitle    string
	SubTitle     string
	Author       string
	Publisher    string
	ISBN         string
	ReleaseDate  string
	Series       string
	SeriesNumber int
	Rating       int
	ReadPercent  int
	LastRead     string
	FileSize     int
	Source       string
}

type Answer struct {
	Format string `survey:"BookList"`
}

func ExportAction() {
	var qs = []*survey.Question{
		{
			Name: "BookList",
			Prompt: &survey.Select{
				Message: "選擇你想匯出書本清單的格式:",
				Options: []string{"JSON", "Markdown", "CSV"},
				Default: "JSON",
			},
		},
	}
	Answers := Answer{}
	err := survey.Ask(qs, &Answers)
	error.CheckErr("Answer Fail: ", err)
	bookList(Answers.Format)
}

func bookList(action string) {
	db, err := sql.Open("sqlite3", "KoboReader.sqlite")
	error.CheckErr("Failed to open database:", err)
	defer db.Close()

	count := 0
	err = db.QueryRow(`SELECT COUNT(*) 
	FROM content
	WHERE ContentType=6 AND ___UserId IS NOT NULL AND ___UserId != '' AND ___UserId != 'removed'`).Scan(&count)
	error.CheckErr("Failed to run sql: ", err)

	row, err := db.Query(`
	SELECT
		IFNULL(Title,'') as 'Book Title', 
		IFNULL(Subtitle,''), 
		IFNULL(Attribution,'') as 'Author', 
		IFNULL(Publisher,''), 
		IFNULL(ISBN,0), 
		IFNULL(date(DateCreated),'') as 'Release Date',
		IFNULL(Series,''), 
		IFNULL(SeriesNumber,0) as 'SeriesNumber', 
		IFNULL(AverageRating,0) as 'Rating', 
		IFNULL(___PercentRead,0) as 'ReadPercent',
		IFNULL(CASE WHEN ReadStatus>0 THEN datetime(DateLastRead) END,'') as 'Last Read',
		IFNULL(___FileSize,0) as 'File Size',
		IFNULL(CASE WHEN Accessibility=1 THEN 'Store' ELSE CASE WHEN Accessibility=-1 THEN 'Import' ELSE CASE WHEN Accessibility=6 THEN 'Preview' ELSE 'Other' END END END,'') as 'Source'
	FROM content
	WHERE ContentType=6 AND ___UserId IS NOT NULL AND ___UserId != '' AND ___UserId != 'removed'
	ORDER BY Source desc, Title`)
	error.CheckErr("Failed to run sql: ", err)
	defer row.Close()

	if action == "JSON" {
		exportJSON(row, count)
	} else if action == "Markdown" {
		exportMarkdown(row)
	} else if action == "CSV" {
		exportCSV(row)
	}
}

func exportJSON(row *sql.Rows, count int) {
	f, err := os.Create("./output/BookListExport.json")
	error.CheckErr("Failed to create json: ", err)

	rowCount := 0
	for row.Next() {
		Books := Book{}
		err := row.Scan(&Books.BookTitle, &Books.SubTitle, &Books.Author, &Books.Publisher, &Books.ISBN, &Books.ReleaseDate, &Books.Series, &Books.SeriesNumber, &Books.Rating, &Books.ReadPercent, &Books.LastRead, &Books.FileSize, &Books.Source)
		error.CheckErr("Failed to insert Data: ", err)

		BookBytes, err := json.Marshal(Books)
		error.CheckErr("Failed to encode json: ", err)

		if rowCount == 0 {
			_, _ = f.WriteString("[")
		}
		_, err = f.Write(BookBytes)
		error.CheckErr("Failed to write json: ", err)

		if rowCount+1 == count {
			_, _ = f.WriteString("]")
		} else {
			_, _ = f.WriteString(",")
		}
		rowCount += 1
	}
}

func exportMarkdown(row *sql.Rows) {
	f, err := os.Create("./output/BookListExport.md")
	error.CheckErr("Failed to create markdown: ", err)
	_, _ = f.WriteString("| BookTitle | SubTitle | Author | Publisher | ISBN | ReleaseDate | Series | SeriesNumber | Rating | ReadPercent | LastRead | FileSize | Source | \n")
	_, _ = f.WriteString("| --- | --- | ---	| --- | ---	| --- | ---	| --- | ---	| --- | ---	| --- | ---	| \n")
	for row.Next() {
		Books := Book{}
		err := row.Scan(&Books.BookTitle, &Books.SubTitle, &Books.Author, &Books.Publisher, &Books.ISBN, &Books.ReleaseDate, &Books.Series, &Books.SeriesNumber, &Books.Rating, &Books.ReadPercent, &Books.LastRead, &Books.FileSize, &Books.Source)
		error.CheckErr("Failed to insert Data: ", err)
		_, _ = f.WriteString("| " + Books.BookTitle + " | " + Books.SubTitle + " | " + Books.Author + " | " + Books.Publisher + " | " + Books.ISBN + " | " + Books.ReleaseDate + " | " + Books.Series + " | " + strconv.Itoa(Books.SeriesNumber) + " | " + strconv.Itoa(Books.Rating) + " | " + strconv.Itoa(Books.ReadPercent) + " | " + Books.LastRead + " | " + strconv.Itoa(Books.FileSize) + " | " + Books.Source + " | \n")
	}
}

func exportCSV(row *sql.Rows) {
	f, err := os.Create("./output/BookListExport.csv")
	error.CheckErr("Failed to create json: ", err)
	_, _ = f.WriteString("\"BookTitle\",\"SubTitle\",\"Author\",\"Publisher\",\"ISBN\",\"ReleaseDate\",\"Series\",\"SeriesNumber\",\"Rating\",\"ReadPercent\",\"LastRead\",\"FileSize\",\"Source\" \n")
	for row.Next() {
		Books := Book{}
		err := row.Scan(&Books.BookTitle, &Books.SubTitle, &Books.Author, &Books.Publisher, &Books.ISBN, &Books.ReleaseDate, &Books.Series, &Books.SeriesNumber, &Books.Rating, &Books.ReadPercent, &Books.LastRead, &Books.FileSize, &Books.Source)
		error.CheckErr("Failed to insert Data: ", err)
		_, _ = f.WriteString("\"" + Books.BookTitle + "\",\"" + Books.SubTitle + "\",\"" + Books.Author + "\",\"" + Books.Publisher + "\",\"" + Books.ISBN + "\",\"" + Books.ReleaseDate + "\",\"" + Books.Series + "\",\"" + strconv.Itoa(Books.SeriesNumber) + "\",\"" + strconv.Itoa(Books.Rating) + "\",\"" + strconv.Itoa(Books.ReadPercent) + "\",\"" + Books.LastRead + "\",\"" + strconv.Itoa(Books.FileSize) + "\",\"" + Books.Source + "\" \n")
	}
}

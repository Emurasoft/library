package main

import (
	"bytes"
	"github.com/gocarina/gocsv"
	"os"
	"strings"
	"text/template"
)

func main() {
	file, err := os.Open("wpen13_wpfb_files.csv")
	if err != nil {
		panic(err)
	}

	var records []*Record

	if err := gocsv.UnmarshalFile(file, &records); err != nil {
		panic(err)
	}

	formatDescriptions(records)

	result := recordsToTable(records)
	if err := os.WriteFile("output.md", result, 0666); err != nil {
		panic(err)
	}
}

type Record struct {
	FilePath        string `csv:"file_path"`
	Category        string `csv:"file_category_name"`
	FileName        string `csv:"file_name"`
	DisplayName     string `csv:"file_display_name"`
	FileDescription string `csv:"file_description"`
	Author          string `csv:"file_author"`
}

func formatDescriptions(records []*Record) {
	replacer := strings.NewReplacer(
		"&nbsp;", " ",
		"\r\n", "<br>",
		"\n", "<br>",
	)
	for _, record := range records {
		record.FileDescription = replacer.Replace(record.FileDescription)
	}
}

const tableTemplate = `| Description | File | Author |
|-|-|-|
{{range . -}}
| <details><summary>{{.DisplayName}}</summary> {{.FileDescription}}</details> | {{.FileName}} | {{.Author}} |
{{end -}}
`

func recordsToTable(records []*Record) []byte {
	tableTemplate, err := template.New("table").Parse(tableTemplate)
	if err != nil {
		panic(err)
	}

	b := &bytes.Buffer{}
	if err := tableTemplate.Execute(b, &records); err != nil {
		panic(err)
	}

	return b.Bytes()
}

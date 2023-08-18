package main

import (
	"bytes"
	"cmp"
	"fmt"
	"github.com/gocarina/gocsv"
	"os"
	"slices"
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

	categories := splitCategories(records)

	text := &bytes.Buffer{}
	if _, err := text.WriteString(heading); err != nil {
		panic(err)
	}

	for _, category := range categories {
		if _, err := text.WriteString(fmt.Sprintf("# %s\n\n", category.Name)); err != nil {
			panic(err)
		}

		tableText := recordsToTable(category.Records)
		if _, err := text.Write(tableText); err != nil {
			panic(err)
		}
	}

	textBytes := text.Bytes()
	textBytes = bytes.ReplaceAll(textBytes, []byte("\n"), []byte("\r\n"))

	if err := os.WriteFile("../README.md", textBytes, 0666); err != nil {
		panic(err)
	}
}

const heading = `# Library
- [Macros](#macros)
- [Plug-ins (32-bit)](#plug-ins-32-bit)
- [Plug-ins (64-bit)](#plug-ins-64-bit)
- [Plug-ins (source code)](#plug-ins-source-code)
- [Related Software](#related-software)
- [Snippets](#snippets)
- [Syntax Files](#syntax-files)
- [Template Files](#template-files)
- [Theme](#theme)
- [Uploaded](#uploaded)

`

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
		"\u00A0", " ",
		"\r\n", "<br>",
		"\n", "<br>",
		"[", `\[`,
		"]", `\]`,
		"*", `\*`,
		"-", `\-`,
		"_", `\_`,
		"|", `\|`,
	)
	for _, record := range records {
		record.FileDescription = replacer.Replace(record.FileDescription)
	}
}

type Category struct {
	Name    string
	Records []*Record
}

func splitCategories(records []*Record) []*Category {
	categoryMap := make(map[string][]*Record)

	for _, record := range records {
		categoryMap[record.Category] = append(categoryMap[record.Category], record)
	}

	var result []*Category
	for name, records := range categoryMap {
		result = append(result, &Category{
			Name:    name,
			Records: records,
		})
	}

	slices.SortFunc(result, func(a, b *Category) int {
		return cmp.Compare(a.Name, b.Name)
	})

	return result
}

const tableTemplate = `| Description | File | Author |
|-|-|-|
{{range . -}}
| <details><summary>{{.DisplayName}}</summary> {{.FileDescription}}</details> | [{{.FileName}}](https://raw.githubusercontent.com/Emurasoft/library/main/{{.FilePath}}) | {{.Author}} |
{{end}}
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

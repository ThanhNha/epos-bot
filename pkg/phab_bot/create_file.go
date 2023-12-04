package phab_bot

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
)

type TableRevisions struct {
	Title      string
	Heading    string
	ColTitle   []string
	ColContent []TableContent
	Style      string
}
type TableContent struct {
	Name   string
	Author string
	URL    string
}

func CreateHtmlFile() {

	// Init create File
	// fileName := "example.html"
	outputPath := filepath.Join("static", "list-revision.html")
	// // Create the file
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	// Get data
	tableContent, _ := GetListRevisionsOfWeek()
	var tableRevisions TableRevisions

	colHeading := []string{
		"Name",
		"Author",
		"URL",
	}
	style := `table { border-collapse: collapse; width: 100%; } table, table th, table td { border: 1px solid #ccc; } table th, table td { padding: 0.5rem; } table th { position: relative; cursor: grab; user-select: none; }table>tbody>tr:hover{background: rgba(0, 0, 0, 0.075);}`

	tableRevisions.ColTitle = colHeading
	tableRevisions.Title = "EPOS - LIST - REVISIONS"
	tableRevisions.Heading = "Hello"
	tableRevisions.Style = style

	for _, item := range tableContent {
		tableRevisions.ColContent = append(tableRevisions.ColContent, TableContent{
			Name:   item.Name,
			Author: item.Author,
			URL:    item.URL,
		})
	}

	htmlTemplate := `
		<!DOCTYPE html>
		<html>
		<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>{{.Title}}</title>
		</head>
		<body>
		<style>{{.Heading}}</style>
			<h1>{{.Heading}}</h1>
			<table>
				<thead>
					<tr>
						{{range .ColTitle}}
						<th>{{.}}</th>
						{{end}}
					</tr>
				</thead>
				<tbody>
					{{range .ColContent}}
					<tr>
						<td>{{.Name}}</td>
						<td>{{.Author}}</td>
						<td>{{.URL}}</td>
					</tr>
					{{end}}
				</tbody>
			</table>
		</body>
		</html>
	`

	tmpl := template.New("example")
	tmpl, err = tmpl.Parse(htmlTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return
	}
	err = tmpl.Execute(file, tableRevisions)

	if err != nil {
		fmt.Println("Error executing template:", err)
		return
	}

}

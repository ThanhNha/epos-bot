package phab_bot

import (
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"time"
)

type TableRevisions struct {
	Title      string
	Heading    string
	ColTitle   []string
	ColContent []TableContent
	Date       string
}
type TableContent struct {
	Name   string
	Author string
	URL    string
}

func CreateHtmlFile(tableContent []TableContent) error {

	// Init create File
	outputPath := filepath.Join("static", "revisions.html")
	// Create the file

	file, err := os.Create(outputPath)

	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}

	defer file.Close()

	// Get data
	var tableRevisions TableRevisions

	colHeading := []string{
		"Name",
		"Author",
		"URL",
	}

	tableRevisions.ColTitle = colHeading
	tableRevisions.Title = "EPOS - LIST - REVISIONS"
	tableRevisions.Heading = "List revision on active"
	tableRevisions.Date = time.Now().Format("01/02/06")

	for _, item := range tableContent {
		tableRevisions.ColContent = append(tableRevisions.ColContent, TableContent{
			Name:   item.Name,
			Author: item.Author,
			URL:    item.URL,
		})
	}

	tmpl, err := template.ParseFiles("templates/template.html")

	if err != nil {
		fmt.Println("Error parsing template:", err)
		return err
	}

	err = tmpl.Execute(file, tableRevisions)

	if err != nil {
		fmt.Println("Error executing template:", err)
		return err
	}

	return nil

}

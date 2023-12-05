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
	Total      int
}
type TableContent struct {
	Name    string
	Author  string
	URL     string
	Status  string
	Project string
}

func CreateHtmlFile(tableContent []TableContent) error {

	// Init create File
	outputPath := filepath.Join("static", "daily-report-revisions.html")

	file, err := os.Create(outputPath)

	if err != nil {
		fmt.Println("Error creating file:", err)
		return err
	}

	if len(tableContent) == 0 {
		fmt.Println("Error creating file:", err)
		return err
	}

	defer file.Close()

	// Get data
	var tableRevisions TableRevisions

	colHeading := []string{
		"Title Revison",
		"Author",
		"URL",
		"Status",
	}

	tableRevisions.ColTitle = colHeading
	tableRevisions.Title = "EPOS - LIST - REVISIONS"
	tableRevisions.Heading = "Daily report revision on active"
	tableRevisions.Date = time.Now().Format("02/01/06")

	for _, item := range tableContent {
		tableRevisions.ColContent = append(tableRevisions.ColContent, TableContent{
			Name:   item.Name,
			Author: item.Author,
			URL:    item.URL,
			Status: item.Status,
		})
	}
	tableRevisions.Total = len(tableContent)

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

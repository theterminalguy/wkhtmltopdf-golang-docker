package main

import (
	"bytes"
	"html/template"
	"net/http"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type Person struct {
	FirstName string
	LastName  string
	JobTitle  string
	Skills    []string
}

func GeneratePDF(htmlTemplate string, person *Person) ([]byte, error) {
	var htmlBuffer bytes.Buffer
	err := template.Must(template.New("profile").
		Parse(htmlTemplate)).
		Execute(&htmlBuffer, person)
	if err != nil {
		return nil, err
	}

	// Initialize the converter
	pdfGen, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return nil, err
	}
	page := wkhtmltopdf.NewPageReader(bytes.NewReader(htmlBuffer.Bytes()))
	pdfGen.AddPage(page)

	// Set PDF page attributes
	// Refer to http://wkhtmltopdf.org/usage/wkhtmltopdf.txt for more information
	pdfGen.MarginLeft.Set(10)
	pdfGen.MarginRight.Set(10)
	pdfGen.Dpi.Set(300)
	pdfGen.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfGen.Orientation.Set(wkhtmltopdf.OrientationPortrait)

	// Generate the PDF
	err = pdfGen.Create()
	if err != nil {
		return nil, err
	}
	return pdfGen.Bytes(), nil
}

func main() {
	person := &Person{
		FirstName: "Simon Peter",
		LastName:  "Damian",
		JobTitle:  "Software Developer",

		Skills: []string{
			"Go",
			"Ruby",
			"Python",
			"JavaScript",
		},
	}
	htmlTemplate := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
		</head>
		<body>
			<h1>{{.FirstName}} {{.LastName}}</h1>
			<p>{{.JobTitle}}</p>

			<h2>Skills</h2>
			<ul>
				{{range .Skills}}
				<li>{{.}}</li>
				{{end}}
			</ul>

			<a href="/download">Download PDF</a>
		</body>
		</html>
	`

	// render the HTML template on the index route
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := template.Must(template.New("profile").
			Parse(htmlTemplate)).
			Execute(w, person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})

	// serve the PDF on the download route
	http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
		pdf, err := GeneratePDF(htmlTemplate, person)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=profile.pdf")
		w.Write(pdf)
	})
	println("listening on port", os.Getenv("PORT"))
	http.ListenAndServe(":"+os.Getenv("PORT"), nil)
}

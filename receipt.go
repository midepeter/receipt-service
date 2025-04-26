package main

import (
	"bytes"
	"fmt"
	"html/template"
	"io/ioutil"
	"os"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

// CompileTemplate compiles a template file with provided data
func CompileTemplate(templatePath string, data interface{}) (string, error) {
	// Read the template file
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("error parsing template: %v", err)
	}

	// Execute the template with data
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("error executing template: %v", err)
	}

	return buf.String(), nil
}

// GeneratePDF converts HTML content to PDF
func GeneratePDF(htmlContent, outputPath string) error {
	// Create a temporary HTML file
	tmpFile, err := ioutil.TempFile("", "template-*.html")
	if err != nil {
		return fmt.Errorf("error creating temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write HTML content to the temp file
	if _, err := tmpFile.WriteString(htmlContent); err != nil {
		return fmt.Errorf("error writing to temp file: %v", err)
	}
	tmpFile.Close()

	// Initialize PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return fmt.Errorf("error creating PDF generator: %v", err)
	}

	// Set PDF properties
	pdfg.Dpi.Set(400)
	///pdfg.Orientation.Set(wkhtmltopdf.OrientationPortrait)
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Cover.DisableSmartShrinking.Set(true)
	//pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	// Add the HTML file to the PDF generato
	page := wkhtmltopdf.NewPage(tmpFile.Name())
	page.EnableLocalFileAccess.Set(true)

	pdfg.AddPage(page)

	// Create the PDF
	if err := pdfg.Create(); err != nil {
		return fmt.Errorf("error creating PDF: %v", err)
	}

	// Write the PDF to file
	if err := pdfg.WriteFile(outputPath); err != nil {
		return fmt.Errorf("error writing PDF to file: %v", err)
	}

	return nil
}

package pdfGenerator

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

//pdf requestpdf struct
type RequestPdf struct {
	body string
}

//new request to pdf function
func NewRequestPdf(body string) *RequestPdf {
	return &RequestPdf{
		body: body,
	}
}

//parsing template function
func (r *RequestPdf) ParseTemplate(templateFileName string, data interface{}) error {

	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return err
	}
	r.body = buf.String()
	return nil
}

//generate pdf function
func (r *RequestPdf) GeneratePDF(pdfPath string) (bool, error) {

	// generate a filename (not necessary)
	t := time.Now().Unix()
	filename := "storage/" + strconv.FormatInt(int64(t), 10) + ".html"

	// write whole the body into a file
	err1 := ioutil.WriteFile(filename, []byte(r.body), 0644)
	if err1 != nil {
		panic(err1)
	}

	// read the file content
	f, err := os.Open(filename)
	if f != nil {
		defer f.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		_ = os.Remove(filename)
		log.Fatal(err)
	}
	pdfg.AddPage(wkhtmltopdf.NewPageReader(f))

	pdfg.AddPage(wkhtmltopdf.NewPageReader(strings.NewReader(r.body)))

	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)

	pdfg.Dpi.Set(300)

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	err = pdfg.WriteFile(pdfPath)
	if err != nil {
		log.Fatal(err)
	}
	_ = os.Remove(filename)

	return true, nil
}

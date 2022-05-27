package store

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	goepub "github.com/bmaupin/go-epub"

	"api/pkg/book/utils"
)

func Conv(src Store, outpath string) (err error) {
	f, err := os.Create(outpath)
	if err != nil {
		return err
	}
	defer f.Close()

	temp := template.New("txt_fiction")
	temp = temp.Funcs(template.FuncMap{
		"split": strings.Split,
	})
	temp, err = temp.Parse(TxtTemplate)
	if err != nil {
		return err
	}

	return temp.Execute(
		f, src)
}

func EPUBConv(src Store, outpath string) (err error) {
	e := goepub.NewEpub(src.BookName)
	e.SetLang("中文")
	e.SetAuthor(src.Author)

	if src.CoverURL != "" {

		body, err := utils.GetWebPageBodyReader(src.CoverURL)
		if err != nil {
			return err
		}
		tempfile, err := ioutil.TempFile("", "book_cover_*.jpg")
		if err != nil {
			return err
		}
		coverBuf, _ := ioutil.ReadAll(body)
		ioutil.WriteFile(tempfile.Name(), coverBuf, 0775)

		log.Printf("Save Cover Image: %#v", tempfile.Name())

		e.AddImage(tempfile.Name(), "cover.jpg")
		e.SetCover("cover.jpg", "")
	}

	d := ""
	dlist := strings.Split(src.Description, "\n")
	for _, cc := range dlist {
		d += fmt.Sprintf(`<p style="text-indent:2em">%s</p>`, cc)
	}
	// Description := fmt.Sprintf(`<h1><a href=%#v>%s</a></h1>%s`, src.BiqugeURL, src.BookName, d)
	Description := fmt.Sprintf(`<h1>%s</h1>%s`, src.BookName, d)
	_, err = e.AddSection(Description, "简介", "Cover.xhtml", "")
	if err != nil {
		return err
	}
	for k1, v1 := range src.Volumes {
		for k2 := range v1.Chapters {
			s := ""
			// s += fmt.Sprintf(`<h1><a href=%#v>%s</a></h1>`, v2.URL, v2.Name)
			s += fmt.Sprintf(`<h1>%s</h1>`, v1.Chapters[k2].Name)
			for _, cc := range v1.Chapters[k2].Text {
				s += fmt.Sprintf(`<p style="text-indent:2em">%s</p>`, cc)
			}
			_, err = e.AddSection(s, v1.Chapters[k2].Name, fmt.Sprintf("%d-%d.xhtml", k1, k2), "")
			if err != nil {
				return err
			}
		}
	}
	err = e.Write(outpath)
	return
}

var TxtTemplate = `书名：{{.BookName}}
作者：{{.Author}}
链接：{{.BookURL}}
简介：
{{range split .Description "\n"}}	{{.}}
{{end}}
{{- range .Volumes }}
{{if .IsVIP}}付费{{else}}免费{{end}}卷 {{.Name}}
{{range .Chapters}}
{{.Name}}
{{range .Text}}	{{.}}
{{end}}{{end}}{{end}}`

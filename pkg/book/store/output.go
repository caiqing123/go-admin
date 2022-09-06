package store

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"strings"

	goepub "github.com/bmaupin/go-epub"
	"gopkg.in/yaml.v3"

	"api/pkg/book/utils"
)

//TXTConv txt格式
func TXTConv(src Store, outpath string) (err error) {
	f, err := os.Create(outpath + ".txt")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

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

var MdTemplate = `
{{- if not .Opts.NoEPUBMetadata -}}
---
{{.EPUBMeta | yaml_marshal -}}
---
{{end -}}

书名: [{{.Store.BookName | markdown}}]({{.Store.BookURL}})

作者: {{.Store.Author}}
简介:
{{range split .Store.Description "\n" -}}
<p style="text-indent:2em">{{. | markdown}}</p>
{{end -}}
{{range .Store.Volumes }}
# {{.Name | markdown}} {{if .IsVIP}}付费{{else}}免费{{end}}卷
{{range .Chapters}}
## {{.Name | markdown}}

{{range .Text -}}
<p style="text-indent:2em">{{. | markdown}}</p>
{{end}}
{{end}}{{end}}
`

type MarkdownEPUBmeta struct {
	Title       string `yaml:"title"`
	Description string `yaml:"description,omitempty"`
	Author      string `yaml:"creator,omitempty"`
	Lang        string `yaml:"lang,omitempty"`
	Cover       string `yaml:"cover-image,omitempty"`
}

// Option is Convert output options
type Option struct {
	IgnoreCover    bool // 忽略封面
	NoEPUBMetadata bool // 不添加EPUB元数据
}

type MarkdownTemplateValues struct {
	Store    Store
	Opts     Option
	EPUBMeta MarkdownEPUBmeta
}

//EPUBConv epub格式
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
		_ = ioutil.WriteFile(tempfile.Name(), coverBuf, 0775)

		log.Printf("Save Cover Image: %#v", tempfile.Name())

		_, _ = e.AddImage(tempfile.Name(), "cover.jpg")

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
	err = e.Write(outpath + ".epub")
	return
}

//MdConv md格式
func MdConv(src Store, outpath string) (err error) {
	var (
		meta MarkdownEPUBmeta
		temp *template.Template
	)
	f, err := os.Create(outpath + ".md")
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	meta = MarkdownEPUBmeta{
		Title:       src.BookName,
		Description: src.Description,
		Author:      src.Author,
		Lang:        "zh-CN",
		Cover:       src.CoverURL,
	}

	temp = template.New("md_fiction")

	temp = temp.Funcs(template.FuncMap{
		"yaml_marshal": func(in interface{}) (string, error) {
			a, err := yaml.Marshal(in)
			return string(a), err
		},
		"split":    strings.Split,
		"markdown": MarkdownEscape,
	})

	temp, err = temp.Parse(MdTemplate)
	if err != nil {
		return err
	}

	return temp.Execute(
		f, MarkdownTemplateValues{
			Store:    src,
			Opts:     Option{IgnoreCover: false, NoEPUBMetadata: false},
			EPUBMeta: meta,
		})
}

func MarkdownEscape(s string) string {
	for _, v := range "\\!\"#$%&'()*+,./:;<=>?@[]^_`{|}~-" {
		s = strings.Replace(s, string(v), "\\"+string(v), -1)
	}
	return s
}

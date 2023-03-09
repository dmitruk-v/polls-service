package web

import (
	"html/template"
	"io"
	"path"
)

type HTMLRenderer interface {
	ShowTemplate(w io.Writer, name string, data any) error
}

type BaseHTMLRenderer struct {
	templatesDir string
}

func NewBaseHTMLRender(templatesDir string) *BaseHTMLRenderer {
	return &BaseHTMLRenderer{
		templatesDir: templatesDir,
	}
}

func (rdr *BaseHTMLRenderer) ShowTemplate(w io.Writer, name string, data any) error {
	tplFile := path.Join(rdr.templatesDir, name)
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return err
	}
	if err := tpl.Execute(w, data); err != nil {
		return err
	}
	return nil
}

type StubHTMLRenderer struct {
	ShowTemplateFn func(w io.Writer, name string, data any) error
}

func NewStubHTMLRender() *StubHTMLRenderer {
	return &StubHTMLRenderer{
		ShowTemplateFn: func(w io.Writer, name string, data any) error {
			return nil
		},
	}
}

func (rdr *StubHTMLRenderer) ShowTemplate(w io.Writer, name string, data any) error {
	return rdr.ShowTemplateFn(w, name, data)
}

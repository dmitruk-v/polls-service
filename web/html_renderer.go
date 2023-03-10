package web

import (
	"embed"
	"fmt"
	"html/template"
	"io"
)

type HTMLRenderer interface {
	ShowTemplate(w io.Writer, name string, data any) error
}

//go:embed templates
var templatesFS embed.FS

type BaseHTMLRenderer struct{}

func NewBaseHTMLRender() *BaseHTMLRenderer {
	return &BaseHTMLRenderer{}
}

func (rdr *BaseHTMLRenderer) ShowTemplate(w io.Writer, name string, data any) error {
	tpl, err := template.ParseFS(templatesFS, fmt.Sprintf("templates/%s", name))
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

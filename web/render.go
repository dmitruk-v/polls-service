package web

import (
	"html/template"
	"io"
	"path"
)

const templatesDir = "./static/templates"

func showTemplate(w io.Writer, name string, data any) error {
	tplFile := path.Join(templatesDir, name)
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return err
	}
	if err := tpl.Execute(w, data); err != nil {
		return err
	}
	return nil
}

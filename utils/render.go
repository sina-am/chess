package utils

import (
	"io"
	"path"
	"text/template"
)

type Renderer interface {
	Render(w io.Writer, name string, content any) error
}

type templateRenderer struct {
	rootPath string
	debug    bool
	tmpl     *template.Template
}

func NewTemplateRenderer(debug bool, rootPath string) (Renderer, error) {
	tmpl, err := template.ParseGlob(path.Join(rootPath, "*"))
	if err != nil {
		return nil, err
	}
	return &templateRenderer{
		debug:    debug,
		rootPath: rootPath,
		tmpl:     tmpl,
	}, nil
}

func (t *templateRenderer) Render(w io.Writer, name string, content any) error {
	if t.debug {
		tmpl, err := template.New("").ParseFiles(path.Join(t.rootPath, name), path.Join(t.rootPath, "base.html"))
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(w, "base", content)
	}
	return t.tmpl.ExecuteTemplate(w, name, content)
}

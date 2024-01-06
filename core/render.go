package core

import (
	"path"
	"text/template"

	"github.com/labstack/echo/v4"
)

type Renderer interface {
	Render(c echo.Context, name string, content map[string]any) error
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

func (t *templateRenderer) Render(c echo.Context, name string, content map[string]any) error {
	if content == nil {
		content = map[string]any{}
	}
	content["csrfToken"] = c.Get("csrf")
	if t.debug {
		tmpl, err := template.New("").ParseFiles(path.Join(t.rootPath, name), path.Join(t.rootPath, "base.html"))
		if err != nil {
			return err
		}
		return tmpl.ExecuteTemplate(c.Response().Writer, "base", content)
	}
	return t.tmpl.ExecuteTemplate(c.Response().Writer, name, content)
}

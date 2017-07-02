package helpers

import (
	"github.com/gin-gonic/gin/render"
	"html/template"
)

type Render struct {
	render.HTMLDebug
	FuncMap template.FuncMap
}

func New() Render {
	return Render{}
}

func (r Render) Instance(name string, data interface{}) render.Render {
	return render.HTML{
		Template: r.loadTemplate(),
		Name:     name,
		Data:     data,
	}
}

func (r Render) loadTemplate() *template.Template {
	if len(r.Files) > 0 {
		return template.Must(template.New("").Funcs(r.FuncMap).ParseFiles(r.Files...))
	}
	if len(r.Glob) > 0 {
		return template.Must(template.New("").Funcs(r.FuncMap).ParseGlob(r.Glob))
	}
	panic("the HTML debug render was created without files or glob pattern")
}

package templates

import (
	"bytes"
	"fmt"
	"github.com/itelman/forum/internal/dto"
	"net/http"
	"time"
)

const (
	AuthenticatedUser = "AuthenticatedUser"
	CurrentYear       = "CurrentYear"
	Flash             = "Flash"
	Error             = "Error"
	Form              = "Form"
	Posts             = "Posts"
	Post              = "Post"
	Categories        = "Categories"
	Comments          = "Comments"
	PostReactions     = "PostReactions"
	Comment           = "Comment"
)

type TemplateData map[string]any

type TemplateRender interface {
	RenderData(w http.ResponseWriter, r *http.Request, tmplName string, td TemplateData) error
}

type templateRender struct {
	templateCache TemplateCache
}

func NewTemplateRender(templateCache TemplateCache) *templateRender {
	return &templateRender{templateCache: templateCache}
}

func (tr *templateRender) RenderData(w http.ResponseWriter, r *http.Request, tmplName string, td TemplateData) error {
	ts, ok := tr.templateCache[tmplName]
	if !ok {
		return fmt.Errorf("TEMPLATE CACHE: template not found (%s)", tmplName)
	}

	addDefaultData(r, td)

	buf := new(bytes.Buffer)
	err := ts.Execute(buf, td)
	if err != nil {
		return err
	}

	_, err = buf.WriteTo(w)
	if err != nil {
		return err
	}

	return nil
}

func addDefaultData(r *http.Request, td TemplateData) {
	td[AuthenticatedUser] = dto.GetAuthUser(r)
	td[CurrentYear] = time.Now().Year()
}

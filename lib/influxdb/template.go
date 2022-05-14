package influxdb

import (
	"bytes"
	"text/template"

	"github.com/pkg/errors"
)

type QueryTemplate struct {
	name           string
	templateString string
	template       *template.Template
}

func NewQueryTemplate(name string, tmpl string) *QueryTemplate {
	return &QueryTemplate{
		name:           name,
		templateString: tmpl,
	}
}

func (q *QueryTemplate) getTemplate() (*template.Template, error) {
	if q.template != nil {
		return q.template, nil
	}
	tmpl, err := template.New(q.name).Parse(q.templateString)
	if err != nil {
		return nil, errors.Wrap(err, "can't parse template")
	}
	q.template = tmpl
	return q.template, nil
}

func (q *QueryTemplate) GetQueryString(args any) (string, error) {
	tmpl, err := q.getTemplate()
	if err != nil {
		return "", errors.Wrap(err, "can't get template")
	}
	buf := new(bytes.Buffer)
	err = tmpl.Execute(buf, args)
	if err != nil {
		return "", errors.Wrap(err, "can't execute template with given args")
	}
	query := buf.String()
	return query, nil
}

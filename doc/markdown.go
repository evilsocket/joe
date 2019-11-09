package doc

import (
	"encoding/json"
	"fmt"
	"github.com/evilsocket/islazy/log"
	"github.com/evilsocket/joe/api"
	"github.com/evilsocket/joe/doc/templates"
	"github.com/evilsocket/joe/models"
	"os"
	"text/template"
)

var funcMap = template.FuncMap{
	"json": func(arg interface{}) string {
		if raw, err := json.MarshalIndent(arg, "", "  "); err != nil {
			return fmt.Sprintf("template error: %v", err)
		} else {
			return string(raw)
		}
	},
}

func tpl(name string) (*template.Template, error) {
	raw, err := templates.Asset(name)
	if err != nil {
		return nil, err
	}

	t, err := template.New(name).Funcs(funcMap).Parse(string(raw))
	if err != nil {
		return nil, err
	}

	return t, nil
}

func ToMarkdown(address, fileName string) (err error) {
	log.Info("generating markdown documentation for %d queries to %s ...", models.NumQueries, fileName)
	t, err := tpl("doc/templates/doc.md")
	if err != nil {
		return
	}

	queries := make([]*models.Query, 0)
	models.Queries.Range(func(key, value interface{}) bool {
		queries = append(queries, value.(*models.Query))
		return true
	})

	out, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer out.Close()

	return t.Execute(out, map[string]interface{}{
		"Address": address,
		"Version": api.Version,
		"Queries": queries,
	})
}

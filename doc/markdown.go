package doc

import (
	"encoding/json"
	"fmt"
	"github.com/evilsocket/islazy/log"
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

func ToMarkdown(fileName string) (err error) {
	log.Info("generating markdown documentation for %d queries to %s ...", models.NumQueries, fileName)

	raw, err := templates.Asset("doc/templates/doc.md")
	if err != nil {
		return
	}

	queries := make([]*models.Query, 0)
	models.Queries.Range(func(key, value interface{}) bool {
		queries = append(queries, value.(*models.Query))
		return true
	})

	t, err := template.New("template").Funcs(funcMap).Parse(string(raw))
	if err != nil {
		return
	}

	out, err := os.Create(fileName)
	if err != nil {
		return
	}
	defer out.Close()

	return t.Execute(out, map[string]interface{}{"Queries": queries})
}

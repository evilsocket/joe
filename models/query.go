package models

import (
	"database/sql"
	"fmt"
	"github.com/evilsocket/islazy/log"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/djherbis/times.v1"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"path"
	"regexp"
	"strings"
	"sync"
	"time"
)

var paramParser = regexp.MustCompile("\\{([^}]+)\\}")

type Query struct {
	sync.Mutex

	CreatedAt   time.Time              `yaml:"-" json:"created_at"`
	UpdatedAt   time.Time              `yaml:"-" json:"updated_at"`
	Name        string                 `yaml:"-" json:"name"`
	Cache       *CachePolicy           `yaml:"cache" json:"cache"`
	Description string                 `yaml:"description" json:"description"`
	Defaults    map[string]interface{} `yaml:"defaults" json:"defaults"`
	Expression  string                 `yaml:"expression" json:"expression"`

	fileName   string
	statement  string
	parameters map[string]int
	numParams  int
	prepared   *sql.Stmt
	explainer  *Query
}

func LoadQuery(fileName string) (*Query, error) {
	log.Debug("loading %s ...", fileName)

	query := &Query{
		fileName:   fileName,
		parameters: make(map[string]int),
	}

	if raw, err := ioutil.ReadFile(fileName); err != nil {
		return nil, err
	} else if err = yaml.Unmarshal(raw, query); err != nil {
		return nil, err
	} else if err = query.load(); err != nil {
		return nil, err
	}
	return query, nil
}

func (q *Query) prepare() (err error) {
	q.statement = q.Expression
	for _, match := range paramParser.FindAllStringSubmatch(q.Expression, -1) {
		tok, par := match[0], match[1]
		if _, found := q.parameters[par]; found {
			return fmt.Errorf("token %s has been used more than once", tok)
		} else {
			q.parameters[par] = q.numParams
			q.numParams++
			q.statement = strings.Replace(q.statement, tok, "?", 1)
		}
	}

	if q.prepared, err = DB.Prepare(q.statement); err != nil {
		return fmt.Errorf("error preparing statement for %s: %v", q.Name, err)
	}

	return
}

func (q *Query) load() error {
	if t, err := times.Stat(q.fileName); err != nil {
		return err
	} else if raw, err := ioutil.ReadFile(q.fileName); err != nil {
		return err
	} else if err = yaml.Unmarshal(raw, q); err != nil {
		return err
	} else {
		q.CreatedAt = t.BirthTime()
		q.UpdatedAt = t.ModTime()
		q.Name = strings.ReplaceAll(path.Base(q.fileName), ".yml", "")

		if err := q.prepare(); err != nil {
			return err
		}

		explain := fmt.Sprintf("EXPLAIN %s", q.Expression)
		q.explainer = &Query{
			Expression: explain,
			Cache:      &CachePolicy{Type: None},
			statement:  explain,
			parameters: make(map[string]int),
		}

		if err := q.explainer.prepare(); err != nil {
			return err
		}

		log.Debug("loaded %v", q)
	}
	return nil
}

func (q *Query) toQueryArgs(params map[string]interface{}) ([]interface{}, error) {
	// assign statement parameters
	args := make([]interface{}, q.numParams)
	for name, value := range params {
		if order, found := q.parameters[name]; !found {
			return nil, fmt.Errorf("unknown parameter '%s'", name)
		} else {
			args[order] = value
		}
	}

	// assign missing
	for name, defValue := range q.Defaults {
		if _, found := params[name]; !found {
			if order, found := q.parameters[name]; !found {
				return nil, fmt.Errorf("unknown parameter '%s'", name)
			} else {
				args[order] = defValue
			}
		}
	}

	return args, nil
}

func (q *Query) Query(params map[string]interface{}) (*Rows, error) {
	log.Debug("running '%s' with %s", q.statement, params)

	begin := time.Now()

	if cached := q.Cache.Get(params); cached != nil {
		rows := cached.Data.(*Rows)
		rows.CachedAt = &cached.At
		rows.ExecutionTime = time.Since(begin)
		return rows, nil
	}

	q.Lock()
	defer q.Unlock()

	args, err := q.toQueryArgs(params)
	if err != nil {
		return nil, err
	}

	dbRows, err := q.prepared.Query(args...)
	if err != nil {
		return nil, err
	}
	defer dbRows.Close()

	rows := &Rows{
		Rows: make([]Row, 0),
	}

	if rows.ColumnNames, err = dbRows.Columns(); err != nil {
		return nil, err
	} else {
		rows.NumColumns = len(rows.ColumnNames)

		for dbRows.Next() {
			columnValues := make([]interface{}, rows.NumColumns)
			for idx, _ := range columnValues {
				var dummy interface{}
				columnValues[idx] = &dummy
			}

			if err := dbRows.Scan(columnValues...); err != nil {
				return nil, err
			}

			row := make(Row)
			for idx, name := range rows.ColumnNames {
				v := *(columnValues[idx].(*interface{}))
				if rawBytes, ok := v.([]uint8); ok {
					row[name] = string(rawBytes)
				} else {
					row[name] = v
				}
			}

			rows.Rows = append(rows.Rows, row)
			rows.NumRows++
		}
	}

	rows.ExecutionTime = time.Since(begin)

	q.Cache.Set(params, rows)

	return rows, nil
}

func (q *Query) Explain(params map[string]interface{}) (*Rows, error) {
	return q.explainer.Query(params)
}

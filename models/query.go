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
	Expression  string                 `yaml:"query" json:"query"`
	Views       map[string]string      `yaml:"views" json:"views"`
	Access      []string               `yaml:"access" json:"access"`

	Parameters map[string]string `yaml:"-" json:"-"`

	access        map[string]*User
	views         map[string]*View
	compiledViews bool
	fileName      string
	statement     string
	parameters    map[string]int
	numParams     int
	prepared      *sql.Stmt
	explainer     *Query
}

func LoadQuery(fileName string, compileViews bool) (*Query, error) {
	log.Debug("loading %s ...", fileName)

	query := &Query{
		fileName:      fileName,
		Parameters:    make(map[string]string),
		parameters:    make(map[string]int),
		access:        make(map[string]*User),
		views:         make(map[string]*View),
		compiledViews: compileViews,
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

func (q *Query) QueryString() string {
	parts := []string{}
	for name, value := range q.Parameters {
		if value == "" {
			parts = append(parts, fmt.Sprintf("%s=VALUE", name))
		} else {
			parts = append(parts, fmt.Sprintf("%s=%s", name, value))
		}
	}

	if len(parts) > 0 {
		return fmt.Sprintf("?%s", strings.Join(parts, "&"))
	} else {
		return ""
	}
}

func (q *Query) prepare() (err error) {
	q.statement = q.Expression
	for _, match := range paramParser.FindAllStringSubmatch(q.Expression, -1) {
		tok, par := match[0], match[1]
		if _, found := q.parameters[par]; found {
			return fmt.Errorf("token %s has been used more than once", tok)
		} else {
			if def, found := q.Defaults[par]; found {
				q.Parameters[par] = fmt.Sprintf("%v", def)
			} else {
				q.Parameters[par] = ""
			}
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
	} else if len(q.Access) == 0 {
		return fmt.Errorf("%s doens't declare an access section", q.Name)
	} else {
		if t.HasBirthTime() {
			q.CreatedAt = t.BirthTime()
		} else {
			q.CreatedAt = time.Now()
		}
		q.UpdatedAt = t.ModTime()
		q.Name = strings.ReplaceAll(path.Base(q.fileName), ".yml", "")

		for _, username := range q.Access {
			if u, found := Users.Load(username); !found {
				return fmt.Errorf("user %s not found", username)
			} else {
				q.access[username] = u.(*User)
				if username == "anonymous" {
					log.Warning("query %s allows anonymous access", q.Name)
				}
			}
		}

		// prepare the main statement
		if err := q.prepare(); err != nil {
			return err
		}

		// prepare the explain statement
		explain := fmt.Sprintf("EXPLAIN %s", q.Expression)
		q.explainer = &Query{
			Expression: explain,
			Cache:      &CachePolicy{Type: None},
			statement:  explain,
			parameters: make(map[string]int),
			Parameters: make(map[string]string),
		}

		if err := q.explainer.prepare(); err != nil {
			return err
		}

		// load views
		for viewName, viewFileName := range q.Views {
			if viewFileName != "" {
				if viewFileName[0] != '/' && viewFileName[0] != '.' {
					viewFileName = path.Join(path.Dir(q.fileName), viewFileName)
				}
			}
			if view, err := PrepareView(q.Name, viewName, viewFileName, q.compiledViews); err != nil {
				return fmt.Errorf("%s: %v", viewName, err)
			} else {
				q.views[viewName] = view
			}
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

func (q *Query) AuthRequired() bool {
	// does allow anonymous access?
	if _, found := q.access["anonymous"]; found {
		return false
	}
	return true
}

func (q *Query) Authorized(user *User) bool {
	if !q.AuthRequired() {
		return true
	}
	// sanity check
	if user == nil {
		return false
	}
	// check if included
	if _, found := q.access[user.Username]; found {
		return true
	}
	return false
}

func (q *Query) Query(params map[string]interface{}) (*Results, error) {
	log.Debug("running '%s' with %s", q.statement, params)

	begin := time.Now()

	if cached := q.Cache.Get(params); cached != nil {
		rows := cached.Data.(*Results)
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

	rows := &Results{
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

func (q *Query) View(name string) *View {
	if v, found := q.views[name]; found {
		return v
	}
	return nil
}

func (q *Query) Explain(params map[string]interface{}) (*Results, error) {
	return q.explainer.Query(params)
}

# Joe

Joe is a tool to automatically build a REST API around SQL queries and their resulting data. 

In a way it is an anti-[ORM](https://en.wikipedia.org/wiki/Object-relational_mapping): its purpose is to help backend engineers versioning
, annotating, exposing and charting the queries that can't be implemented or aren't worth implementing in the backend
 main business logic and that they would normally keep on .txt or .sql files.

**WORK IN PROGRESS, DO NOT USE IN PRODUCTION**

## How to Install

```sh
go get -u https://github.com/evilsocket/joe/cmd/joe
```

## Example

First create an `/etc/joe/joe.conf` configuration file with the access credentials for the database:

```conf
# CHANGE THIS: use a complex secret for the JWT token generation
API_SECRET=02zygnJs5e0bBLJjaHCinWTjfRdheTYO

DB_HOST=joe-mysql
DB_DRIVER=mysql
DB_USER=joe
DB_PASSWORD=joe
DB_NAME=joe
DB_PORT=3306
```

Then create the `admin` user (this command will generate the file `/etc/joe/users/admin.yml`):

```sh
sudo mkdir -p /etc/joe/users
sudo joe -new-user admin -token-ttl 6 # JWT tokens for this user expire after 6 hours
```

The next step is creating a `/etc/joe/queries/example.yml` file with our first example query (this query
 selects the top players [for this project](https://github.com/evilsocket/pwngrid)):

```yaml
description: "Top players by access points."

# who can access this? use 'anonymous' to allow unauthenticated access
access:
  - admin

# optional cache
cache:
  # 0 = None, 1 = by keys (+ optional ttl), 2 = by ttl
  type: 1
  # expression parameters to use for caching
  keys: [limit]
  # ttl in seconds
  ttl: 30

# the query itself
query:
  SELECT  u.updated_at as active_at, u.name, u.fingerprint, u.country, COUNT(a.id) AS networks FROM units u
  INNER JOIN access_points a ON u.id = a.unit_id GROUP BY u.id ORDER BY networks DESC LIMIT {limit}

# optional default values for the parameters
defaults:
  limit: 25

# define a chart
views:
  bars: example_view.go
```

The `example_view.go` is a view plugin that generates a chart from the selected records, by defining these views joe
 can generate PNG or SVG charts:

```go
package main

import (
	"github.com/wcharczuk/go-chart"
	"github.com/evilsocket/joe/models"
)

func View(res *models.Results) models.Chart {
	ch := chart.BarChart{
		Title: "Top Players",
		TitleStyle: chart.Style {
			Hidden: false,
		},
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Bars: make([]chart.Value, res.NumRows),
	}

	for i, row := range res.Rows {
		ch.Bars[i] = chart.Value{
			Label: row["name"].(string),
			Value: float64(row["networks"].(int64)),
		}
	}

	return ch
}
```

Now you can just start joe with:

```sh
joe -conf /etc/joe/joe.conf -data /etc/joe/queries -users /etc/joe/users
```
    
This will load the queries, compile the views and expose the following API endpoints automatically:

* http://localhost:8080/api/v1/auth?user=admin&pass=somecomplexpasswordhere this will allow users to authenticate to
 the API with username and password credentials. The endpoint will return a JWT `token` that must be passed to all
  other endpoints for authentication.
* http://localhost:8080/api/v1/query/example.json?limit=20 to get JSON data.
* http://localhost:8080/api/v1/query/example.csv?limit=20 to get CSV data.
* http://localhost:8080/api/v1/query/example/explain to explain the query.
* http://localhost:8080/api/v1/query/example/bars.png to get a PNG chart.
* http://localhost:8080/api/v1/query/example/bars.svg to get a SVG chart.

While the JWT `token` can be passed to each endpoint either as a GET or POST parameter, it is recommended to pass it via
the `Authorization` header (`Authorization: bearer put-the-token-here`).

Both the authentication process and the `token` parameter itself are optional for queries that have the user `anonymous
` in their
 `access` list.

## License

`joe` is made with ♥  by [@evilsocket](https://twitter.com/evilsocket) and it is released under the GPL3 license.
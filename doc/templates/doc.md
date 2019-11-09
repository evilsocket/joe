# Joe Queries

This document has been automatically generated with [joe v1.0.0](https://github.com/evilsocket/joe).

## Main

These are the authentication and visualization generic endpoints.

### GET|POST /api/v1/auth  

Authenticate to the API with username and password in order to get a JWT token.

|Parameter|Default|
|--|--|
| `user` | _none_ |
| `pass` | _none_ |

#### Request

    curl --data "user=admin&pass=admin" http://localhost:8080/api/v1/auth

#### Response

```json
{
  "token": "..."
}
```

### GET /api/v1/queries/

Get a list of the queries that are currently loaded in the system.

#### Request

    curl -H "Authorization: Bearer ..token.." http://localhost:8080/api/v1/queries

#### Response

```json
[
  {
    "created_at": "2019-11-09T10:57:01.784331255+01:00",
    "updated_at": "2019-11-09T10:57:01.784494235+01:00",
    "name": "top_players",
    "cache": {
      "type": 1,
      "keys": [
        "limit"
      ],
      "ttl": 3600
    },
    "description": "Top players by access points.",
    "defaults": {
      "limit": 25
    },
    "query": "SELECT  u.updated_at as active_at, u.name, u.fingerprint, u.country, COUNT(a.id) AS networks FROM units u INNER JOIN access_points a ON u.id = a.unit_id GROUP BY u.id ORDER BY networks DESC LIMIT {limit}",
    "views": {
      "bars": "top_players_bars.go"
    },
    "access": [
      "admin"
    ]
  },
  ...
}
```
{{range .Queries}}

### GET /api/v1/query/{{.Name}}/view  

Show information about the {{.Name}} query.

#### Request

{{if .AuthRequired }}
    curl -H "Authorization: Bearer ..token.." http://localhost:8080/api/v1/query/{{.Name}}/view
{{else}}
    curl http://localhost:8080/api/v1/query/{{.Name}}/view
{{end}}

#### Response

```json
{{json .}}
```

### GET|POST /api/v1/query/{{.Name}}(.json|.csv)

{{.Description}}

{{if .Parameters }}
|Parameter|Default|
|--|--|
{{range $name, $def := .Parameters }}| `{{$name}}` | {{if $def}}{{$def}}{{else}}_none_{{end}} |
{{end}}
{{end}}

#### Request

{{if .AuthRequired }}
    curl -H "Authorization: Bearer ..token.." http://localhost:8080/api/v1/query/{{.Name}}.json{{.QueryString}}
{{else}}
    curl http://localhost:8080/api/v1/query/{{.Name}}.json
{{end}}

#### Response

```json
{{json .}}
```

### GET|POST /api/v1/query/{{.Name}}/explain

Return results for an EXPLAIN operation on the {{.Name}} main query.

{{if .Parameters }}
|Parameter|Default|
|--|--|
{{range $name, $def := .Parameters }}| `{{$name}}` | {{if $def}}{{$def}}{{else}}_none_{{end}} |
{{end}}
{{end}}

#### Request

{{if .AuthRequired }}
    curl -H "Authorization: Bearer ..token.." http://localhost:8080/api/v1/query/{{.Name}}/explain{{.QueryString}}
{{else}}
    curl http://localhost:8080/api/v1/query/{{.Name}}/explain{{.QueryString}}
{{end}}

#### Response

```json
{
  "cached_at": null,
  "exec_time": 2763263,
  "num_records": 1,
  "records": [
    {
      "Extra": "Zero limit",
      "filtered": null,
      "id": 1,
      "key": null,
      "key_len": null,
      "partitions": null,
      "possible_keys": null,
      "ref": null,
      "rows": null,
      "select_type": "SIMPLE",
      "table": null,
      "type": null
    }
  ]
}
```
 
{{ $queryName := .Name }}
{{ $queryParams := .Parameters }}
{{ $queryString := .QueryString }}
{{ $queryAuthRequired := .AuthRequired }}

{{range $viewName, $viewFile := .Views }}
### GET /api/v1/query/{{$queryName}}/{{$viewName}}.(png|svg)

Return a PNG or SVG representation of a {{$viewName}} chart for the {{$queryName}} query.

{{if $queryParams }}
|Parameter|Default|
|--|--|
{{range $name, $def := $queryParams }}| `{{$name}}` | {{if $def}}{{$def}}{{else}}_none_{{end}} |
{{end}}
{{end}}

#### Request

{{if $queryAuthRequired }}
    curl -H "Authorization: Bearer ..token.." http://localhost:8080/api/v1/query/{{$queryName}}/{{$viewName}}.png{{$queryString}}
{{else}}
    curl http://localhost:8080/api/v1/query/{{$queryName}}/{{$viewName}}.png{{$queryString}}
{{end}}

{{end}}

{{end}}
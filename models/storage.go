package models

import "sync"

var (
	Queries    = sync.Map{}
	NumQueries = 0
)

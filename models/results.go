package models

import "time"

type Row map[string]interface{}

type Results struct {
	CachedAt      *time.Time    `json:"cached_at"`
	ExecutionTime time.Duration `json:"exec_time"`
	ColumnNames   []string      `json:"-"`
	NumColumns    int           `json:"-"`
	NumRows       int           `json:"num_records"`
	Rows          []Row         `json:"records"`
}
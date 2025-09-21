package store

import "sync"

type Store struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

var Clients_sse_stores sync.Map
var once_stores sync.Once

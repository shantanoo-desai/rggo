package main

import (
	"encoding/json"
	"time"

	"github.com/shantanoo-desai/rggo/todo"
)

type todoResponse struct {
	Results todo.List `json:"results"`
}

func (r *todoResponse) MarshalJSON() ([]byte, error) {
	resp := struct {
		Results	todo.List	`json:"results"`
		Date	int64		`json:"data"`
		TotalResults int	`json:"total_results"`
	}{
		Results:	r.Results,
		Date:		time.Now().Unix(),
		TotalResults: len(r.Results),
	}

	return json.Marshal(resp)
}

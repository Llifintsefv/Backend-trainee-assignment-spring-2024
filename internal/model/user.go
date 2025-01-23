package model

import "time"

type User struct {
	id         string   `json:"id"`
	username   string   `json:"username"`
	first_name string   `json:"first_name"`
	last_name  string   `json:"last_name"`
	created_at time.Time `json:"created_at"`
	updated_at time.Time   `json:"updated_at"`
}


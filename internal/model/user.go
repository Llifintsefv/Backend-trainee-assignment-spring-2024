package model

import "time"

type User struct {
	Id         string   `json:"id"`
	Username   string   `json:"username"`
	First_name string   `json:"first_name"`
	Last_name  string   `json:"last_name"`
	Created_at time.Time `json:"created_at"`
	Updated_at time.Time   `json:"updated_at"`
}


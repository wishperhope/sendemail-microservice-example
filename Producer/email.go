package main

import "time"

// Email Struct
type Email struct {
	To       string `json:"to"`
	From     string `json:"from"`
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Note     string `json:"note"`
}

// Job Db Struct
type Job struct {
	ID        string    `json:"id"`
	To        string    `json:"to"`
	From      string    `json:"from"`
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    string    `json:"status"`
	Note      string    `json:"note"`
}

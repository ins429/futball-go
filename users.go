package main

import ()

type User struct {
	Id        int64         `json:"id"`
	Username  string        `json:"username"`
	Password  string        `json:"password,omitempty"`
	FirstName string        `json:"firstname"`
	LastName  string        `json:"lastname"`
	Players   []UserPlayers `json:"players,omitempty"`
}

type UserPlayers []string

type UserResponse struct {
	Status int  `json:"status"`
	User   User `json:"user"`
}

type UserForm struct {
	Username string `form:"username"`
	Password string `form:"password,omitempty"`
}

type FbForm struct {
	Token string `form:"token"`
}

type FbUser struct {
	Id        string `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

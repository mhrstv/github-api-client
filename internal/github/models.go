package github

import "time"

type User struct {
	Login       string    `json:"login"`
	Name        string    `json:"name"`
	Followers   int       `json:"followers"`
	PublicRepos int       `json:"public_repos"`
	CreatedAt   time.Time `json:"created_at"`
}

type Repo struct {
	Name       string    `json:"name"`
	Fork       bool      `json:"fork"`
	ForksCount int       `json:"forks_count"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Languages map[string]int

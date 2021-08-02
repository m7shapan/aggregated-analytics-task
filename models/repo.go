package models

type Repo struct {
	Name    string
	Watch   int
	Commits int
}

type RepoCommit struct {
	Name    string
	Commits int
}

type RepoWatch struct {
	Name  string
	Watch int
}

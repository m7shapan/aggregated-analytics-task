package repositories

import (
	"sort"

	"github.com/m7shapan/aggregated-analytics-task/models"
)

type RepoRepository interface {
	GetTopReposByCommits(int) []models.RepoCommit
	GetTopReposByWatch(int) []models.RepoWatch
}

type repoRepository struct {
	repos []*models.Repo
}

func NewRepoRepository(repos []*models.Repo) RepoRepository {
	return &repoRepository{
		repos: repos,
	}
}

func (r repoRepository) GetTopReposByCommits(x int) (repos []models.RepoCommit) {
	if x > len(r.repos) {
		x = len(r.repos)
	}

	sort.Slice(r.repos, func(i, j int) bool {
		return r.repos[i].Commits > r.repos[j].Commits
	})

	for i := 0; i < x; i++ {
		repos = append(repos, models.RepoCommit{
			Name:    r.repos[i].Name,
			Commits: r.repos[i].Commits,
		})
	}

	return repos
}

func (r repoRepository) GetTopReposByWatch(x int) (repos []models.RepoWatch) {
	if x > len(r.repos) {
		x = len(r.repos)
	}

	sort.Slice(r.repos, func(i, j int) bool {
		return r.repos[i].Watch > r.repos[j].Watch
	})

	for i := 0; i < x; i++ {
		repos = append(repos, models.RepoWatch{
			Name:  r.repos[i].Name,
			Watch: r.repos[i].Watch,
		})
	}

	return repos
}

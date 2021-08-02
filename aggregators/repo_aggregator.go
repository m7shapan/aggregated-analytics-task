package aggregators

import (
	"io"

	"github.com/m7shapan/aggregated-analytics-task/models"
	"github.com/m7shapan/aggregated-analytics-task/storage"
)

type RepoAggregator interface {
	Aggregate() ([]*models.Repo, error)
}

type repoAggregator struct {
	repos   storage.Reader
	events  storage.Reader
	commits storage.Reader
}

func NewRepoAggregator(repos storage.Reader, events storage.Reader, commits storage.Reader) RepoAggregator {
	a := repoAggregator{
		repos:   repos,
		events:  events,
		commits: commits,
	}

	a.removeHeaders()

	return &a
}

func (a *repoAggregator) removeHeaders() {
	_, _ = a.repos.Read()
	_, _ = a.events.Read()
	_, _ = a.commits.Read()
}

func (a repoAggregator) Aggregate() ([]*models.Repo, error) {
	var reposMap = make(map[string]*models.Repo, 0)
	var reposList = make([]*models.Repo, 0)

	var commit []string
	for {
		event, err := a.events.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		repo, err := a.repos.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		var commitCount int
		if event[1] == "PushEvent" {
			if commit != nil {
				commitCount++
			}

			for {
				var err error
				commit, err = a.commits.Read()
				if err != nil {
					if err == io.EOF {
						break
					}

					return nil, err
				}

				if commit[2] != event[0] {
					break
				}
				commitCount++
			}

			if r, found := reposMap[repo[1]]; found {
				r.Commits += commitCount
				continue
			}

			a := &models.Repo{
				Name:    repo[1],
				Commits: commitCount,
			}
			reposMap[repo[1]] = a
			reposList = append(reposList, a)
			continue
		}
	}

	return reposList, nil
}

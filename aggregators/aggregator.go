package aggregators

import (
	"io"

	"github.com/m7shapan/aggregated-analytics-task/models"
	"github.com/m7shapan/aggregated-analytics-task/storage"
)

type Aggregator interface {
	AggregateAll() ([]*models.Actor, []*models.Repo, error)
}

type aggregator struct {
	actors  storage.Reader
	repos   storage.Reader
	events  storage.Reader
	commits storage.Reader
}

func NewAggregator(actors storage.Reader, repos storage.Reader, events storage.Reader, commits storage.Reader) Aggregator {
	a := aggregator{
		actors:  actors,
		repos:   repos,
		events:  events,
		commits: commits,
	}

	a.removeHeaders()

	return &a
}

func (a *aggregator) removeHeaders() {
	_, _ = a.actors.Read()
	_, _ = a.repos.Read()
	_, _ = a.events.Read()
	_, _ = a.commits.Read()
}

func (a aggregator) AggregateAll() ([]*models.Actor, []*models.Repo, error) {
	var actorsMap = make(map[string]*models.Actor, 0)
	var actorsList = make([]*models.Actor, 0)

	var reposMap = make(map[string]*models.Repo, 0)
	var reposList = make([]*models.Repo, 0)

	var commit []string
	for {
		event, err := a.events.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, nil, err
		}

		actor, err := a.actors.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, nil, err
		}

		repo, err := a.repos.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, nil, err
		}

		if event[1] == "PullRequestEvent" {
			if a, found := actorsMap[actor[1]]; found {
				a.PullRequests++
				continue
			}

			a := &models.Actor{
				Name:         actor[1],
				PullRequests: 1,
			}
			actorsMap[actor[1]] = a
			actorsList = append(actorsList, a)
			continue
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

					return nil, nil, err
				}

				if commit[2] != event[0] {
					break
				}
				commitCount++
			}

			if a, found := actorsMap[actor[1]]; found {
				a.Commits += commitCount
			} else {
				a := &models.Actor{
					Name:    actor[1],
					Commits: commitCount,
				}
				actorsMap[actor[1]] = a
				actorsList = append(actorsList, a)
			}

			if r, found := reposMap[repo[1]]; found {
				r.Commits += commitCount
			} else {
				r := &models.Repo{
					Name:    repo[1],
					Commits: commitCount,
				}
				reposMap[repo[1]] = r
				reposList = append(reposList, r)
			}

			continue
		}

		if event[1] == "WatchEvent" {
			if r, found := reposMap[repo[1]]; found {
				r.Watch++
				continue
			}

			a := &models.Repo{
				Name:  repo[1],
				Watch: 1,
			}
			reposMap[repo[1]] = a
			reposList = append(reposList, a)
			continue
		}
	}

	return actorsList, reposList, nil
}

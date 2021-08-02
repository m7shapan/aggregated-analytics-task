package aggregators

import (
	"io"

	"github.com/m7shapan/aggregated-analytics-task/models"
	"github.com/m7shapan/aggregated-analytics-task/storage"
)

type ActorAggregator interface {
	Aggregate() ([]*models.Actor, error)
}

type actorAggregator struct {
	actors  storage.Reader
	events  storage.Reader
	commits storage.Reader
}

func NewActorAggregator(actors storage.Reader, events storage.Reader, commits storage.Reader) ActorAggregator {
	a := actorAggregator{
		actors:  actors,
		events:  events,
		commits: commits,
	}

	a.removeHeaders()

	return &a
}

func (a actorAggregator) removeHeaders() {
	_, _ = a.actors.Read()
	_, _ = a.events.Read()
	_, _ = a.commits.Read()
}

func (a actorAggregator) Aggregate() ([]*models.Actor, error) {
	var actorsMap = make(map[string]*models.Actor, 0)
	var actorsList = make([]*models.Actor, 0)

	var commit []string
	for {
		event, err := a.events.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		actor, err := a.actors.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
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

					return nil, err
				}

				if commit[2] != event[0] {
					break
				}
				commitCount++
			}

			if a, found := actorsMap[actor[1]]; found {
				a.Commits += commitCount
				continue
			}

			a := &models.Actor{
				Name:    actor[1],
				Commits: commitCount,
			}
			actorsMap[actor[1]] = a
			actorsList = append(actorsList, a)
			continue
		}
	}

	return actorsList, nil
}

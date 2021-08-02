package repositories

import (
	"sort"

	"github.com/m7shapan/aggregated-analytics-task/models"
)

type ActorRepository interface {
	GetTopActorsByPRAndCommits(int) []models.Actor
}

type actorRepository struct {
	actors []*models.Actor
}

func NewActorRepository(actors []*models.Actor) ActorRepository {
	return &actorRepository{
		actors: actors,
	}
}

func (a actorRepository) GetTopActorsByPRAndCommits(x int) (actors []models.Actor) {
	if x > len(a.actors) {
		x = len(a.actors)
	}
	sort.Slice(a.actors, func(i, j int) bool {
		return (a.actors[i].PullRequests + a.actors[i].Commits) > (a.actors[j].PullRequests + a.actors[j].Commits)
	})

	for i := 0; i < x; i++ {
		actors = append(actors, *a.actors[i])
	}

	return
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/m7shapan/aggregated-analytics-task/aggregators"
	"github.com/m7shapan/aggregated-analytics-task/models"
	"github.com/m7shapan/aggregated-analytics-task/repositories"
	"github.com/m7shapan/aggregated-analytics-task/storage"
	"github.com/olekukonko/tablewriter"
)

func main() {
	flag.Parse()
	filePath := flag.Arg(0)

	aggregator := openFilesAndGetAggregator(filePath)

	actorsData, reposData, err := aggregator.AggregateAll()
	if err != nil {
		log.Fatal(err)
	}

	actorRepository := repositories.NewActorRepository(actorsData)
	repoRepository := repositories.NewRepoRepository(reposData)

	displayTopUsersByCommentsAndPRs(actorRepository.GetTopActorsByPRAndCommits(10))
	displayTopReposByComments(repoRepository.GetTopReposByCommits(10))
	displayTopReposByWatchNumbers(repoRepository.GetTopReposByWatch(10))
}

func openFilesAndGetAggregator(filePath string) aggregators.Aggregator {
	events, err := storage.NewFileReader(filePath, storage.EventsFile)
	if err != nil {
		log.Fatal(err)
	}

	commits, err := storage.NewFileReader(filePath, storage.CommitsFile)
	if err != nil {
		log.Fatal(err)
	}

	actors, err := storage.NewFileReader(filePath, storage.ActorsFile)
	if err != nil {
		log.Fatal(err)
	}

	repos, err := storage.NewFileReader(filePath, storage.ReposFile)
	if err != nil {
		log.Fatal(err)
	}

	return aggregators.NewAggregator(actors, repos, events, commits)
}

func displayTopUsersByCommentsAndPRs(top []models.Actor) {
	fmt.Printf("==========================================================================================\n")
	fmt.Printf("Top %d active users sorted by amount of PRs created and commits pushed\n", len(top))
	fmt.Printf("==========================================================================================\n")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"#", "User Name", "Commits", "Pull Requests", "Comments + PRs"})
	for i := 0; i < len(top); i++ {
		table.Append([]string{fmt.Sprint(i + 1), top[i].Name, fmt.Sprint(top[i].Commits), fmt.Sprint(top[i].PullRequests), fmt.Sprint(top[i].Commits + top[i].PullRequests)})
	}
	table.Render()
}

func displayTopReposByComments(top []models.RepoCommit) {
	fmt.Printf("==========================================================================================\n")
	fmt.Printf("Top %d repositories sorted by amount of commits pushed\n", len(top))
	fmt.Printf("==========================================================================================\n")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"#", "Repo Name", "Commits"})
	for i := 0; i < len(top); i++ {
		table.Append([]string{fmt.Sprint(i + 1), top[i].Name, fmt.Sprint(top[i].Commits)})
	}
	table.Render()
}

func displayTopReposByWatchNumbers(top []models.RepoWatch) {
	fmt.Printf("==========================================================================================\n")
	fmt.Printf("Top %d repositories sorted by amount of watch events\n", len(top))
	fmt.Printf("==========================================================================================\n")

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowLine(true)
	table.SetHeader([]string{"#", "Repo Name", "Watch numbers"})
	for i := 0; i < len(top); i++ {
		table.Append([]string{fmt.Sprint(i + 1), top[i].Name, fmt.Sprint(top[i].Watch)})
	}
	table.Render()
}

# Aggregated Analytics Task
is a CLI application that aggregates 1 hour of Github events and writes outputs to the console :
- Top 10 active users sorted by the amount of PRs created and commits pushed
- Top 10 repositories sorted by the amount of commits pushed
- Top 10 repositories sorted by the amount of watch events


## Installation

Install Aggregated Analytics Task with git

```bash
  git clone https://github.com/m7shapan/aggregated-analytics-task
  cd aggregated-analytics-task
```
then build it
```bash
    go build -o aggregator
```

Now it's ready to be used
    
## Usage Examples

```bash
    ./aggregator
```
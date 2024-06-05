package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
)

type Owner struct {
	Login string `json:"login"`
}

type License struct {
	Key    string `json:"key"`
	Name   string `json:"name"`
	SpdxID string `json:"spdx_id"`
	URL    string `json:"url"`
}

type Repository struct {
	ID              int         `json:"id"`
	NodeID          string      `json:"node_id"`
	Name            string      `json:"name"`
	FullName        string      `json:"full_name"`
	Owner           Owner       `json:"owner"`
	Parent          *Repository `json:"parent,omitempty"`
	IsFork          bool        `json:"fork"`
	URL             string      `json:"url"`
	ForksCount      int         `json:"forks_count"`
	StargazersCount int         `json:"stargazers_count"`
	WatchersCount   int         `json:"watchers_count"`
	OpenIssuesCount int         `json:"open_issues_count"`
	License         *License    `json:"license,omitempty"`
	DefaultBranch   string      `json:"default_branch"`
}

func fetchRepoInfo(repo string, token string, w *tabwriter.Writer) (*Repository, error) {
	parts := strings.Split(repo, "/")
	if len(parts) < 2 {
		return nil, fmt.Errorf("invalid repo format: %s", repo)
	}

	repo = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])

	url := fmt.Sprintf("https://api.github.com/repos/%s", repo)
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request for repo %s: %w", repo, err)
	}

	req.Header.Set("Authorization", "token "+token)
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request for repo %s: %w", repo, err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("received non-200 response from GitHub API for repo %s: %d", repo, resp.StatusCode)
	}

	defer resp.Body.Close()
	var repository Repository
	err = json.NewDecoder(resp.Body).Decode(&repository)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling JSON for repo %s: %w", repo, err)
	}

	licenseName := ""
	if repository.License != nil {
		licenseName = repository.License.Name
	}

	parentRepo := ""
	if repository.IsFork && repository.Parent != nil {
		parentRepo = repository.Parent.FullName
	}

	fmt.Fprintf(w, "%s\t%t\t%s\t%d\t%d\t%d\t%s\t%s\n",
		repository.FullName,
		repository.IsFork,
		parentRepo,
		repository.StargazersCount,
		repository.WatchersCount,
		repository.OpenIssuesCount,
		licenseName,
		repository.DefaultBranch,
	)

	return &repository, nil
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		fmt.Println("GitHub token not found, exiting...")
		os.Exit(1)
	}

	fmt.Println("GitHub token found, proceeding...")

	file, err := os.Open("go.mod")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile(`(?m)^\s*(\S+) v(\S+)`)
	fmt.Println("Scanning go.mod file...")

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Name\tIsFork\tParent Repo\tStargazers\tWatchers\tOpen Issues\tLicense\tDefault Branch")

	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		if match != nil {
			modulePath := match[1]
			parts := strings.Split(modulePath, "/")
			if len(parts) > 2 {
				repo := fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
				_, err := fetchRepoInfo(repo, token, w)
				if err != nil {
					fmt.Printf("Error fetching repo info: %v\n", err)
					continue
				}
			}
		}
	}

	w.Flush()

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

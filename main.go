package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
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
	fmt.Println("Scanning go.mod file...")

	w := writeHeader()

	var wg sync.WaitGroup
	repos := make(chan string)
	repoInfo := make(chan *Repository)

	wg.Add(1)
	go func() {
		defer wg.Done()
		collectRepos(scanner, repos)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for r := range fetchRepoInfo(repos, token) {
			repoInfo <- r
		}
		close(repoInfo)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		repoList := sortRepos(repoInfo)
		printRepos(repoList, w)
	}()

	wg.Wait()
	w.Flush()
}

func writeHeader() *tabwriter.Writer {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', tabwriter.AlignRight|tabwriter.Debug)
	fmt.Fprintln(w, "Name\tIsFork\tParent Repo\tStargazers\tWatchers\tOpen Issues\tLicense\tDefault Branch")
	return w
}

func collectRepos(scanner *bufio.Scanner, repos chan<- string) {
	go func() {
		defer close(repos)

		re := regexp.MustCompile(`(?m)^\s*(\S+) v(\S+)`)
		for scanner.Scan() {
			line := scanner.Text()
			match := re.FindStringSubmatch(line)
			if match != nil {
				modulePath := match[1]
				parts := strings.Split(modulePath, "/")
				if len(parts) > 2 {
					repo := fmt.Sprintf("%s/%s/%s", parts[0], parts[1], parts[2])
					repos <- repo
				}
			}
		}
	}()
}

func fetchRepoInfo(repos <-chan string, token string) <-chan *Repository {
	out := make(chan *Repository)
	go func() {
		defer close(out)
		for repo := range repos {
			parts := strings.Split(repo, "/")
			if len(parts) < 2 {
				log.Printf("invalid repo format: %s", repo)
				continue
			}

			repo = fmt.Sprintf("%s/%s", parts[len(parts)-2], parts[len(parts)-1])

			url := fmt.Sprintf("https://api.github.com/repos/%s", repo)
			client := &http.Client{}
			req, err := http.NewRequest("GET", url, nil)
			if err != nil {
				log.Printf("error creating request for repo %s: %v", repo, err)
				continue
			}

			req.Header.Set("Authorization", "token "+token)
			resp, err := client.Do(req)
			if err != nil {
				log.Printf("error making request for repo %s: %v", repo, err)
				continue
			}

			if resp.StatusCode != 200 {
				log.Printf("received non-200 response from GitHub API for repo %s: %d", repo, resp.StatusCode)
				continue
			}

			defer resp.Body.Close()
			var repository Repository
			err = json.NewDecoder(resp.Body).Decode(&repository)
			if err != nil {
				log.Printf("error unmarshalling JSON for repo %s: %v", repo, err)
				continue
			}

			out <- &repository
		}
	}()
	return out
}

func sortRepos(repos <-chan *Repository) []Repository {
	var repoList []Repository
	for repo := range repos {
		repoList = append(repoList, *repo)
	}

	sort.Slice(repoList, func(i, j int) bool {
		return repoList[i].WatchersCount < repoList[j].WatchersCount
	})

	return repoList
}

func printRepos(repos []Repository, w *tabwriter.Writer) {
	for _, repo := range repos {
		licenseName := ""
		if repo.License != nil {
			licenseName = repo.License.Name
		}

		parentRepo := ""
		if repo.IsFork && repo.Parent != nil {
			parentRepo = repo.Parent.FullName
		}

		fmt.Fprintf(w, "%s\t%t\t%s\t%d\t%d\t%d\t%s\t%s\n",
			repo.FullName,
			repo.IsFork,
			parentRepo,
			repo.StargazersCount,
			repo.WatchersCount,
			repo.OpenIssuesCount,
			licenseName,
			repo.DefaultBranch,
		)
	}
}

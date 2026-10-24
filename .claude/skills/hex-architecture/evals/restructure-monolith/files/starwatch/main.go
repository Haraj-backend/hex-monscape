package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"time"
)

type Repo struct {
	FullName   string    `json:"full_name"`
	Stars      int       `json:"stargazers_count"`
	OpenIssues int       `json:"open_issues_count"`
	PushedAt   time.Time `json:"pushed_at"`
	Score      float64   `json:"score"`
}

func main() {
	token := os.Getenv("GITHUB_TOKEN")
	reposFile := os.Getenv("REPOS_FILE")
	if reposFile == "" {
		reposFile = "repos.txt"
	}
	outFile := os.Getenv("OUT_FILE")
	if outFile == "" {
		outFile = "report.json"
	}

	f, err := os.Open(reposFile)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	var names []string
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		names = append(names, line)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	var repos []Repo
	for _, name := range names {
		req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+name, nil)
		if token != "" {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("skip %s: %v", name, err)
			continue
		}
		if resp.StatusCode == 404 {
			resp.Body.Close()
			log.Printf("skip %s: not found", name)
			continue
		}
		var r Repo
		err = json.NewDecoder(resp.Body).Decode(&r)
		resp.Body.Close()
		if err != nil {
			log.Printf("skip %s: %v", name, err)
			continue
		}
		// trending score: popular, active, and with work to do
		days := time.Since(r.PushedAt).Hours() / 24
		activity := 1.0
		if days > 30 {
			activity = 0
		} else if days > 7 {
			activity = 0.5
		}
		r.Score = float64(r.Stars)/1000*0.6 + activity*30 + float64(min(r.OpenIssues, 500))/500*10
		repos = append(repos, r)
	}

	sort.Slice(repos, func(i, j int) bool { return repos[i].Score > repos[j].Score })
	data, err := json.MarshalIndent(map[string]any{"generated_at": time.Now().UTC(), "repos": repos}, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	if err := os.WriteFile(outFile, data, 0o644); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("wrote %d repos to %s\n", len(repos), outFile)
}

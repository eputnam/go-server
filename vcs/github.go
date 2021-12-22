package vcs

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/eputnam/health-check-server/config"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

const questions_path = "questions.json"

type GitHubBoi struct {
	Client  *github.Client
	Context context.Context
	Config  config.GitHubConfig
}

type questionsResponse struct {
	Questions []string
}

func NewClient(config config.GitHubConfig) GitHubBoi {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.Token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return GitHubBoi{Config: config, Client: github.NewClient(tc), Context: ctx}
}

func (g *GitHubBoi) GetQuestionsList() (questionsResponse, error) {
	byteData, err := g.getGitHubBytes(questions_path)
	if nil != err {
		return questionsResponse{Questions: nil}, err
	}

	var jsonData questionsResponse
	json.Unmarshal(byteData, &jsonData)

	return jsonData, nil
}

func (g *GitHubBoi) getGitHubBytes(path string) ([]byte, error) {
	content, _, _, err := g.Client.Repositories.GetContents(g.Context, g.Config.Username, g.Config.Repository, path, nil)
	checkError(err)

	byteData, err := base64.StdEncoding.DecodeString(*content.Content)
	checkError(err)

	return byteData, err
}

func checkError(err error) {
	if nil != err {
		panic(err)
	}
}

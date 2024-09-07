package context

import (
	"github.com/chancehl/godo/internal/services"
)

type AppContext struct {
	GistID        string
	GodoService   *services.GodoService
	GithubService *services.GithubService
}

func (ctx *AppContext) NewAppContext(gistID string, godoService *services.GodoService, githubService *services.GithubService) *AppContext {
	return &AppContext{
		GistID:        gistID,
		GodoService:   godoService,
		GithubService: githubService,
	}
}

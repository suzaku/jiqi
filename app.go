package main

import (
	"changeme/k8s"
	"context"
)

// App struct
type App struct {
	ctx          context.Context
	nodesManager k8s.NodesManager
}

// NewApp creates a new App application struct
func NewApp() *App {
	nm, err := k8s.NewNodesManager()
	if err != nil {
		panic(err) // FIXME
	}
	return &App{
		nodesManager: nm,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListNodes(shouldClearCache bool) []k8s.Node {
	nodes, err := a.nodesManager.GetNodes(a.ctx, shouldClearCache)
	if err != nil {
		// TODO: Show error message
		panic(err)
	}
	return nodes
}

func (a *App) GetCurrentContext() string {
	currentContext, err := a.nodesManager.GetCurrentContext()
	if err != nil {
		panic(err)
	}
	return currentContext
}

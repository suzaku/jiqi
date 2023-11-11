package main

import (
	"changeme/k8s"
	"context"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) ListNodes() []k8s.Node {
	nodes, err := k8s.GetNodes()
	if err != nil {
		// TODO: Show error message
		panic(err)
	}
	return nodes
}

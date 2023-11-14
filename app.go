package main

import (
	"changeme/k8s"
	"context"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"strings"
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

func (a *App) ListNodes(shouldClearCache bool, labelSelectors string) []k8s.Node {
	nodes, err := a.nodesManager.GetNodes(a.ctx, shouldClearCache)
	if err != nil {
		// TODO: Show error message
		panic(err)
	}
	labelSelectors = strings.TrimSpace(labelSelectors)
	if len(labelSelectors) == 0 {
		return nodes
	}
	var selectors [][2]string
	for _, selector := range strings.Split(labelSelectors, ",") {
		parts := strings.SplitN(selector, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key, value := parts[0], parts[1]
		selectors = append(selectors, [2]string{key, value})
	}
	var selected []k8s.Node
	for _, n := range nodes {
		for _, selector := range selectors {
			k, v := selector[0], selector[1]
			runtime.LogPrintf(a.ctx, "Selector: %v", selector)
			runtime.LogPrintf(a.ctx, "Labels: %v", n.Labels)
			if n.Labels[k] == v {
				selected = append(selected, n)
				break
			}
		}
	}
	return selected
}

func (a *App) GetCurrentContext() string {
	currentContext, err := a.nodesManager.GetCurrentContext()
	if err != nil {
		panic(err)
	}
	return currentContext
}

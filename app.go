package main

import (
	"changeme/k8s"
	"context"
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

type QueriedNodes struct {
	Nodes  []k8s.Node          `json:"nodes"`
	Labels map[string][]string `json:"labels"`
}

func indexLabels(nodes []k8s.Node) map[string][]string {
	labels := make(map[string][]string)
	seen := make(map[string]map[string]struct{})
	for _, n := range nodes {
		for k, v := range n.Labels {
			if _, ok := seen[k]; !ok {
				seen[k] = make(map[string]struct{})
			}
			if _, ok := seen[k][v]; ok {
				continue
			} else {
				seen[k][v] = struct{}{}
				labels[k] = append(labels[k], v)
			}
		}
	}
	return labels
}

func (a *App) ListNodes(shouldClearCache bool, labelSelectors string) QueriedNodes {
	nodes, err := a.nodesManager.GetNodes(a.ctx, shouldClearCache)
	if err != nil {
		// TODO: Show error message
		panic(err)
	}
	labelsIndexForAllNodes := indexLabels(nodes)
	labelSelectors = strings.TrimSpace(labelSelectors)
	if len(labelSelectors) == 0 {
		return QueriedNodes{
			Nodes:  nodes,
			Labels: labelsIndexForAllNodes,
		}
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
			if n.Labels[k] == v {
				selected = append(selected, n)
				break
			}
		}
	}
	return QueriedNodes{
		Nodes:  selected,
		Labels: labelsIndexForAllNodes,
	}
}

func (a *App) GetCurrentContext() string {
	currentContext, err := a.nodesManager.GetCurrentContext()
	if err != nil {
		panic(err)
	}
	return currentContext
}

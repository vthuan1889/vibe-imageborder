package main

import (
	"context"
)

// App struct holds the application state.
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct.
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved for runtime methods.
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

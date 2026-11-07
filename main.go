package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

const filename = "tasks.json"

type Task struct {
	Name    string `json:"Name"`
	Checked bool   `json:"Checked"`
	Version int    `json:"Version"`
}

func saveTasks(path string, tasks []Task) error {
	if tasks == nil {
		tasks = []Task{}
	}
	data, err := json.MarshalIndent(tasks, "", "\t")
	if err != nil {
		return fmt.Errorf("encode tasks: %w", err)
	}

	file, err := os.CreateTemp(filepath.Dir(path), ".tasks-*.json")
	if err != nil {
		return fmt.Errorf("create temporary task file: %w", err)
	}
	defer os.Remove(file.Name())

	if _, err := file.Write(data); err != nil {
		file.Close()
		return fmt.Errorf("write tasks: %w", err)
	}
	if err := file.Close(); err != nil {
		return fmt.Errorf("close task file: %w", err)
	}
	if err := os.Rename(file.Name(), path); err != nil {
		return fmt.Errorf("replace task file: %w", err)
	}
	return nil
}

func loadTasks(path string) ([]Task, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return []Task{}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("read tasks: %w", err)
	}
	if strings.TrimSpace(string(data)) == "" {
		return []Task{}, nil
	}

	var tasks []Task
	if err := json.Unmarshal(data, &tasks); err != nil {
		return nil, fmt.Errorf("decode tasks: %w", err)
	}
	return tasks, nil
}

func main() {
	application := app.New()
	window := application.NewWindow("Task Manager")

	tasks, err := loadTasks(filename)
	if err != nil {
		window.SetContent(widget.NewLabel("Unable to load tasks. Check tasks.json and restart."))
		dialog.ShowError(err, window)
	} else {
		manager := &taskManager{
			window:  window,
			path:    filename,
			tasks:   tasks,
			active:  container.NewVBox(),
			archive: container.NewVBox(),
		}
		window.SetContent(manager.content())
	}

	window.Resize(fyne.NewSize(420, 480))
	window.ShowAndRun()
}

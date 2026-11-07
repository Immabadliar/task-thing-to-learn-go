package main

import (
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

type taskManager struct {
	window  fyne.Window
	path    string
	tasks   []Task
	active  *fyne.Container
	archive *fyne.Container
}

func (manager *taskManager) content() fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter task name...")

	add := func() {
		name := strings.TrimSpace(input.Text)
		if name == "" {
			return
		}
		tasks := append(append([]Task{}, manager.tasks...), Task{Name: name})
		if manager.save(tasks) {
			input.SetText("")
		}
	}
	input.OnSubmitted = func(string) { add() }
	manager.refresh()

	return container.NewBorder(
		container.NewVBox(input, widget.NewButton("Add Task", add)),
		nil, nil, nil,
		container.NewVScroll(container.NewVBox(
			widget.NewLabel("Active"),
			manager.active,
			widget.NewLabel("Archive"),
			manager.archive,
		)),
	)
}

func (manager *taskManager) save(tasks []Task) bool {
	if err := saveTasks(manager.path, tasks); err != nil {
		dialog.ShowError(err, manager.window)
		return false
	}
	manager.tasks = tasks
	manager.refresh()
	return true
}

func (manager *taskManager) refresh() {
	manager.active.RemoveAll()
	manager.archive.RemoveAll()

	for index, task := range manager.tasks {
		checkbox := widget.NewCheck(task.Name, nil)
		checkbox.SetChecked(task.Checked)
		checkbox.OnChanged = func(checked bool) {
			tasks := append([]Task{}, manager.tasks...)
			tasks[index].Checked = checked
			if !manager.save(tasks) {
				manager.refresh()
			}
		}
		if task.Checked {
			manager.archive.Add(checkbox)
		} else {
			manager.active.Add(checkbox)
		}
	}
}

package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func TestTaskPersistence(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tasks.json")
	tasks, err := loadTasks(path)
	if err != nil || len(tasks) != 0 {
		t.Fatalf("load missing file: tasks=%v, err=%v", tasks, err)
	}
	want := []Task{{Name: "First"}, {Name: "Second", Checked: true, Version: 2}}
	if err := saveTasks(path, want); err != nil {
		t.Fatal(err)
	}
	got, err := loadTasks(path)
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("round trip: got=%v, err=%v", got, err)
	}
	if err := os.WriteFile(path, []byte("invalid"), 0600); err != nil {
		t.Fatal(err)
	}
	if _, err := loadTasks(path); err == nil {
		t.Fatal("expected malformed JSON to fail")
	}
}

func TestTaskCompletion(t *testing.T) {
	application := test.NewApp()
	defer application.Quit()
	manager := &taskManager{
		window:  application.NewWindow("Test"),
		path:    filepath.Join(t.TempDir(), "tasks.json"),
		tasks:   []Task{{Name: "Duplicate"}, {Name: "Duplicate"}},
		active:  container.NewVBox(),
		archive: container.NewVBox(),
	}
	manager.refresh()
	manager.active.Objects[1].(*widget.Check).SetChecked(true)
	if manager.tasks[0].Checked || !manager.tasks[1].Checked {
		t.Fatalf("wrong task completed: %v", manager.tasks)
	}
	if len(manager.active.Objects) != 1 || len(manager.archive.Objects) != 1 {
		t.Fatal("completed task did not move to archive")
	}
	manager.archive.Objects[0].(*widget.Check).SetChecked(false)
	if len(manager.active.Objects) != 2 || len(manager.archive.Objects) != 0 {
		t.Fatal("unchecked task did not return to active")
	}
	saved, err := loadTasks(manager.path)
	if err != nil || !reflect.DeepEqual(saved, manager.tasks) {
		t.Fatalf("completion was not saved: tasks=%v, err=%v", saved, err)
	}
}

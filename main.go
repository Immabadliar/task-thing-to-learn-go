package main

import (
	"fmt"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2"
)


func main() {
	appl := app.New()
	wind := appl.NewWindow("my name is terry davis")

	taskInput := widget.NewEntry()
	taskInput.SetPlaceHolder("Enter task name...")

	activeText := widget.NewLabel("active")
	archiveText := widget.NewLabel("archive")

	activeContainer := container.NewVBox()

	archiveContainer := container.NewVBox()

	addBtn := widget.NewButton("Add Task", func() {
		currentText := taskInput.Text

		if currentText != "" {
			AddTask(wind, currentText, activeContainer, archiveContainer)
			taskInput.SetText("")
		}
	})

	mainLayout := container.NewVBox(
		taskInput,
		addBtn,
		activeText,
		activeContainer,
		archiveText,
		archiveContainer,
	)

	wind.SetContent(mainLayout)

	wind.Resize(fyne.NewSize(200, 200))
	wind.ShowAndRun()
	appl.Run()
	tidyUp()
}

func tidyUp() {
	fmt.Println("Exited")
}

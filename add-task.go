package main

import (
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2"
)

func AddTask(win fyne.Window, name string, active *fyne.Container, archive *fyne.Container) {
	var newCheckbox *widget.Check

	newCheckbox = widget.NewCheck(name, func(checked bool) {
		if newCheckbox.Checked == true {
			active.Remove(newCheckbox)
			archive.Add(newCheckbox)
			active.Refresh()
			archive.Refresh()
		}

		for _, child := range archive.Objects {
			if child == newCheckbox {
				if newCheckbox.Checked != true {
					archive.Remove(newCheckbox)
					active.Add(newCheckbox)
					archive.Refresh()
					active.Refresh()
				}
			}
		}
	})

	active.Add(newCheckbox)
	active.Refresh()
	
}
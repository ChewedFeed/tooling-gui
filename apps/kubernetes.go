package apps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func KubernetesView(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle(
				"Kubernetes View",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true})))
}

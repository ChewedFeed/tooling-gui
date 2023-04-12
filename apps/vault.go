package apps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func VaultTokenView(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle(
				"Vault Token",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true})))
}

func VaultDatabaseTokenView(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle(
				"Vault Database Token",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true})))
}

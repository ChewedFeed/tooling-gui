package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/chewedfeed/automated/apps"
)

const currentApp = "unknownApp"

func main() {
	a := app.NewWithID("net.chewed-k8s.automated")
	a.Settings().SetTheme(theme.DarkTheme())
	logWindow(a)

	w := a.NewWindow("Automated ChewedFeed")
	w.SetMaster()
	w.SetMainMenu(makeMenu(a))

	content := container.NewMax()
	title := widget.NewLabel("Automated ChewedFeed")
	title.Alignment = fyne.TextAlignCenter
	setApp := func(app apps.App) {
		title.SetText(app.AppTitle)
		content.Objects = []fyne.CanvasObject{app.View(w)}
		content.Refresh()
	}
	appItem := container.NewBorder(container.NewVBox(title, widget.NewSeparator()), nil, nil, nil, content)
	split := container.NewHSplit(makeNav(setApp), appItem)
	split.Offset = 0.2
	w.SetContent(split)
	w.Resize(fyne.NewSize(800, 600))

	w.ShowAndRun()
}

func logWindow(a fyne.App) {
	a.Lifecycle().SetOnStarted(func() {
		logs.Local().Info("Automated ChewedFeed started")
	})
	a.Lifecycle().SetOnStopped(func() {
		logs.Local().Info("Automated ChewedFeed stopped")
	})
}

func makeMenu(a fyne.App) *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu("File",
			fyne.NewMenuItem("Quit", func() {
				a.Quit()
			}),
		),
	)
}

func makeNav(setApp func(app apps.App)) fyne.CanvasObject {
	ca := fyne.CurrentApp()
	tree := &widget.Tree{
		ChildUIDs: func(uid string) []string {
			return apps.AppsIndex[uid]
		},
		IsBranch: func(uid string) bool {
			children, ok := apps.AppsIndex[uid]
			return ok && len(children) > 0
		},
		CreateNode: func(branch bool) fyne.CanvasObject {
			return widget.NewLabel("tester")
		},
		UpdateNode: func(uid string, branch bool, node fyne.CanvasObject) {
			a, ok := apps.Apps[uid]
			if !ok {
				fyne.LogError("Unknown app: "+uid, nil)
				return
			}
			node.(*widget.Label).SetText(a.MenuTitle)
			node.(*widget.Label).TextStyle = fyne.TextStyle{}
		},
		OnSelected: func(uid string) {
			if a, ok := apps.Apps[uid]; ok {
				ca.Preferences().SetString(currentApp, uid)
				setApp(a)
			}
		},
	}

	currentPref := ca.Preferences().StringWithFallback(currentApp, "Tester")
	tree.Select(currentPref)
	return container.NewBorder(nil, nil, nil, nil, tree)
}

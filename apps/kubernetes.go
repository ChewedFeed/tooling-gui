package apps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/chewedfeed/automated/internal/kubernetes"
)

type minWidthEntry struct {
	widget.Entry
	minWidth float32
}

type doubleEntry struct {
	key   *minWidthEntry
	value *minWidthEntry
}

func newMin(min float32) *minWidthEntry {
	e := &minWidthEntry{
		minWidth: min,
	}
	e.ExtendBaseWidget(e)
	return e
}
func (e *minWidthEntry) MinSize() fyne.Size {
	m := e.Entry.MinSize()
	if m.Width < e.minWidth {
		m.Width = e.minWidth
	}
	return m
}

func KubernetesView(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle(
				"Kubernetes View",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true})))
}

func KubernetesSecretView(w fyne.Window) fyne.CanvasObject {
	type secretDetails struct {
		Name      string
		Namespace string
		Data      map[string]string
	}
	s := &secretDetails{}

	k, err := kubernetes.NewKubernetes()
	if err != nil {
		dialog.ShowError(logs.Local().Errorf("failed to get kubernetes: %+v", err), w)
	}

	ns, err := k.GetNamespaces()
	if err != nil {
		dialog.ShowError(logs.Local().Errorf("failed to get kubernetes namespaces: %+v", err), w)
	}

	secretName := widget.NewEntry()
	namespaceSelect := widget.NewSelect(ns, func(s string) {})
	formItems := []*widget.FormItem{
		{Text: "Secret Name", Widget: secretName},
		{Text: "Namespace", Widget: namespaceSelect},
	}
	dataEntries := []*doubleEntry{}

	form := &widget.Form{
		Items: formItems,
		OnSubmit: func() {
			s.Name = secretName.Text
			s.Namespace = namespaceSelect.Selected

			data := make(map[string]string)
			for _, d := range dataEntries {
				if d.key.Entry.Text == "" || d.value.Entry.Text == "" {
					continue
				}
				data[d.key.Entry.Text] = d.value.Entry.Text
			}
			s.Data = data

			if err := k.CreateSecret(s.Name, s.Namespace, s.Data); err != nil {
				dialog.ShowError(logs.Local().Errorf("failed to create kubernetes secret: %+v", err), w)
			}
			dialog.ShowInformation("Success", "Secret created", w)

			// logs.Local().Infof("Submitted form with robot name: %s, project: %s, create: %b", robotName.Text, projectSelect.Selected, createRobot.Checked)
		},
		OnCancel: func() {
			logs.Local().Infof("Cancelled form")
		},
	}

	addEntry := func() {
		k := newMin(200)
		v := newMin(200)

		d := container.NewHBox(
			widget.NewLabel("Key"),
			k,
			layout.NewSpacer(),
			widget.NewLabel("Value"),
			v)
		form.Append("", d)
		dataEntries = append(dataEntries, &doubleEntry{key: k, value: v})
		d.Resize(fyne.NewSize(500, 100))
	}
	addButton := widget.NewButton("Add", addEntry)

	return container.NewVBox(form, addButton)
}

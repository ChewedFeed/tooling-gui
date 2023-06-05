package apps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/chewedfeed/automated/internal/vault"
)

func VaultTokenView(_ fyne.Window) fyne.CanvasObject {
	return container.NewCenter(
		container.NewVBox(
			widget.NewLabelWithStyle(
				"Vault Token",
				fyne.TextAlignCenter,
				fyne.TextStyle{Bold: true})))
}

func VaultDatabaseTokenView(w fyne.Window) fyne.CanvasObject {
	v, err := vault.NewVault()
	if err != nil {
		fyne.LogError("Failed to get vault", err)
	}
	creds, err := v.GetCreds()
	if err != nil {
		fyne.LogError("Failed to get vault creds", err)
	}
	credsSelect := widget.NewSelect(creds, nil)
	readCapability := widget.NewCheck("Read", nil)
	writeCapability := widget.NewCheck("Write", nil)
	policyName := widget.NewEntry()
	policyName.SetPlaceHolder("Bobs Policy")

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Policy Name", Widget: policyName},
			{Text: "Database Connection Name", Widget: credsSelect},
			{Text: "Read", Widget: readCapability},
			{Text: "Write", Widget: writeCapability},
		},
		OnSubmit: func() {
			if policyName.Text == "" {
				policyName.Text = credsSelect.Selected + "-automated-policy"
			}

			token, err := v.CreateDatabaseToken(credsSelect.Selected, policyName.Text, readCapability.Checked, writeCapability.Checked)
			if err != nil {
				fyne.LogError("Failed to create database token", err)
			}
			dialog.ShowInformation("Database Token", token, w)

			// reset form
			credsSelect.Selected = ""
			readCapability.Checked = false
			writeCapability.Checked = false
			policyName.Text = ""
		},
		OnCancel: func() {
			credsSelect.Selected = ""
			readCapability.Checked = false
			writeCapability.Checked = false
			policyName.Text = ""

			logs.Local().Infof("Cancelled form")
		},
	}

	return container.NewVBox(form)
}

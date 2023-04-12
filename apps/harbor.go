package apps

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
	"github.com/bugfixes/go-bugfixes/logs"
	"github.com/chewedfeed/automated/internal/harbor"
	"time"
)

func getHarborProjects() ([]string, error) {
	h, err := harbor.NewHarbor()
	if err != nil {
		return nil, err
	}
	projects, err := h.GetProjects()
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func HarborRegistrySecretView(w fyne.Window) fyne.CanvasObject {
	robotName := widget.NewEntry()
	robotName.SetPlaceHolder("Robot Name")

	harborProjects, err := getHarborProjects()
	if err != nil {
		fyne.LogError("Failed to get harbor projects", err)
	}
	projectSelect := widget.NewSelect(harborProjects, nil)
	robotSecret := widget.NewEntry()
	robotSecret.SetPlaceHolder("Robot Secret")
	var expireDate time.Time

	robotCreateDetails := harbor.RobotDetails{}
	createRobot := widget.NewButton("Create Robot", func() {
		robotType := widget.NewRadioGroup([]string{"System", "Project"}, nil)
		l := widget.NewLabel("None Expiring")
		l.Alignment = fyne.TextAlignLeading
		i := widget.NewLabel("Select Date")
		i.Alignment = fyne.TextAlignLeading
		datePicker := xwidget.NewCalendar(time.Now(), func(date time.Time) {
			l.SetText(date.Format("2006-01-02"))
			expireDate = date
		})

		formItems := []*widget.FormItem{
			{Text: "Robot Type", Widget: robotType},
			{Text: "Selected Date", Widget: l},
			{Text: "Expire Date", Widget: datePicker},
		}

		robotSecret.Disabled()
		dialog.ShowForm("Create Robot", "Create Robot", "Cancel", formItems, func(submitted bool) {
			if !submitted {
				return
			}
			systemRobot := false
			if robotType.Selected == "System" {
				systemRobot = true
			}

			if robotName.Text == "" {
				dialog.ShowError(logs.Local().Errorf("Robot Name Needs to be input"), w)
				return
			}
			if projectSelect.Selected == "" {
				dialog.ShowError(logs.Local().Errorf("Project Needs to be selected"), w)
				return
			}

			if robotType.Selected == "" {
				dialog.ShowError(logs.Local().Errorf("Robot Type Needs to be selected"), w)
				return
			}

			h, err := harbor.NewHarbor()
			if err != nil {
				fyne.LogError("Failed to create harbor client", err)
			}
			r, err := h.CreateRobot(robotName.Text, projectSelect.Selected, systemRobot, expireDate)
			if err != nil {
				fyne.LogError("Failed to create robot", err)
			}
			robotCreateDetails = *r
		}, w)
	})

	form := &widget.Form{
		Items: []*widget.FormItem{
			{Text: "Robot Name", Widget: robotName},
			{Text: "Project", Widget: projectSelect},
			{Text: "Create Robot", Widget: createRobot},
			{Text: "Robot Secret", Widget: robotSecret},
		},
		OnSubmit: func() {
			robotValid := false
			if robotCreateDetails.Name != "" {
				robotValid = true
			}
			robotDetails := harbor.RobotDetails{}

			if projectSelect.Selected == "" {
				fyne.CurrentApp().SendNotification(&fyne.Notification{
					Title:   "Error",
					Content: "Please select a project",
				})
			}

			h, err := harbor.NewHarbor()
			if err != nil {
				fyne.LogError("Failed to create harbor client", err)
			}

			valid, err := h.ValidateRobot(robotName.Text, robotSecret.Text)
			if err != nil {
				fyne.LogError("Failed to validate robot", err)
			}
			if !valid {
				fyne.LogError("Invalid robot", err)
			}
			robotValid = valid
			robotDetails = harbor.RobotDetails{
				Name:   robotName.Text,
				Secret: robotSecret.Text,
			}

			if robotValid {
				logs.Local().Infof("robot details: %+v", robotDetails)
			}

			// logs.Local().Infof("Submitted form with robot name: %s, project: %s, create: %b", robotName.Text, projectSelect.Selected, createRobot.Checked)
		},
		OnCancel: func() {
			logs.Local().Infof("Cancelled form")
		},
	}

	return container.NewVBox(form)
}

package apps

import (
	"fyne.io/fyne/v2"
)

type App struct {
	MenuTitle string
	AppTitle  string
	View      func(w fyne.Window) fyne.CanvasObject
}

var (
	Apps = map[string]App{
		"welcome": {
			MenuTitle: "Welcome",
			AppTitle:  "Welcome",
			View:      WelcomeView,
		},
		"app": {
			MenuTitle: "App",
			AppTitle:  "App",
			View:      AppView,
		},
		"kubernetes": {
			MenuTitle: "Kubernetes",
			AppTitle:  "Kubernetes",
			View:      KubernetesView,
		},
		"harbor_registry_secret": {
			MenuTitle: "Harbor Registry Secret",
			AppTitle:  "Harbor Registry Secret",
			View:      HarborRegistrySecretView,
		},
		"vault_token": {
			MenuTitle: "Vault Token",
			AppTitle:  "Vault Token",
			View:      VaultTokenView,
		},
		"vault_database_token": {
			MenuTitle: "Vault Database Token",
			AppTitle:  "Vault Database Token",
			View:      VaultDatabaseTokenView,
		},
		"kubernetes_secret": {
			MenuTitle: "Kubernetes Secret",
			AppTitle:  "Kubernetes Secret",
			View:      KubernetesSecretView,
		},
	}
	AppsIndex = map[string][]string{
		"": {
			"welcome",
			"app",
			"kubernetes",
		},
		"app": {
			"vault_database_token",
		},
		"kubernetes": {
			"harbor_registry_secret",
			"vault_token",
			"kubernetes_secret",
		},
	}
)

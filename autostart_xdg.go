// +build !windows,!darwin

package autostart

import (
	"os"
	"path/filepath"
	"text/template"
)

const desktopTemplate = `[Desktop Entry]
Type=Application
Name={{.DisplayName}}
Exec={{.Exec}}
{{- if .Icon}}Icon={{.Icon}}{{end}}
X-GNOME-Autostart-enabled=true
`

var autostartDir string

func init() {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		autostartDir = os.Getenv("XDG_CONFIG_HOME")
	} else {
		autostartDir = filepath.Join(os.Getenv("HOME"), ".config")
	}
	autostartDir = filepath.Join(autostartDir, "autostart")
}

func (a *App) path() string {
	return filepath.Join(autostartDir, a.Name + ".desktop")
}

// Check if the app is enabled on startup.
func (a *App) IsEnabled() bool {
	_, err := os.Stat(a.path())
	return !os.IsNotExist(err)
}

type app struct {
	*App
}

// Override App.Exec to return a string.
func (a *app) Exec() string {
	return quote(a.App.Exec)
}

// Enable this app on startup.
func (a *App) Enable() error {
	t := template.Must(template.New("desktop").Parse(desktopTemplate))

	f, err := os.Create(a.path())
	if err != nil {
		return err
	}
	defer f.Close()

	return t.Execute(f, &app{a})
}

// Disable this app on startup.
func (a *App) Disable() error {
	return os.Remove(a.path())
}

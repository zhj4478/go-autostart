package autostart

import (
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
)

const jobTemplate = `<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
  <dict>
    <key>Label</key>
    <string>{{.Name}}</string>
    <key>ProgramArguments</key>
      <array>
        {{range .Exec -}}
        <string>{{.}}</string>
        {{end}}
      </array>
    <key>RunAtLoad</key>
    <true/>
  </dict>
</plist>`

var launchDir string

func init() {
	launchDir = filepath.Join(os.Getenv("HOME"), "Library", "LaunchAgents")
}

func (a *App) path() string {
	return filepath.Join(launchDir, a.Name+".plist")
}

// IsEnabled Check is app enabled startup.
func (a *App) IsEnabled() bool {
	_, err := os.Stat(a.path())
	return !os.IsNotExist(err)
}

// Enable this app on startup.
func (a *App) Enable() error {
	t := template.Must(template.New("job").Parse(jobTemplate))

	f, err := os.Create(a.path())
	if err != nil {
		return err
	}
	defer f.Close()

	if err := t.Execute(f, a); err != nil {
		return err
	}

	cmd := exec.Command("launchctl", "load", a.path())
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}

// Disable this app on startup.
func (a *App) Disable() error {
	cmd := exec.Command("launchctl", "unload", a.path())
	if err := cmd.Run(); err != nil {
		return err
	}

	return os.Remove(a.path())
}

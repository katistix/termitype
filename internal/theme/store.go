package theme

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

// settings is the on-disk config file. Only the theme is stored for now.
type settings struct {
	Theme string `json:"theme"`
}

// ConfigPath returns $XDG_CONFIG_HOME/termitype/config.json,
// falling back to ~/.config/termitype/config.json.
func ConfigPath() (string, error) {
	dir := os.Getenv("XDG_CONFIG_HOME")
	if dir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		dir = filepath.Join(home, ".config")
	}
	return filepath.Join(dir, "termitype", "config.json"), nil
}

// Load returns the saved theme, or Default if none is saved or it is unknown.
func Load() Theme {
	path, err := ConfigPath()
	if err != nil {
		return MustGet(Default)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return MustGet(Default)
	}
	var s settings
	if json.Unmarshal(data, &s) != nil {
		return MustGet(Default)
	}
	return MustGet(s.Theme)
}

// Save persists name as the selected theme.
func Save(name string) error {
	t, ok := Get(name)
	if !ok {
		return errors.New("unknown theme " + name)
	}
	path, err := ConfigPath()
	if err != nil {
		return err
	}
	// Preserve any other keys a future version may add.
	raw := map[string]any{}
	if data, err := os.ReadFile(path); err == nil {
		_ = json.Unmarshal(data, &raw)
	} else if !errors.Is(err, fs.ErrNotExist) {
		return err
	}
	raw["theme"] = t.Name
	data, err := json.MarshalIndent(raw, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

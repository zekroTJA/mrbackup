package modrinth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zekrotja/rogu/log"
)

type Modrinth struct {
	installDir string
	db         *AppDb
}

func New(installDir string) (*Modrinth, error) {
	installDir, err := findModrinthDir(installDir)
	if err != nil {
		return nil, err
	}

	log.Debug().Field("installDir", installDir).Msg("found Modrinth install dir")

	db, err := NewAppDb(filepath.Join(installDir, "app.db"))
	if err != nil {
		return nil, fmt.Errorf("failed opening app db: %w", err)
	}

	t := &Modrinth{
		installDir: installDir,
		db:         db,
	}

	return t, nil
}

func (t *Modrinth) Close() error {
	return t.db.Close()
}

func (t *Modrinth) Profiles() ([]*Profile, error) {
	profiles, err := t.db.Profiles()
	if err != nil {
		return nil, fmt.Errorf("failed getting profiles from db: %w", err)
	}

	for _, p := range profiles {
		p.FullPath = filepath.Join(t.installDir, "profiles", p.Path)
	}

	return profiles, nil
}

func findModrinthDir(path string) (string, error) {
	if path == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return "", fmt.Errorf("failed getting user config dir: %w", err)
		}
		path = filepath.Join(configDir, "ModrinthApp")
	}

	stat, err := os.Stat(path)
	if err == nil {
		if !stat.IsDir() {
			return "", fmt.Errorf("Modrinth directory is not a directory (%s)", path)
		}
		return path, nil
	}
	return "", fmt.Errorf("Modrinth directory error: %w", err)
}

package modrinth

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zekrotja/rogu/log"
)

type Modrinth struct {
	installDir string
}

func New(installDir string) (*Modrinth, error) {
	installDir, err := findModrinthDir(installDir)
	if err != nil {
		return nil, err
	}

	log.Debug().Field("installDir", installDir).Msg("found Modrinth install dir")

	t := &Modrinth{
		installDir: installDir,
	}

	return t, nil
}

func (t *Modrinth) Profiles() ([]*Profile, error) {
	profileDir := filepath.Join(t.installDir, "profiles")
	dirEntries, err := os.ReadDir(profileDir)
	if err != nil {
		return nil, fmt.Errorf("failed listing profiles directory: %w", err)
	}

	profiles := make([]*Profile, 0, len(dirEntries))
	for _, de := range dirEntries {
		if !de.IsDir() {
			continue
		}
		profile, err := ParseProfile(filepath.Join(profileDir, de.Name()))
		if err != nil {
			return nil, fmt.Errorf("failed parsing profile: %w", err)
		}
		profiles = append(profiles, profile)
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

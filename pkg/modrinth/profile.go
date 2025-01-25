package modrinth

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Profile struct {
	Name     string
	Uuid     string    `json:"uuid"`
	Path     string    `json:"path"`
	Metadata *Metadata `json:"metadata"`

	HasProfileFile bool   `json:"-"`
	FullPath       string `json:"-"`
}

type Metadata struct {
	Name         string    `json:"name"`
	GameVersion  string    `json:"game_version"`
	Loader       string    `json:"loader"`
	DateCreated  time.Time `json:"date_created"`
	DateModified time.Time `json:"date_modified"`
}

func ParseProfile(path string) (*Profile, error) {
	profileFile := filepath.Join(path, "profile.json")

	var profile Profile
	profile.Name = filepath.Base(path)
	profile.FullPath = path

	f, err := os.Open(profileFile)
	if os.IsNotExist(err) {
		return &profile, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed opening file: %w", err)
	}
	defer f.Close()

	err = json.NewDecoder(f).Decode(&profile)
	if err != nil {
		return nil, fmt.Errorf("failed decoding json: %w", err)
	}

	profile.HasProfileFile = true

	return &profile, nil
}

func (t *Profile) String() string {
	var b strings.Builder

	b.WriteString(t.Name)

	if t.Uuid != "" {
		fmt.Fprintf(&b, " [%s]", t.Uuid)
	}

	return b.String()
}

package modrinth

import (
	"fmt"
	"strings"
	"time"
)

type Profile struct {
	Path             string
	Name             string
	GameVersion      string
	ModLoader        string
	ModLoaderVersion *string
	Created          time.Time
	Modified         time.Time
	LastPlayed       time.Time

	FullPath string
}

func (t *Profile) String() string {
	var b strings.Builder

	b.WriteString(t.Name)

	if t.Name != "" {
		fmt.Fprintf(&b, " (%s)", t.Name)
	}

	return b.String()
}

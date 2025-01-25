package modrinth

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/glebarez/go-sqlite"
)

type AppDb struct {
	db *sql.DB
}

func NewAppDb(path string) (*AppDb, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, fmt.Errorf("failed opening database: %w", err)
	}

	t := AppDb{
		db: db,
	}

	return &t, nil
}

func (t *AppDb) Close() error {
	return t.db.Close()
}

func (t *AppDb) Profiles() (profiles []*Profile, err error) {
	rows, err := t.db.Query(`
		SELECT path, name, game_version, mod_loader, mod_loader_version, created, modified, last_played
		FROM profiles`)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var (
			p                             Profile
			created, modified, lastPlayed int64
		)
		err = rows.Scan(&p.Path, &p.Name, &p.GameVersion, &p.ModLoader, &p.ModLoaderVersion, &created, &modified, &lastPlayed)
		if err != nil {
			return nil, err
		}

		p.Created = time.Unix(created, 0)
		p.Modified = time.Unix(modified, 0)
		p.LastPlayed = time.Unix(lastPlayed, 0)

		profiles = append(profiles, &p)
	}

	return profiles, nil
}

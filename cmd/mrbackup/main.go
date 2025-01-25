package main

import (
	"bytes"
	"fmt"
	"html/template"
	"path/filepath"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/zekrotja/mrbackup/pkg/modrinth"
	"github.com/zekrotja/mrbackup/pkg/util"
	"github.com/zekrotja/rogu/level"
	"github.com/zekrotja/rogu/log"
)

const backupFileTimeFormat = "2006-01-02_15-04-05"

type Args struct {
	Profile string `arg:"positional,required" help:"The profile to back up"`
	Target  string `arg:"positional,required,env:TARGET_PATH" help:"Target path for backup files"`

	FilenameFormat string `arg:"--filename-format,-f,env:FILENAME_FORMAT" help:"Format of the file name" default:"{{.Profile.Name}}_{{.Profile.GameVersion}}_{{.Timestamp}}.zip"`

	InstallDir string      `arg:"--install-dir,env:MODRINTH_DIR" help:"Custom Modrinth install directory"`
	LogLevel   level.Level `arg:"--log-level,-l" help:"Log level" default:"info"`
}

func main() {
	var args Args
	arg.MustParse(&args)

	log.SetLevel(args.LogLevel)

	mr, err := modrinth.New(args.InstallDir)
	if err != nil {
		log.Fatal().Err(err).Msg("failed initializing Modrinth handler")
	}
	defer mr.Close()

	profiles, err := mr.Profiles()
	if err != nil {
		log.Fatal().Err(err).Msg("failed getting profile list")
	}

	log.Debug().Field("profiles", profiles).Msg("found profiles")

	profile, ok := findProfile(profiles, args.Profile)
	if !ok {
		log.Fatal().Fields("name", args.Profile).Msg("no profile found with given name")
	}
	if err != nil {
		log.Fatal().Fields("name", args.Profile).Msg("no profile found with given name")
	}

	targetFileName, err := getTargetFileName(args.FilenameFormat, profile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed parsing target file name")
	}

	targetFile := filepath.Join(args.Target, targetFileName)
	log.Info().Msg("creating backup zip archive ...")
	err = util.ZipDirectory(profile.FullPath, targetFile)
	if err != nil {
		log.Fatal().Err(err).Msg("failed creating zip archive")
	}

	log.Info().Field("dir", targetFile).Msg("backup created")
}

func getTargetFileName(nameTemplate string, profile *modrinth.Profile) (string, error) {
	timestamp := time.Now().Format(backupFileTimeFormat)

	tpl, err := template.New("main").Parse(nameTemplate)
	if err != nil {
		return "", fmt.Errorf("failed parsing template: %w", err)
	}

	type payload struct {
		Timestamp string
		Profile   *modrinth.Profile
	}

	var buf bytes.Buffer
	err = tpl.Execute(&buf, payload{Timestamp: timestamp, Profile: profile})
	if err != nil {
		return "", fmt.Errorf("failed executing template: %w", err)
	}

	return buf.String(), nil
}

func findProfile(profiles []*modrinth.Profile, name string) (*modrinth.Profile, bool) {
	for _, profile := range profiles {
		if profile.Name == name || profile.Path != "" && profile.Path == name {
			return profile, true
		}
	}
	return nil, false
}

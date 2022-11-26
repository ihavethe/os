package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	"github.com/hjson/hjson-go/v4"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/samber/lo"

	"github.com/sveltinio/prompti/choose"
)

type PackageInfo struct {
	Name     string            `json:"name"`
	Versions map[string]string `json:"versions"`
}

var packages []PackageInfo

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	configFile, err := os.ReadFile("packages.hjson")
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to read package manifest")
	}

	err = hjson.Unmarshal(configFile, &packages)
	if err != nil {
		log.Fatal().
			Err(err).
			Msg("Failed to parse package manifest")
	}

	pkg, ok := lo.Find(packages, func(pkg PackageInfo) bool {
		return pkg.Name == os.Args[2]
	})
	if !ok {
		log.Fatal().
			Str("package", os.Args[2]).
			Msg("Could not find package")
	}

	versionSelectionPrompt := &choose.Config{
		Title:    "Please choose a version:",
		ErrorMsg: "Please select a valid version.",
	}

	entries := lo.MapToSlice(
		pkg.Versions,
		func(ver string, cid string) list.Item {
			return choose.Item{Name: ver, Desc: ver}
		})

	result, err := choose.Run(versionSelectionPrompt, entries)
	fmt.Println(result)
}

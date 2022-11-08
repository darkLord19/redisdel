// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package commands

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	Version = "1.0"
	// Build date in ISO8601 format, output of $(date -u +'%Y-%m-%dT%H:%M:%SZ')
	buildDate = "unknown"
)

var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of redisdel.",
	RunE:  versionCmdF,
}

func init() {
	RootCmd.AddCommand(VersionCmd)
}

type Info struct {
	Version   string
	BuildDate string
	GoVersion string
	Compiler  string
	Platform  string
}

func getVersionInfo() Info {
	return Info{
		Version:   Version,
		BuildDate: buildDate,
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

func versionCmdF(cmd *cobra.Command, args []string) error {
	v := getVersionInfo()
	fmt.Printf("redisdel:\nVersion:\t%v\nBuildDate:\t%v\nGoVersion:\t%v"+
		"\nCompiler:\t%v\nPlatform:\t%v", v.Version, v.BuildDate, v.GoVersion, v.Compiler, v.Platform)
	return nil
}

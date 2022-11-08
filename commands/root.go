// Copyright (c) 2015-present Mattermost, Inc. All Rights Reserved.
// See LICENSE.txt for license information.

package commands

import (
	"fmt"
	"os"
	"runtime/debug"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var RootCmd = &cobra.Command{
	Use:               "redisdel",
	Short:             "Command line tool to delete given key patterns from your redis server, cluster or sentinel",
	DisableAutoGenTag: true,
	SilenceUsage:      true,
}

func Run(args []string) error {
	viper.SetEnvPrefix("mmctl")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// RootCmd.PersistentFlags().String("config", filepath.Join(xdgConfigHomeVar, configParent, configFileName), "path to the configuration file")
	// _ = viper.BindPFlag("config", RootCmd.PersistentFlags().Lookup("config"))
	// RootCmd.PersistentFlags().String("config-path", xdgConfigHomeVar, "path to the configuration directory.")
	// _ = viper.BindPFlag("config-path", RootCmd.PersistentFlags().Lookup("config-path"))
	// _ = RootCmd.PersistentFlags().MarkHidden("config-path")

	RootCmd.SetArgs(args)

	defer func() {
		if x := recover(); x != nil {
			fmt.Println("Uh oh! Something unexpected happened :( Would you mind reporting it?\n")
			fmt.Println(`https://github.com/mattermost/mmctl/issues/new?title=%5Bbug%5D%20panic%20on%20mmctl%20v` + Version + "&body=%3C!---%20Please%20provide%20the%20stack%20trace%20--%3E\n")
			fmt.Println(string(debug.Stack()))

			os.Exit(1)
		}
	}()

	return RootCmd.Execute()
}

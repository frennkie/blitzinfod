package main

import (
	"github.com/frennkie/blitzd/internal/config"
	"github.com/frennkie/blitzd/pkg/cmd/blitzd"
	"github.com/spf13/cobra"
)

func main() {
	cobra.OnInitialize(config.InitConfig)

	blitzd.RootCmd.PersistentFlags().StringVar(&config.BlitzdDir, "dir",
		config.DefaultBlitzdDir, "blitzd home directory")

	blitzd.RootCmd.PersistentFlags().BoolVarP(&config.Verbose, "verbose", "v",
		false, "print more log messages")

	blitzd.RootCmd.AddCommand(blitzd.DemoCmd)
	blitzd.RootCmd.AddCommand(blitzd.GenCertCmd)
	blitzd.RootCmd.AddCommand(blitzd.GraceCmd)

	_ = blitzd.RootCmd.Execute()

}

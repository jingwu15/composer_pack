package cmd

import (
	//"fmt"
	"github.com/spf13/cobra"
	sv "github.com/jingwu15/composer_pack/service"
)

var checkCmd = &cobra.Command {
	Use:   "check",
	Short: "check the composer update",
	Long:  `check the composer update`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperServer()
		sv.Check()
	},
}

var packCmd = &cobra.Command {
	Use:   "pack",
	Short: "pack composer.json,composer.lock,vendor to md5.tar.gz",
	Long:  `pack composer.json,composer.lock,vendor to md5.tar.gz`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperServer()
		sv.Pack()
	},
}

var upCmd = &cobra.Command {
	Use:   "up",
	Short: "up the *.tar.gz file to the file server",
	Long:  `up the *.tar.gz file to the file server`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperServer()
		sv.Up()
	},
}

var serverCmd = &cobra.Command {
	Use:   "server",
	Short: "the file server",
	Long:  `the file server`,
	//Run: func(cmd *cobra.Command, args []string) {
	//	//mergeViperServer()
	//	//sv.Up()
	//},
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "run the server, use nohup ... &",
	Long:  `run the server, use nohup ... &`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperServer()
		sv.Run()
	},
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "start the server, use nohup ... &",
	Long:  `start the server, use nohup ... &`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperServer()
		sv.Start()
	},
}

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "restart the crontabd server",
	Long:  `restart the crontabd server `,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperServer()
		sv.Restart()
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop the server",
	Long:  `stop the server`,
	Run: func(cmd *cobra.Command, args []string) {
		mergeViperServer()
		sv.Stop()
	},
}

func init() {
	serverCmd.AddCommand(runCmd)
	serverCmd.AddCommand(startCmd)
	serverCmd.AddCommand(restartCmd)
	serverCmd.AddCommand(stopCmd)

	rootCmd.AddCommand(serverCmd)
	rootCmd.AddCommand(checkCmd)
	rootCmd.AddCommand(packCmd)
	rootCmd.AddCommand(upCmd)
}

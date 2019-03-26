package cmd

import (
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "augit",
}

var log *logrus.Logger

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func RegisterLogger(lgr *logrus.Logger) {
	log = lgr
}

func generateSuccessString(cmdName string) string {
	return fmt.Sprintf("100 : %s Cron Success", cmdName)
}

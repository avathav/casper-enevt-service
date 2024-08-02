package cmd

import (
	"github.com/spf13/cobra"
)

var serviceCmd = &cobra.Command{
	Use:   "start-event-service",
	Short: "cli that starts events service",
	Run:   startEventService,
}

func init() {
	rootCmd.AddCommand(serviceCmd)
}

func startEventService(_ *cobra.Command, _ []string) {
	//ctx := gorm.ContextWithConnection(cmd.Context(), di.GORM())
}

package cmd

import (
	emailexchange "event-service/email/echange"
	"event-service/internal/di"

	"github.com/spf13/cobra"
)

var consumerCmd = &cobra.Command{
	Use:   "start-email-service",
	Short: "cli that starts amqp consumer",
	Run:   startConsumer,
}

func init() {
	rootCmd.AddCommand(consumerCmd)
}

func startConsumer(_ *cobra.Command, _ []string) {
	eventUpdateConsumer := di.NewEventUpdateConsumer()

	handler := emailexchange.NewEventQueryHandler()

	if err := eventUpdateConsumer.Consume(handler.Handle); err != nil {
		panic(err)
	}
}

package main

import (
	"log"

	"github.com/alifcapital/keycloak_module/cmd/amqp"
	"github.com/spf13/cobra"
)

func runCobra() {
	starter := &cobra.Command{
		Use:   "start [amqp]",
		Short: "Starts rabbitmq consumer",
		Long:  "All configuration parsed via .env",
		Args:  cobra.MinimumNArgs(1),
	}

	amqpApp := &cobra.Command{
		Use:   "amqp",
		Short: "Start listening for messages.",
		Long:  "Starts rabbitmq consumer which takes user CRUD commands from HR and broadcasts changes.",
		Run: func(cmd *cobra.Command, args []string) {
			amqp.Run()
		},
	}

	rootCmd := &cobra.Command{Use: "start"}
	rootCmd.AddCommand(starter)
	starter.AddCommand(amqpApp)

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

package cmd

import (
"context"
"log"

"github.com/spf13/cobra"
"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:              "governet",
	Short:            "Governet FEC Fetcher",
	Long:             `Use the Governet FEC Fetcher to fetch FEC bulk data files from the FEC bulk data stores`,
}

func init() {
	rootCmd.PersistentFlags().IntP("verbosity", "v", 0, "Set the log level verbosity")
	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Fatalf("failed to bind flags for root: %v", err)
	}
}

func Execute() error {
	return rootCmd.ExecuteContext(context.Background())
}

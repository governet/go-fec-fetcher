package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type fetchDataOptions struct {
	filename string
}

var fd = &fetchDataOptions{}

var fetchDataCmd = &cobra.Command{
	Use:          "cluster",
	Short:        "Create workload cluster",
	Long:         "This command is used to create workload clusters",
	PreRunE:      preRunCreateCluster,
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := fd.fetchData(cmd.Context()); err != nil {
			return fmt.Errorf("failed to create cluster: %v", err)
		}
		return nil
	},
}

func init() {
	fetchCmd.AddCommand(fetchDataCmd)
	fetchDataCmd.Flags().StringVarP(&fd.filename, "filename", "f", "", "Filename that contains EKS-A cluster configuration")
	err := fetchDataCmd.MarkFlagRequired("filename")
	if err != nil {
		log.Fatalf("Error marking flag as required: %v", err)
	}
}

func preRunCreateCluster(cmd *cobra.Command, args []string) error {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		err := viper.BindPFlag(flag.Name, flag)
		if err != nil {
			log.Fatalf("Error initializing flags: %v", err)
		}
	})
	return nil
}

func (fd *fetchDataOptions) fetchData(ctx context.Context) error {
	fmt.Println("THIS IS ANNOYING")
	return nil
}

var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "fetch resources",
	Long:  "fetch resources",
}

func init() {
	rootCmd.AddCommand(fetchCmd)
}
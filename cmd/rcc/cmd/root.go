package cmd

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/evertras/rcc/pkg/repository"
	"github.com/evertras/rcc/pkg/server"
)

var rootCmd = &cobra.Command{
	Use:   "rcc",
	Short: "Host and serve code coverage badges",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		var repo server.CoverageRepository
		var err error

		switch config.StorageType {
		case "in-memory":
			log.Println("Using in-memory storage")
			repo = repository.NewInMemory()

		case "file":
			log.Println("Using local file storage in directory", config.FileStorageBaseDir)
			repo = repository.NewFile(config.FileStorageBaseDir)

		case "dynamodb":
			log.Println("Using DynamoDB storage")
			repo, err = repository.NewDynamoDB(config.DynamoDBTable)

			if err != nil {
				return fmt.Errorf("failed to create DynamoDB repository: %w", err)
			}

		default:
			return fmt.Errorf("unknown storage type %q", config.StorageType)
		}
		cfg := server.NewDefaultConfig()

		s := server.New(cfg, repo)

		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		err = s.ListenAndServe(ctx)

		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	flags := rootCmd.Flags()

	flags.StringP("storage-type", "s", "in-memory", "The storage type to use.  One of 'in-memory' or 'file', or 'dynamodb'")
	flags.String("dynamodb-table", "evertras-rcc", "The DynamoDB table to use, when using the dynamodb storage type")
	flags.String("file-storage-base-dir", "./rcc-storage", "The base directory to store files in when using local file storage")

	viper.BindPFlags(flags)
}

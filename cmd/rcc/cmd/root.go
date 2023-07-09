package cmd

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/spf13/cobra"

	"github.com/evertras/rcc/pkg/repository"
	"github.com/evertras/rcc/pkg/server"
)

var rootCmd = &cobra.Command{
	Use:   "rcc",
	Short: "Host and serve code coverage badges",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo := repository.NewInMemory()
		cfg := server.NewDefaultConfig()

		s := server.New(cfg, repo)

		ctx, cancel := signal.NotifyContext(cmd.Context(), os.Interrupt)
		defer cancel()

		err := s.ListenAndServe(ctx)

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

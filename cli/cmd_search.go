package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) searchCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "search <query>",
		Short: "Search for videos on iQIYI",
		Long:  `Search iQIYI for movies, TV series, and other content.`,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			query := args[0]
			n := a.effectiveLimit(20)
			a.progressf("searching iQIYI for %q...", query)
			videos, err := a.client.Search(cmd.Context(), query, n)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(videos, len(videos))
		},
	}
}

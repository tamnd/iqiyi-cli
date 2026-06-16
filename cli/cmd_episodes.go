package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) episodesCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "episodes <album-id>",
		Short: "List episodes for a show by album ID",
		Long: `Fetch the episode list for an iQIYI series by its album/qipu ID.

The album ID is the numeric identifier found in iQIYI URLs:
  https://www.iqiyi.com/a_<album-id>.html`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			albumID := args[0]
			n := a.effectiveLimit(0)
			a.progressf("fetching episodes for album %s...", albumID)
			episodes, err := a.client.Episodes(cmd.Context(), albumID, n)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(episodes, len(episodes))
		},
	}
}

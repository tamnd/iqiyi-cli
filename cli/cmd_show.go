package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) showCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "show <album-id>",
		Short: "Fetch detail for a show or movie by album ID",
		Long: `Fetch full metadata for an iQIYI album (series, movie, etc.) by its album/qipu ID.

The album ID is the numeric identifier found in iQIYI URLs:
  https://www.iqiyi.com/a_<album-id>.html`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			albumID := args[0]
			a.progressf("fetching show %s...", albumID)
			video, err := a.client.Show(cmd.Context(), albumID)
			if err != nil {
				return mapFetchErr(err)
			}
			if video == nil {
				return codeError(exitNoData, nil)
			}
			return a.render([]any{video})
		},
	}
}

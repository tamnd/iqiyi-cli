package cli

import (
	"github.com/spf13/cobra"
)

func (a *App) hotCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "hot",
		Short: "List trending/hot videos from iQIYI",
		Long:  `Fetch the current hot/trending content list used by the iQIYI homepage.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			a.progressf("fetching hot list...")
			videos, err := a.client.Hot(cmd.Context())
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(videos, len(videos))
		},
	}
}

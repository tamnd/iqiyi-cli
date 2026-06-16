package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

var validRankTypes = []string{"tvSeries", "movie", "variety", "cartoon", "documentary"}

func (a *App) rankCmd() *cobra.Command {
	var period, dimension, contentType string
	cmd := &cobra.Command{
		Use:   "rank",
		Short: "Show content rankings",
		Long: `Fetch content rankings from iQIYI.

Period:    day|week|month
Dimension: hot|playcount|score
Type:      tvSeries|movie|variety|cartoon|documentary`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			a.progressf("fetching rank list (period=%s dimension=%s type=%s)...", period, dimension, contentType)
			items, err := a.client.Rank(cmd.Context(), period, dimension, contentType)
			if err != nil {
				return mapFetchErr(err)
			}
			n := a.effectiveLimit(len(items))
			if n < len(items) {
				items = items[:n]
			}
			if len(items) == 0 {
				return fmt.Errorf("no rank data returned")
			}
			return a.renderOrEmpty(items, len(items))
		},
	}
	cmd.Flags().StringVar(&period, "period", "day", "ranking period: day|week|month")
	cmd.Flags().StringVar(&dimension, "dimension", "hot", "ranking dimension: hot|playcount|score")
	cmd.Flags().StringVar(&contentType, "type", "tvSeries", "content type: tvSeries|movie|variety|cartoon|documentary")
	return cmd
}

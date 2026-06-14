package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/tamnd/iqiyi-cli/iqiyi"
)

func (a *App) listCmd() *cobra.Command {
	var channel string
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List recommended shows for a channel",
		Long: `List the top recommended shows for an iQIYI channel.

Available channels: movie, tv_drama, animation, kids, documentary, variety`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			chID, ok := iqiyi.ChannelID(channel)
			if !ok {
				return codeError(exitUsage, fmt.Errorf("unknown channel %q -- use one of: movie, tv_drama, animation, kids, documentary, variety", channel))
			}
			n := a.effectiveLimit(20)
			a.progressf("fetching %d shows from channel %q...", n, channel)
			shows, err := a.client.List(cmd.Context(), chID, n)
			if err != nil {
				return mapFetchErr(err)
			}
			return a.renderOrEmpty(shows, len(shows))
		},
	}
	cmd.Flags().StringVar(&channel, "channel", "movie", "channel: movie|tv_drama|animation|kids|documentary|variety")
	return cmd
}

func (a *App) channelsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "channels",
		Short: "List available channels with their IDs",
		RunE: func(_ *cobra.Command, _ []string) error {
			channels := iqiyi.Channels()
			return a.render(channels)
		},
	}
}

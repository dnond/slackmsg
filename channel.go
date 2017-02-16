package slackmsg

import (
	"github.com/k0kubun/pp"
	"github.com/nlopes/slack"
)

func GetChannels(api *slack.Client, channelNames []string) ([]slack.Channel, error) {
	channels, err := api.GetChannels(true)
	if err != nil {
		pp.Println("GetChannels")
		return nil, err
	}

	var targetChannels []slack.Channel
	for _, channel := range channels {
		for _, channelName := range channelNames {
			if channel.Name == channelName {
				targetChannels = append(targetChannels, channel)
				// break
			}
		}
	}

	return targetChannels, nil
}

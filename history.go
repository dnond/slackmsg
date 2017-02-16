package slackmsg

import (
	"fmt"
	"log"

	"github.com/k0kubun/pp"
	"github.com/nlopes/slack"
	"github.com/syndtr/goleveldb/leveldb"
)

type slackHistoryLatest struct {
	channelId string
	latest    string
}

func saveLatestTime(channel slack.Channel, history *slack.History) error {
	historyLatest := NewSlackHistoryLatest(channel.ID, history.Latest)
	err := historyLatest.save()
	return err
}

func NewSlackHistoryLatest(channelId, latest string) *slackHistoryLatest {
	historyLatest := &slackHistoryLatest{
		channelId: channelId,
		latest:    latest,
	}
	return historyLatest
}

const LATEST_SAVE_FILE = "/tmp/slack_latest"

func (historyLatest *slackHistoryLatest) save() error {
	db, err := leveldb.OpenFile(LATEST_SAVE_FILE, nil)
	defer db.Close()

	if err != nil {
		if err.Error() == "leveldb: not found" {
		} else {
			log.Fatal(err)
			return err
		}
	}
	err = db.Put([]byte(historyLatest.channelId), []byte(historyLatest.latest), nil)
	pp.Println("save!!:" + historyLatest.latest)
	return err
}

func (historyLatest *slackHistoryLatest) createSaveString() string {
	return fmt.Sprintf("%s\t%s", historyLatest.channelId, historyLatest.latest)
}

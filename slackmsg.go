package slackmsg

import (
	"log"
	"time"

	"github.com/k0kubun/pp"
	"github.com/nlopes/slack"
	"github.com/syndtr/goleveldb/leveldb"
)

type SlackMessages []slack.Message

// sort interface
func (s SlackMessages) Len() int {
	return len(s)
}

// sort interface
func (s SlackMessages) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// sort interface
func (s SlackMessages) Less(i, j int) bool {
	return s[i].Timestamp < s[j].Timestamp
}

//////////

func GetMesssages(api *slack.Client, channel slack.Channel, resume bool) (SlackMessages, error) {
	historyParams, err := getHistoryParams(channel, resume)
	pp.Println(historyParams)
	if err != nil {
		pp.Println("getHistoryParams")
		log.Fatal(err)
		return nil, err
	}

	history, err := api.GetChannelHistory(channel.ID, *historyParams)
	if err != nil {
		pp.Println("GetChannelHistory")
		log.Fatal(err)
		return nil, err
	}

	var newMsg SlackMessages
	for _, message := range history.Messages {
		message.Channel = channel.Name
		newMsg = append(newMsg, message)
	}

	if newMsg != nil {
		historyLatest := NewSlackHistoryLatest(channel.ID, history.Latest)
		err = historyLatest.save()
		if err != nil {
			pp.Println("NewSlackHistoryLatest")
			log.Fatal(err)
			return nil, err
		}
	}

	return newMsg, nil
}

func getHistoryParams(channel slack.Channel, resume bool) (*slack.HistoryParameters, error) {
	var err error
	var oldest = "0"
	var latest string

	// oldestの決定
	if resume == true {
		oldest, err = getOldest(channel)
		if err != nil {
			return nil, err
		}
	}

	// latestの決定
	latest = getLatest()

	// var historyParams slack.HistoryParameters
	historyParams := slack.HistoryParameters{
		Latest:    latest,
		Oldest:    oldest, //キャッシュ等を使って最後に取得した時間にしたい
		Count:     100,    //ポーリングするので多くなくていい
		Inclusive: true,   //取得したメッセージのLatest・Oldestのtimestampを入れるか
		Unreads:   true,   //新規メッセージの数を入れるか
	}

	return &historyParams, nil
}

func getOldest(channel slack.Channel) (string, error) {
	oldest := "0"

	db, err := leveldb.OpenFile(LATEST_SAVE_FILE, nil)
	defer db.Close()
	if err != nil {
		if err.Error() == "file missing" {
			return oldest, nil
		} else if err != nil {
			pp.Println("getOldest")
			log.Fatal(err)
			return "", err
		}
	}

	data, err := db.Get([]byte(channel.ID), nil)
	if err != nil {
		if err.Error() == "leveldb: not found" {
			pp.Println("leveldb: not found")
		} else {
			pp.Println(err)

			log.Fatal(err)
			return "", err
		}
	} else {
		oldest = string(data)
	}
	return oldest, nil
}

func getLatest() string {
	//latestは常に現在時刻
	return time.Now().Format("2006-01-02 15:04:05")
}

package qiita

import (
	"mackerel-adcal-counter/counter"
	"strconv"
	"strings"

	"github.com/gocolly/colly"
)

type Counter struct {
	url string
}

var _ counter.Counter = &Counter{}

func NewCounter(url string) *Counter {
	return &Counter{
		url: url,
	}
}

func (c *Counter) Count() (*counter.CountResult, error) {
	entriesCnt := uint64(0)
	participantsCnt := uint64(0)
	subscribersCnt := uint64(0)

	co := colly.NewCollector(
		colly.AllowedDomains("qiita.com"),
	)

	// TODO: カレンダーが2ページ以上になった時の対応
	entriesQuery := "main > div:nth-child(3) > div > div:nth-child(2) > div > div:nth-child(2) > div > table:nth-child(1) td > div:nth-child(2) > div:nth-child(1)"
	co.OnHTML(entriesQuery, func(e *colly.HTMLElement) {
		entriesCnt++
	})

	participantsQuery := "main > div:nth-child(3) > div > div:nth-child(1) > div:nth-child(1) > div:nth-child(3) > div > div:nth-child(1) > span:nth-child(3)"
	co.OnHTML(participantsQuery, func(e *colly.HTMLElement) {
		text := e.DOM.Text()
		cntStr := strings.TrimSpace(strings.Replace(text, "People", "", 1))
		cnt, err := strconv.ParseUint(cntStr, 10, 64)
		if err != nil {
			return
		}
		participantsCnt = cnt
	})

	subscribersQuery := "main > div:nth-child(3) > div > div:nth-child(1) > div:nth-child(1) > div:nth-child(3) > div > div:nth-child(2) > span:nth-child(3)"
	co.OnHTML(subscribersQuery, func(e *colly.HTMLElement) {
		text := e.DOM.Text()
		cntStr := strings.TrimSpace(strings.Replace(text, "People", "", 1))
		cnt, err := strconv.ParseUint(cntStr, 10, 64)
		if err != nil {
			return
		}
		subscribersCnt = cnt
	})

	co.Visit(c.url)
	result := &counter.CountResult{
		Entries:      entriesCnt,
		Participants: participantsCnt,
		Subscribers:  subscribersCnt,
	}
	return result, nil
}

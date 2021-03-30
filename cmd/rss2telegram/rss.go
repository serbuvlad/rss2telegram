package main

import (
	"fmt"
	"log"
	"time"

	"github.com/mmcdole/gofeed"
)

type Feed struct {
	Link string
	Top  int
}

func serveFeeds(feeds []Feed, timeout time.Duration, out chan<- string) {
	fp := gofeed.NewParser()

	for {
		for _, f := range feeds {
			feed, err := fp.ParseURL(f.Link)
			if err != nil {
				log.Println(fmt.Errorf("failed to read feed %s: %w", f.Link))
				continue
			}

			if f.Top != 0 {
				feed.Items = feed.Items[:f.Top]
			}

			for _, i := range feed.Items {
				if i.GUID == "" {
					i.GUID = i.Link
				}

				if haveEntry(f.Link, i.GUID) {
					continue
				}

				out <- i.Link

				addEntry(f.Link, i.GUID)
			}
		}

		time.Sleep(timeout)
	}
}

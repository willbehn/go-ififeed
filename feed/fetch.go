package feed

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"

	h2m "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/mmcdole/gofeed"
	"github.com/willbehn/go-ifi-feed/models"
)

type Message struct {
	Timestamp time.Time
	Content   string
	Course    string
}

var httpClient = &http.Client{
	Transport: &http.Transport{},
	Timeout:   15 * time.Second,
}

func fetchRssFeed(course string) *gofeed.Feed {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL(fmt.Sprintf(
		"https://www.uio.no/studier/emner/matnat/ifi/%s/h25/beskjeder/?vrtx=feed",
		course,
	))
	if err != nil {
		return nil
	}

	return feed
}

func fetchHttpItem(resp *http.Response, course string, title string) string {
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Println("Html parse error", err)
		return ""
	}

	headerSelection := doc.Find("h1")
	contentSelection := doc.Find("#vrtx-field-message")

	if headerSelection.Length() == 0 || contentSelection.Length() == 0 {
		fmt.Println("Error finding html tag(s)")
		return ""
	}

	htmlH, _ := headerSelection.Html()
	htmlC, _ := contentSelection.Html()

	htmlCourse := fmt.Sprintf(
		`<p><code>%s</code> – <b>%s</b></p>`,
		course,
		title,
	)

	return fmt.Sprintf("<h1>🔔 %s</h1>%s%s", htmlH, htmlCourse, htmlC)
}

func ConvertToMarkdown(html string) string {
	converter := h2m.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)

	if err != nil {
		fmt.Println("Error converting html to markdown:", err)
		return ""
	}
	return markdown
}

func singleFeed(course models.Course) []Message {
	var results []Message

	feed := fetchRssFeed(course.Code)

	for _, item := range feed.Items {
		resp, err := httpClient.Get(item.Link)
		if err != nil {
			fmt.Println("Error fetching link", err)
			continue
		}

		html := fetchHttpItem(resp, course.Code, course.Title)
		resp.Body.Close()

		timePublished := item.PublishedParsed
		htmlTs := html + "<p>" + timePublished.Format("2006-01-02 15:04") + "</p><hr/>"

		newMessage := Message{Content: htmlTs, Timestamp: *timePublished}

		results = append(results, newMessage)
	}
	return results
}

func Fetch(courses models.Courses) []Message {
	var results []Message

	msgChannel := make(chan Message)
	var wg sync.WaitGroup
	wg.Add(len(courses.Courses))

	for _, course := range courses.Courses {
		go func(c models.Course) {
			defer wg.Done()
			for _, msg := range singleFeed(c) {
				msgChannel <- msg
			}
		}(course)
	}

	go func() {
		wg.Wait()
		close(msgChannel)
	}()

	for m := range msgChannel {
		results = append(results, m)
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.After(results[j].Timestamp)
	})

	return results
}

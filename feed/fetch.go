package feed

import (
	"fmt"
	"net/http"
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

func fetchHttpItem(resp *http.Response, course string, title string, published time.Time) string {
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

	htmlCourse := fmt.Sprintf(`<p><code>%s</code> - %s</p>`, course, title)
	htmlDate := fmt.Sprintf(`<p><em>%s</em></p>`, published.Format("2 Jan 2006, 15:04"))

	return fmt.Sprintf("%s<h1>%s</h1>%s%s<hr/>", htmlCourse, htmlH, htmlDate, htmlC)
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
	if feed == nil {
		return nil
	}

	for _, item := range feed.Items {
		resp, err := httpClient.Get(item.Link)
		if err != nil {
			fmt.Println("Error fetching link", err)
			continue
		}

		timePublished := item.PublishedParsed
		html := fetchHttpItem(resp, course.Code, course.Title, *timePublished)
		resp.Body.Close()

		newMessage := Message{Content: html, Timestamp: *timePublished}

		results = append(results, newMessage)
	}
	return results
}

func FetchStream(courses models.Courses, out chan<- Message) {
	var wg sync.WaitGroup
	wg.Add(len(courses.Courses))

	for _, course := range courses.Courses {
		go func(c models.Course) {
			defer wg.Done()
			for _, msg := range singleFeed(c) {
				out <- msg
			}
		}(course)
	}

	wg.Wait()
	close(out)
}

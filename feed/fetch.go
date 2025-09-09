package feed

import (
	"fmt"
	"net/http"
	"sort"
	"time"

	h2m "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/glamour"
	"github.com/mmcdole/gofeed"
)

type Message struct {
	Timestamp time.Time
	Content   string
	Course    string
}

func fetchRssFeed(course string) *gofeed.Feed {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL("https://www.uio.no/studier/emner/matnat/ifi/" + course + "/h25/beskjeder/?vrtx=feed")
	if err != nil {
		return nil
	}
	return feed
}

func fetchHttpItem(resp *http.Response) string {
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

	return "<h1>" + htmlH + "</h1>" + htmlC
}

func convertToMarkdown(html string) string {
	converter := h2m.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)

	if err != nil {
		fmt.Println("Error converting html to markdown:", err)
		return ""
	}
	return markdown
}

func Fetch(courses []string) []Message {
	var results []Message

	for _, course := range courses {

		feed := fetchRssFeed(course)

		for _, item := range feed.Items {
			resp, err := http.Get(item.Link)
			if err != nil {
				fmt.Println("Error fetching link", err)
				continue
			}
			html := fetchHttpItem(resp)
			resp.Body.Close()

			timePublished := item.PublishedParsed
			htmlTs := html + "<p>" + timePublished.Format("2006-01-02 15:04") + "</p>"

			markdown := convertToMarkdown(htmlTs)

			markdown = "## " + course + "\n" + markdown

			out, err := glamour.Render(markdown, "dark")
			if err != nil {
				fmt.Println("Error rendering markdown:", err)
				continue
			}

			newMessage := Message{Content: out, Timestamp: *timePublished}

			results = append(results, newMessage)

			resp.Body.Close()
		}
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].Timestamp.After(results[j].Timestamp)
	})

	return results
}

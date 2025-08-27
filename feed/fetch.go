package feed

import (
	"fmt"
	"net/http"

	h2m "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/glamour"
	"github.com/mmcdole/gofeed"
)

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

	htmlHeader := doc.Find("h1")
	htmlContent := doc.Find("#vrtx-field-message")

	if htmlHeader.Length() == 0 || htmlContent.Length() == 0 {
		fmt.Println("Error finding html tag(s)")
		return ""
	}

	htmlH, _ := htmlHeader.Html()
	htmlC, _ := htmlContent.Html()

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

func Fetch(course string) []string {

	feed := fetchRssFeed(course)

	var results []string

	for _, item := range feed.Items {

		resp, err := http.Get(item.Link)
		if err != nil {
			fmt.Println("Error fetching link", err)
			continue
		}
		defer resp.Body.Close()

		html := fetchHttpItem(resp)
		markdown := convertToMarkdown(html)

		out, err := glamour.Render(markdown, "dark")
		if err != nil {
			fmt.Println("Error rendering markdown:", err)
			continue
		}

		results = append(results, out)

		resp.Body.Close()
	}
	return results
}

package feed

import (
	"fmt"
	"net/http"

	h2m "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/glamour"
	"github.com/mmcdole/gofeed"
)

func Fetch(course string) []string {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL("https://www.uio.no/studier/emner/matnat/ifi/" + course + "/h25/beskjeder/?vrtx=feed")
	if err != nil {
		return nil
	}

	var results []string

	for _, item := range feed.Items {

		resp, err := http.Get(item.Link)
		if err != nil {
			fmt.Println("Error fetching link", err)
			continue
		}
		defer resp.Body.Close()

		doc, err := goquery.NewDocumentFromReader(resp.Body)
		if err != nil {
			fmt.Println("Html parse error", err)
			continue
		}

		htmlSelection := doc.Find("#right-main")
		if htmlSelection.Length() == 0 {
			fmt.Println("Error findin html tag", item.Link)
			continue
		}

		html, err := htmlSelection.Html()
		if err != nil {
			continue
		}

		converter := h2m.NewConverter("", true, nil)
		markdown, err := converter.ConvertString(html)
		if err != nil {
			fmt.Println("Error converting html to markdown:", err)
			continue
		}

		out, err := glamour.Render(markdown, "dark")
		if err != nil {
			fmt.Println("Error rendering markdown:", err)
			continue
		}

		results = append(results, out)
	}
	return results
}

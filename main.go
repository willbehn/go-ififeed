package main

import (
	"fmt"
	"net/http"

	h2m "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/glamour"
	"github.com/mmcdole/gofeed"
)

func main() {
	/*
			in := `# Hello World

				This is a simple example of Markdown rendering with Glamour!
				Check out the [other examples](https://github.com/charmbracelet/glamour/tree/master/examples) too.

				Bye!
		   	 `

			out, err := glamour.Render(in, "dark")
			if err != nil {
				fmt.Println("Error rendering markdown:", err)
				return
			}
			fmt.Print(out)*/

	fp := gofeed.NewParser()

	feed, err := fp.ParseURL("https://www.uio.no/studier/emner/matnat/ifi/IN5040/h25/beskjeder/?vrtx=feed")
	if err != nil {
		return
	}

	for _, item := range feed.Items {
		//fmt.Println(item.Title, item.Link)

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
		fmt.Print(out)
	}
}

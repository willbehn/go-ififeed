package main

import (
	"fmt"
	"net/http"

	h2m "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/mmcdole/gofeed"
)

type model struct {
	vp      viewport.Model
	content string
	ready   bool
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:

		if !m.ready {
			m.vp = viewport.New(msg.Width, msg.Height-2) // leave room for header/footer
			m.vp.MouseWheelEnabled = true
			m.vp.SetContent(m.content)
			m.ready = true
		} else {
			m.vp.Width = msg.Width
			m.vp.Height = msg.Height - 2
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.vp, cmd = m.vp.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if !m.ready {
		return "loading...\n"
	}
	header := "Scrollable output (viewport)\n"
	footer := fmt.Sprintf("\nScroll: %.0f%% — press q to quit", m.vp.ScrollPercent()*100)
	return header + m.vp.View() + footer
}

func getMarkdown() string {
	fp := gofeed.NewParser()

	feed, err := fp.ParseURL("https://www.uio.no/studier/emner/matnat/ifi/IN5040/h24/beskjeder/?vrtx=feed")
	if err != nil {
		return ""
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
		//fmt.Print(out)
		return out
	}
	return ""
}

func main() {
	p := tea.NewProgram(model{content: getMarkdown()},
		tea.WithAltScreen(),
		tea.WithMouseCellMotion())

	p.Run()
}

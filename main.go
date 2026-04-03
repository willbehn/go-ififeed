package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/willbehn/go-ififeed/feed"
	"github.com/willbehn/go-ififeed/models"
	yaml "gopkg.in/yaml.v3"
)

func runTui(courses models.Courses) {
	ch := make(chan feed.Message, 50)
	go feed.FetchStream(courses, ch)

	p := tea.NewProgram(
		newModel(ch),
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)
	p.Run()
}

func cmdSubscribe(courses []string) {
	for _, course := range courses {
		fmt.Print(course)
	}
}

const sampleConfig = `
Courses:
  - code: "IN1000"
    title: "Introduksjon til objektorientert programmering"
    semester: h25

  - code: "IN1010"
    title: "Objektorientert programmering"
    semester: h25
`

func configPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("kunne ikke finne config directory %w", err)
	}
	return filepath.Join(dir, "ififeed", "courses.yaml"), nil
}

func readCourses() (models.Courses, error) {
	path, err := configPath()
	if err != nil {
		return models.Courses{}, err
	}

	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return models.Courses{}, fmt.Errorf("kunne ikke lage config directory: %w", err)
		}
		if err := os.WriteFile(path, []byte(sampleConfig), 0644); err != nil {
			return models.Courses{}, fmt.Errorf("kunne ikke skrive til config fil: %w", err)
		}
		return models.Courses{}, fmt.Errorf("lagde config.yaml på %s\nendre den for å legge til emnene du tar, deretter kjør ififeed igjen", path)
	}
	if err != nil {
		return models.Courses{}, err
	}
	defer file.Close()

	var data models.Courses
	if err := yaml.NewDecoder(file).Decode(&data); err != nil {
		return models.Courses{}, err
	}
	return data, nil
}

func main() {
	args := os.Args

	if len(args) > 1 {
		if args[1] == "subscribe" || args[1] == "-s" {
			cmdSubscribe(args[2:])
		}
	} else {
		courses, err := readCourses()
		if err != nil {
			fmt.Println(err)
			return
		}
		runTui(courses)
	}
}

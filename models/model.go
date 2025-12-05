package models

type Course struct {
	Code  string `yaml:"code"`
	Title string `yaml:"title"`
}

type Courses struct {
	Courses []Course `yaml:"Courses"`
}

const Banner = `
  ██╗███████╗██╗███████╗███████╗███████╗██████╗ 
  ██║██╔════╝██║██╔════╝██╔════╝██╔════╝██╔══██╗
  ██║█████╗  ██║█████╗  █████╗  █████╗  ██║  ██║
  ██║██╔══╝  ██║██╔══╝  ██╔══╝  ██╔══╝  ██║  ██║
  ██║██║     ██║██║     ███████╗███████╗██████╔╝v0
  ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
`

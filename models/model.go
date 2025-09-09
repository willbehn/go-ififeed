package models

type Course struct {
	Code  string `yaml:"Code"`
	Title string `yaml:"Title"`
}

type Courses struct {
	Courses []Course `yaml:"Courses"`
}

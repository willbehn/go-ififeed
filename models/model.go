package models

type Course struct {
	Code     string `yaml:"code"`
	Title    string `yaml:"title"`
	Semester string `yaml:"semester"`
}

type Courses struct {
	Courses []Course `yaml:"Courses"`
}

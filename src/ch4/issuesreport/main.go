package main

import (
	"text/template"
	"time"
	"log"
	"os"
	"ch4/github"
)

const temp1 = `{{.TotalCount}} issues:
{{range .Items}}--------------------------
Number:{{.Number}}
User:{{.User.Login}}
Title:{{.Title | printf "%.64s"}}
Age:{{.CreatedAt | daysAgo}} days
{{end}}`

func daysAgo(t time.Time) int{
	return int(time.Since(t).Hours() / 24)
}

func main(){

	report := template.Must(template.New("report").
		Funcs(template.FuncMap{"daysAgo":daysAgo}).
		Parse(temp1))

	result,err := github.SearchIssues(os.Args[1:])

	if err != nil{
		log.Fatal(err)
	}

	if err := report.Execute(os.Stdout,result); err != nil{
		log.Fatal(err)
	}
}

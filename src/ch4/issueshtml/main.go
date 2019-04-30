package main

import (
	"html/template"
	"time"
	"log"
	"os"
	"ch4/github"
)

const htmltemplate = `
<h1>{{.TotalCount}} issues</h1>
<table>
<tr style='text-align: left'>
  <th>#</th>
  <th>State</th>
  <th>User</th>
  <th>Title</th>
</tr>
{{range .Items}}
<tr>
  <td><a href='{{.HTMLURL}}'>{{.Number}}</a></td>
  <td>{{.State}}</td>
  <td><a href='{{.User.HTMLURL}}'>{{.User.Login}}</a></td>
  <td><a href='{{.HTMLURL}}'>{{.Title}}</a></td>
</tr>
{{end}}
</table>
`
func daysAgo(t time.Time) int{
	return int(time.Since(t).Hours() / 24)
}

func main(){

	report := template.Must(template.New("report").
		Funcs(template.FuncMap{"daysAgo":daysAgo}).
		Parse(htmltemplate))

	result,err := github.SearchIssues(os.Args[1:])

	if err != nil{
		log.Fatal(err)
	}

	if err := report.Execute(os.Stdout,result); err != nil{
		log.Fatal(err)
	}
}

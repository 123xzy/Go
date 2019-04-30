package main

import (
	"net/http"
	"ch4/github"
	"html/template"
	"os"
	"log"
)

func handle(w http.ResponseWriter,r *http.Request){
	result,err := github.SearchIssues(os.Args[1:])

	if err != nil{
		log.Fatal(err)
	}

	var issuelist = template.Must(template.New("issuelist").Parse(`
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
	`))

	issuelist.Execute(w,result)
}

func main(){
	http.HandleFunc("/",handle)
	http.ListenAndServe("0.0.0.0:8000",nil)
}

package main

import (
	"flag"
	"fmt"
	"os"
	"log"
	"time"

	"ch4/github"
)

// input:./ex4.10 -m=10 
var m = flag.Int("m",0,"positive: query issues created during recently m months;\nnegative: query issues created before -m months")

func main(){
	flag.Parse()
	q:= os.Args[1:]

	for i,v := range q{
		if v == "-m"{
			copy(q[i:],q[i+2:])
			q = q[0:len(q)-2]
			break
		}
	}

	if *m > 0{
		end := time.Now()
		start := end.AddDate(0,-*m,0)
		log.Println(start)
		// Github API:...?q=xxx,created:=xxx..yyy
		q = append(q,"created:" + start.Format("2019-04-29") + ".." + end.Format("2019-04-29"))
	}else if *m < 0{
		end := time.Now().AddDate(0,*m,0)
		// Github API:...?q=xxx,created:<xxx
		q = append(q,"created:<" + end.Format("2019-04-29"))
	}
	log.Println(q)

	result,err := github.SearchIssues(q)

	if err != nil{
		log.Fatal(err)
	}

	fmt.Println("%d issues\n",result.TotalCount)

	for _,item := range result.Items{
		fmt.Printf("#%-5d %9.9s %.55s\n",
			item.Number,item.User.Login,item.Title)
		}
}







package main

import(
	"fmt"
	"log"
	"os"

	"ch5/links"
)

// only 20 links.Extract can happen concurrently
var tokens = make(chan struct{},20)

func crawl(url string) []string{
	fmt.Println(url)
	tokens <- struct{}{}	// acquire a token
	list,err := links.Extract(url)
	<-tokens		// release a token

	if err != nil{
		log.Print(err)
	}
	return list
}

func main(){
	worklist := make(chan []string)

	var n int // number of pending sends to worklist

	n++
	// start with the command-line arguments.
	go func(){ worklist <- os.Args[1:] }()

	// crawl the web concurrently.
	seen := make(map[string]bool)
	for ;n > 0;n-- {
		list := <-worklist
		for _,link := range list {
			if !seen[link]{
				seen[link] = true
				n++
				go func(link string){
					worklist <- crawl(link)
				}(link)
			}
		}
	}
}

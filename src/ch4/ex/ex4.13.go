package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const apikey = "837a1b8b"
const apiurl = "http://www.omdbapi.com"

type movie_info struct{
	Title string
	Year string
	Poster string
}

func getPoster(name string){
	resp,err := http.Get(fmt.Sprintf("%s?t=%s&apikey=%s",apiurl,url.QueryEscape(name),apikey))

	if err != nil{
		fmt.Println(err)
		return
	}

	//defer resp.Boby.Close()

	info,err := ioutil.ReadAll(resp.Body)

	if err != nil{
		fmt.Println(err)
		return
	}
	minfo := movie_info{}
	err = json.Unmarshal(info,&minfo)
	if err != nil{
		fmt.Println(err)
		return
	}

	poster := minfo.Poster
	if poster != ""{
		downloadPoster(poster)
	}
	resp.Body.Close()
}

func downloadPoster(url string){

	resp,err := http.Get(url)
	if err != nil{
		fmt.Println(err)
		return
	}

	content,err := ioutil.ReadAll(resp.Body)
	if err != nil{
		fmt.Println(err)
		return
	}

	pos := strings.LastIndex(url,"/")
	if pos == -1{
		fmt.Println("failed to parse the name of the Poster's url")
		return
	}

	file,err := os.Create(url[pos+1:])
	if err != nil{
		fmt.Println(err)
		return
	}

	_,err = file.Write(content)
	if err != nil {
		fmt.Println(err)
	}

	file.Close()
}


func main(){
	names := os.Args[1:]
	for _,name := range names{
		getPoster(name)
	}
}


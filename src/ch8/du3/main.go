package main

import(
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
	"sync"
)

var verbose = flag.Bool("v",false,"show verbose progress messages")

func main(){
	// determine the initial dir
	flag.Parse()
	roots := flag.Args()

	if len(roots) == 0{
		roots = []string{"."}
	}

	// traverse the file tree
	filesize := make(chan int64)
	var n sync.WaitGroup
	for _,root := range roots{
		n.Add(1)
		go walkDir(root,&n,filesize)
	}

	go func(){
		n.Wait()
		close(filesize)
	}()


	var tick <-chan time.Time
	if *verbose{
		tick = time.Tick(time.Millisecond)
	}
	var nfiles,nbytes int64
loop:
	for{
		select{
		case size,ok := <-filesize:
			if !ok{
				break loop
			}
			nfiles++
			nbytes += size
		case <- tick:
			printDiskUsage(nfiles,nbytes)
		}
	}

	printDiskUsage(nfiles,nbytes)
}

func printDiskUsage(nfiles,nbytes int64){
	fmt.Printf("%d files %.1f GB\n",nfiles,float64(nbytes)/1e9)
}

func walkDir(dir string,n *sync.WaitGroup,filesize chan<- int64){
	defer n.Done()
	for _,entry := range dirents(dir){
		if entry.IsDir(){
			n.Add(1)
			subdir := filepath.Join(dir,entry.Name())
			go walkDir(subdir,n,filesize)
		}else{
			filesize<- entry.Size()
		}
	}
}

var sema = make(chan struct{},20)

func dirents(dir string) []os.FileInfo{
	sema <- struct{}{}		// acquire token
	defer func() { <-sema }()	// release token

	entries,err := ioutil.ReadDir(dir)
	if err != nil{
		fmt.Fprintf(os.Stderr,"du1:%v\n",err)
		return nil
	}
	return entries
}

package main

import(
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main(){
	// determine the initial dir
	flag.Parse()
	roots := flag.Args()

	if len(roots) == 0{
		roots = []string{"."}
	}

	// traverse the file tree
	filesize := make(chan int64)
	go func(){
		for _,root := range roots{
			walkDir(root,filesize)
		}
		close(filesize)
	}()

	// print the results
	var nfiles,nbytes int64
	for size := range filesize{
		nfiles++
		nbytes += size
	}
	printDiskUsage(nfiles,nbytes)
}

func printDiskUsage(nfiles,nbytes int64){
	fmt.Printf("%d files %.1f GB\n",nfiles,float64(nbytes)/1e9)
}

func walkDir(dir string,filesize chan<- int64){
	for _,entry := range dirents(dir){
		if entry.IsDir(){
			subdir := filepath.Join(dir,entry.Name())
			walkDir(subdir,filesize)
		}else{
			filesize<- entry.Size()
		}
	}
}

func dirents(dir string) []os.FileInfo{
	entries,err := ioutil.ReadDir(dir)
	if err != nil{
		fmt.Fprintf(os.Stderr,"du1:%v\n",err)
		return nil
	}
	return entries
}

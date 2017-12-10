package main

import(
	"os"
    "fmt"
	"path/filepath"
	"os/exec"
	"io/ioutil"
)

const file ="bin/temppaths"

func main() {
	if _,err := os.Stat(file); err == nil{
		os.Remove(file)
	}
	pf, err := os.Create(file)
	if err != nil{
		fmt.Println(err)
		return
	}

	for _, path := range os.Args[1:]{
		readfiles(path,pf)
	}
	pf.Close()

	cmd := exec.Command("bin/notepad++")
	err = cmd.Start()
	if err != nil{
		fmt.Println(err)
		return
	}

}

func readfiles(path string,pf *os.File){
	filepath.Walk(path,func(path string, info os.FileInfo, err error) error{		
		if !info.IsDir(){
			pf.WriteString(path+"\n")
		}
		return nil
	})
}
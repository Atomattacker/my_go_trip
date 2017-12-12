package main

import(
	"os"
    "fmt"
	"path/filepath"
	"os/exec"
)

const file ="temppaths"



func main() {

	currentpath,_ :=filepath.Abs(filepath.Dir(os.Args[0]))
	tmpfile := filepath.Join(currentpath,"bin",file)
	

	if _,err := os.Stat(tmpfile); err == nil{
		os.Remove(tmpfile)
	}
	pf, err := os.Create(tmpfile)
	if err != nil{
		fmt.Println(err)
		return
	}

	for _, path := range os.Args[1:]{
		path, _ = filepath.Abs(path)	
		readFiles(path,pf)
	}
	pf.Close()

	cmd := exec.Command(filepath.Join(currentpath,"bin","notepad++"))
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
	os.Remove(tmpfile)	
}

func readFiles(path string,pf *os.File){
	filepath.Walk(path,func(path string, info os.FileInfo, err error) error{		
		if !info.IsDir(){
			pf.WriteString(path+"\n")
		}
		return nil
	})
}
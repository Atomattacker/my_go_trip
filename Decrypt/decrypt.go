package main

import(
	"os"
    "fmt"
	"path/filepath"
	"os/exec"
	"bufio"
	"strings"
	"encoding/xml"
	"io/ioutil"
)

const file ="temppaths"

/*
<?xml version="1.0" encoding="utf-8"?>
<configs>
	<goroutineconfig>
		<goroutinecount>8</goroutinecount>
	</goroutineconfig>
</configs>
*/
type GoroutineConfig struct{
	GoroutineCount int `xml:"goroutineconfig>goroutinecount"`
}

var allowsExts map[string]bool

func main() {
	currentpath,_ :=filepath.Abs(filepath.Dir(os.Args[0]))
	tmpfile := filepath.Join(currentpath,"bin",file)
	allowsExts = make(map[string]bool)

	if f,err := os.Open(filepath.Join(currentpath,"allow")); err==nil{
		scanner := bufio.NewScanner(f)
		for scanner.Scan(){
			t :=strings.TrimSpace(scanner.Text())
			allowsExts[t] = true	
		}
	}

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


	var cfg GoroutineConfig
	xmlFile :=filepath.Join(currentpath,"app.config")
	if bytes, err :=ioutil.ReadFile(xmlFile);err ==nil{
		xml.Unmarshal(bytes,&cfg)
	}	

	arg := fmt.Sprint(cfg.GoroutineCount)
	cmd := exec.Command(filepath.Join(currentpath,"bin","notepad++"),arg)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.Run()
	os.Remove(tmpfile)	
}

func isAllow(ext string) bool{
	_,ok := allowsExts[ext]
	return ok
}

func readFiles(path string,pf *os.File){
	filepath.Walk(path,func(path string, info os.FileInfo, err error) error{		
		if !info.IsDir() && isAllow(filepath.Ext(path)){
			pf.WriteString(path+"\n")
		}
		return nil
	})
}
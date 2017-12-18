package main

import(
	"fmt"
	"net/rpc"
	"os/exec"
	"os"
	"bufio"
	"io/ioutil"
	"path/filepath"
	"strconv"
)

type Args struct{
	OldPath,NewPath string
}

type Reply struct{
	Err error
}

const (
	ext = ".mybak"
	tmpext =".mybak.txt"
	pathfile = "temppaths"
	routinecount = 10
)

var currentDir string

var client *rpc.Client
func main() {
	var routinecount = routinecount
	if len(os.Args) >= 2 {
		 if i,err := strconv.Atoi(os.Args[1]); err ==nil{
			routinecount = i
		 }
	}
	currentDir, _ = filepath.Abs(filepath.Dir(os.Args[0]))

	if startProcess() && dialHTTP() {
		if f, err :=openPathfile(); err == nil{
			defer f.Close()
			scanner := bufio.NewScanner(f)
			ch := make(chan struct{})
			i :=0
			for index := 0; index < routinecount; index++ {				
				if scanner.Scan(){
					i++
					t := scanner.Text()
					go func(){
						convertAndRename(t)
						ch<-struct{}{}
					} ()
				} else {
					break
				}
			}

			for {
				if scanner.Scan(){					
					<-ch
					t :=scanner.Text()
					go func(){
						convertAndRename(t)
						ch<-struct{}{}
					} ()
				} else {
					i--
					break
				}				
			}

			for i>=0{
				<-ch
				i--
			}
			
		}else{
			fmt.Println(err)
		}
		cmd.Process.Kill()
	}
}

var cmd *exec.Cmd

func startProcess()bool{	
	filePath := filepath.Join(currentDir,"rename.exe")
	cmd = exec.Command(filePath)	
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Start()
	if err != nil{
		fmt.Println(err)
		return false
	}
	return true
}

func dialHTTP() bool{
	var err error
	client, err = rpc.DialHTTP("tcp","localhost:12306")
	if err != nil{
		fmt.Println(err)
		return false
	}
	return true
}

func openPathfile() (f *os.File,err error){
	if _, err = os.Stat(filepath.Join(currentDir,pathfile)); os.IsNotExist(err){
		return
	}
	
	f,err = os.Open(filepath.Join(currentDir,pathfile))
	return 
}

func convertAndRename(path string){
	tmpPath := path + tmpext
	if !rename(path,tmpPath) {
		return
	}
	content, err := ioutil.ReadFile(tmpPath)
	if err != nil{
		fmt.Println(err)
		return
	}
	oldPath := path+ext
	f, ce := os.Create(oldPath)
	if ce !=nil{
		fmt.Println(ce)	
		return
	}	
	if _, we := f.Write(content);we != nil{
		fmt.Println(we)	
		return
	}	
	f.Close()
	os.Remove(tmpPath)
	rename(oldPath,path)
	fmt.Println(path)
}

func rename(oldpath,newpath string)bool{
	arg := Args{oldpath,newpath}
	var reply Reply
	err := client.Call("Arith.Rename",arg,&reply)
	if err !=nil {
		fmt.Println(err)
		return false
	}
	if reply.Err != nil{
		fmt.Println(err)
		return false
	}
	return true
}
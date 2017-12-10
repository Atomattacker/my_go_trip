package main

import(
	"fmt"
	"io/ioutil"
	"net/rpc"
	"os/exec"
	"os"
	"bufio"
)

type Args struct{
	OldPath,NewPath string
}

type Reply struct{
	Err error
}

const (
	pathfile = "temppaths"
	routinecount = 10
)

var client *rpc.Client

func main() {

	if startprocess() && dialHTTP() && parsefile() {
		
	}
	
}

func startprocess()bool{
	rename := exec.Command("raname")	
	err := rename.Start()
	if err != nil{
		fmt.Println(err)
		return false
	}
	return true
}

func dialHTTP() bool{
	var err error
	client, err = rpc.DialHTTP("tcp","localhost:12345")
	if err != nil{
		fmt.Println(err)
		return false
	}
	return true
}

func parsefile() bool{
	if _, err := os.Stat(pathfile); os.IsNotExist(err){
		fmt.Println("路径文件不存在")
		return false
	}
	f,err := os.Open(pathfile)
	if err != nil{
		fmt.Println(err)
		return false
	}
	defer f.Close()
	scanner := bufio.NewReader(f)
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
	}
	return true
}
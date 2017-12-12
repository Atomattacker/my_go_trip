package main

import(	
	"fmt"	
    "os"
	"net/rpc"
	"net/http"
)

type Args struct{
	OldPath, NewPath string
}

type Reply struct{
	Err error
}

type Arith int

func (ar *Arith) Rename(arg *Args, reply *Reply) error{
	reply.Err = os.Rename(arg.OldPath,arg.NewPath)
	return reply.Err
}


func main() {
	arith := new(Arith)
	rpc.Register(arith)
	rpc.HandleHTTP()
	err := http.ListenAndServe("localhost:12306",nil)
	if err !=nil{
		fmt.Println(err.Error())
	}
}

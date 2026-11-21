package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "github.com/farhatcodes/apiresilience-lab/internal/load"
)

func main(){
    url:=flag.String("url","http://localhost:8080","authorized target URL")
    n:=flag.Int("n",10,"total requests")
    c:=flag.Int("c",2,"concurrent workers")
    flag.Parse()
    r,err:=load.Run(context.Background(),*url,*n,*c,nil)
    if err!=nil {log.Fatal(err)}
    fmt.Printf("total=%d succeeded=%d failed=%d p50=%s p95=%s p99=%s\n",r.Requests,r.Succeeded,r.Failed,r.P50,r.P95,r.P99)
}

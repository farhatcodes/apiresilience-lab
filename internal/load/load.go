package load

import (
 "context"
 "errors"
 "net/http"
 "sort"
 "sync"
 "time"
)

type Result struct{
 Requests int
 Succeeded int
 Failed int
 P50 time.Duration
 P95 time.Duration
 P99 time.Duration
}
func Run(ctx context.Context,url string,requests,concurrency int,client *http.Client)(Result,error){
 if requests<1||concurrency<1{return Result{},errors.New("requests and concurrency must be positive")}
 if client==nil{client=&http.Client{Timeout:5*time.Second}}
 jobs:=make(chan struct{})
 type outcome struct{duration time.Duration;ok bool}
 results:=make(chan outcome,requests)
 var wg sync.WaitGroup
 for i:=0;i<concurrency;i++{
  wg.Add(1)
  go func(){
   defer wg.Done()
   for range jobs {
    start:=time.Now()
    req,err:=http.NewRequestWithContext(ctx,http.MethodGet,url,nil)
    if err!=nil {results<-outcome{time.Since(start),false};continue}
    resp,err:=client.Do(req)
    if err!=nil {results<-outcome{time.Since(start),false};continue}
    resp.Body.Close()
    results<-outcome{time.Since(start),resp.StatusCode>=200&&resp.StatusCode<400}
   }
  }()
 }
 for i:=0;i<requests;i++{jobs<-struct{}{}}
 close(jobs)
 wg.Wait()
 close(results)
 out:=Result{}
 durations:=[]time.Duration{}
 for result:=range results{
  out.Requests++
  durations=append(durations,result.duration)
  if result.ok{out.Succeeded++}else{out.Failed++}
 }
 sort.Slice(durations,func(i,j int)bool{return durations[i]<durations[j]})
 percentile:=func(p float64)time.Duration{
  index:=int(float64(len(durations)-1)*p+0.5)
  return durations[index]
 }
 out.P50,out.P95,out.P99=percentile(.5),percentile(.95),percentile(.99)
 return out,nil
}

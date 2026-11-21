package load

import (
 "context"
 "net/http"
 "net/http/httptest"
 "testing"
)
func TestRun(t *testing.T){
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){w.WriteHeader(200)}))
 defer server.Close()
 got,err:=Run(context.Background(),server.URL,12,3,nil)
 if err!=nil||got.Requests!=12||got.Succeeded!=12||got.Failed!=0 {t.Fatalf("got=%+v err=%v",got,err)}
}
func TestBadInputs(t *testing.T){
 _,err:=Run(context.Background(),"http://localhost",0,1,nil)
 if err==nil {t.Fatal("expected error")}
}

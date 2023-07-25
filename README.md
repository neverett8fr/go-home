# go-home

_a thing to allow automation of http requests (maybe some of the home automation standards in future)_  
  
This is a library to allow for __home__ automation, etc.  
  
```go
func main() {
 home, _ := NewHome()

 home.RegisterEndpoint("endpoint1", "192.168....", "GET")
 home.AddCondition("endpoint1", func() bool {
  return true
 })

 go home.StartHandlers()

 time.Sleep(5 * time.Minute)
 home.StopHandlers()

}

```

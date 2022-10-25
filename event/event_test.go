package event

import (
    "fmt"
    "testing"

    "github.com/mitchellh/mapstructure"
)

//  for example
type Info struct{}

func (i Info) Parse() {
    fmt.Println("catch info log")
}

func LogLoadWrapper() LoadWrapper {
    return func(d interface{}, v map[string]interface{}) {
        if err := mapstructure.Decode(v, &d); err == nil {
            switch d.(type) {
            case Info:
                d.(Info).Parse()
            default:
                fmt.Println("sorry we did not support this type!")
            }
        }
    }
}

func TestEvent(t *testing.T) {
    c := NewContract("xxx", "xxx", "xxxx")
    c.GetLogs(LogLoadWrapper(), Info{})
}

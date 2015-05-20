package utils

import "fmt"

// A fake response object that will be used to unmarshal test data
type TestResponse1 struct {
	Name string `json:"name" xml:"name"`
}

type TestResponse2 struct {
	Age int `json:"age"`
}

type TestResponseArray []TestResponse1

func (tr TestResponse1) String() string {
	return fmt.Sprintf("TestResponse1{ Name=%s }", tr.Name)
}

func (tr TestResponse2) String() string {
	return fmt.Sprintf("TestResponse2{ Age=%d }", tr.Age)
}

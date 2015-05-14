package utils

import "fmt"

// A fake response object that will be used to unmarshal test data
type TestResponse struct {
	Name string `json:"name" xml:"name"`
}

func (tr TestResponse) String() string {
	return fmt.Sprintf("TestResponse{ Name=%s }", tr.Name)
}

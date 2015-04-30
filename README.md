# GoRest
An immutable rest client for Go

Example:

    package main
    import (
        "github.com/Falkenfighter/GoRest"
        "fmt"
    )
    
    type TwitterTerms struct {
        Terms string `json:"tos"`
    }
    
    type TwitterErrors struct {
        Errors []TwitterError `json:"errors"`
    }
    
    type TwitterError struct {
        Code    int `json:"code"`
        Message string `json:"message"`
    }
    
    func (te TwitterError) String() string {
        return fmt.Sprintf("TwitterError{ Code=%v, Message=%s }", te.Code, te.Message)
    }
    
    func main() {
        baseClient := GoRest.MakeClient("https://api.twitter.com/").Path("1.1")
    
        // Create Help client
        helpClient := baseClient.Path("help")
    
        // Create the pointer to our struct
        twitterTOS := new(TwitterTerms)
        twitterErrors := new(TwitterErrors)
    
        // Finish the path and make a GET call on it.
        // Note: JSON is the default accept type so we do not need to specify it
        err := helpClient.Path("tos.json").Get(twitterTOS, twitterErrors)
        if err != nil {
            fmt.Println("Client Error:", err)
        }
    
        if len(twitterErrors.Errors) != 0 {
            fmt.Println("Twitter Error:", twitterErrors.Errors[0])
        }
    
        fmt.Println(twitterTOS.Terms)
    }
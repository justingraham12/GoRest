package main
import (
    "github.com/Falkenfighter/GoRest"
    "fmt"
    "encoding/json"
)

// Wrapper for events
type Events []Event

// The Base response object from /events
type Event struct {
    Type        string          `json:"type"`
    Public      bool            `json:"public"`
    Payload     interface{}     `json:"payload"`
    Repo        Repository      `json:"repo"`
    Actor       User            `json:"actor"`
    Org         Organization    `json:"org"`
    CreatedAt   string          `json:"created_at"`
    Id          string          `json:"id"`
}
type Repository struct {
    Id      int     `json:"id"`
    Name    string  `json:"name"`
    Url     string  `json:"url"`
}

type User struct {
    Id          int     `json:"id"`
    Login       string  `json:"login"`
    GravatarId  string  `json:"gravatar_id"`
    AvatarUrl   string  `json:"avatar_url"`
    Url         string  `json:"url"`
}
type Organization struct {
    Id          int     `json:"id"`
    Login       string  `json:"login"`
    GravatarId  string  `json:"gravatar_id"`
    AvatarUrl   string  `json:"avatar_url"`
    Url         string  `json:"url"`
}

// A new RestClient pointed at the base url. The Accept type has also been set as JSON
var Client GoRest.RestClient = GoRest.MakeClient("https://api.github.com/").Accept(GoRest.ApplicationJSON)

func main() {
    // The events object to unmarshal the JSON response into
    allEvents := new(Events)

    // Make the request and validate no error resulted.
    if err := GetEvents(allEvents); err != nil { fmt.Println("Error getting events:", err); return }

    // At this point the events struct has been populated with the response JSON
    fmt.Println("First Event:", (*allEvents)[0], "\n\n")

    // Convert the struct into pretty printed JSON and log it out (used as an example to display the response)
    str, err := json.MarshalIndent(allEvents, "", "\t")
    if err != nil { fmt.Println(err); return }
    fmt.Println("All Events:\n", string(str), "\n\n\n\n")


    // The events for a particular user
    usersEvents := new(Events)

    // Get my events
    if err = GetUsersEvents("Falkenfighter", usersEvents);
    err != nil { fmt.Println("Error getting Users events:",err); return }

    // Convert the struct into pretty printed JSON and log it out (used as an example to display the response)
    str, err = json.MarshalIndent(usersEvents, "", "\t")
    if err != nil { fmt.Println(err); return }
    fmt.Println("Users Events:\n", string(str))
}

// GET /events
func GetEvents(events *Events) (err error) {
    _, err = Client.Path("events").Get(events); return
}

// GET /users/:username/events
func GetUsersEvents(username string, events *Events) (err error) {
    _, err = Client.Path("users", username, "events").Get(events); return
}
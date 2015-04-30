package GoRest
import (
    "encoding/json"
    "encoding/xml"
    "errors"
)

type MediaType interface {
    String() string
    Unmarshal(body []byte, entity interface{}) error
}

type XML        string
type JSON       string
type TEXT_PLAIN string
type TEXT_XML   string
type TEXT_HTML  string

const (
    ApplicationXML  XML         = "application/xml"
    ApplicationJSON JSON        = "application/json"
    TextPlain       TEXT_PLAIN  = "text/plain"
    TextXML         TEXT_XML    = "text/xml"
    TextHTML        TEXT_HTML   = "text/html"
)

// JSON

func (j JSON) String() string {
    return string(j)
}

func (j JSON) Unmarshal(body []byte, entity interface{}) error {
    err := json.Unmarshal(body, entity)
    if err != nil { return err }
    return nil
}

// XML

func (x XML) String() string {
    return string(x)
}

func (x XML) Unmarshal(body []byte, entity interface{}) error {
    err := xml.Unmarshal(body, entity)
    if err != nil { return err }
    return nil
}

// TEXT_PLAIN

func (tp TEXT_PLAIN) String() string {
    return string(tp)
}

func (tp TEXT_PLAIN) Unmarshal(body []byte, entity interface{}) error {
    return errors.New("Unable to unmarshal MediaType{TEXT_PLAIN}")
}

// TEXT_XML

func (tx TEXT_XML) String() string {
    return string(tx)
}

func (tx TEXT_XML) Unmarshal(body []byte, entity interface{}) error {
    err := xml.Unmarshal(body, entity)
    if err != nil { return err }
    return nil
}

// TEXT_HTML

func (tx TEXT_HTML) String() string {
    return string(tx)
}

func (tx TEXT_HTML) Unmarshal(body []byte, entity interface{}) error {
    // TODO: Should we be using an xml parser on HTML?
    err := xml.Unmarshal(body, entity)
    if err != nil { return err }
    return nil
}
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

type XML string
type JSON string
type URL_ENCODED string
type TEXT_PLAIN string
type TEXT_XML string
type NO_CONTENT string

const (
	ApplicationXML        XML         = "application/xml"
	ApplicationJSON       JSON        = "application/json"
	ApplicationURLEncoded URL_ENCODED = "application/x-www-form-urlencoded"
	TextPlain             TEXT_PLAIN  = "text/plain"
	TextXML               TEXT_XML    = "text/xml"
	NoContent			  NO_CONTENT  = ""
)

// JSON

func (j JSON) String() string {
	return string(j)
}

func (j JSON) Unmarshal(body []byte, entity interface{}) error {
	err := json.Unmarshal(body, entity)
	if err != nil {
		return err
	}
	return nil
}

// XML

func (x XML) String() string {
	return string(x)
}

func (x XML) Unmarshal(body []byte, entity interface{}) error {
	err := xml.Unmarshal(body, entity)
	if err != nil {
		return err
	}
	return nil
}

// URL_ENCODED

func (url URL_ENCODED) String() string {
	return string(url)
}

func (url URL_ENCODED) Unmarshal(body []byte, entity interface{}) error {
	return errors.New("Unable to unmarshal MediaType{URL_ENCODED}")
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
	if err != nil {
		return err
	}
	return nil
}

// NO_CONTENT

func(nc NO_CONTENT) String() string {
	return string(nc)
}

func (nc NO_CONTENT) Unmarshal(body []byte, entity interface{}) error {
	return nil
}
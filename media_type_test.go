package GoRest

import (
	u "./utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestXMLString(t *testing.T) {
	assert.Equal(t, "application/xml", ApplicationXML.String())
}

func TestXMLUnmarshal(t *testing.T) {
	entity := new(u.TestResponse)
	ApplicationXML.Unmarshal([]byte(`<Response><name>test</name></Response>`), entity)
	assert.Equal(t, "test", entity.Name)
}

func TestJSONString(t *testing.T) {
	assert.Equal(t, "application/json", ApplicationJSON.String())
}

func TestJSONUnmarshal(t *testing.T) {
	entity := new(u.TestResponse)
	ApplicationJSON.Unmarshal([]byte(`{"name":"test"}`), entity)
	assert.Equal(t, "test", entity.Name)
}

func TestURLEncodedString(t *testing.T) {
	assert.Equal(t, "application/x-www-form-urlencoded", ApplicationURLEncoded.String())
}

func TestURLEncodedUnmarshal(t *testing.T) {
	err := ApplicationURLEncoded.Unmarshal([]byte(`Some Text`), new(u.TestResponse))
	assert.NotNil(t, err)
}

func TestTextPlainString(t *testing.T) {
	assert.Equal(t, "text/plain", TextPlain.String())
}

func TestTextPlainUnmarshal(t *testing.T) {
	err := TextPlain.Unmarshal([]byte(`Some Text`), new(u.TestResponse))
	assert.NotNil(t, err)
}

func TestTextXMLString(t *testing.T) {
	assert.Equal(t, "text/xml", TextXML.String())
}

func TestTextXMLUnmarshal(t *testing.T) {
	entity := new(u.TestResponse)
	TextXML.Unmarshal([]byte(`<Response><name>test</name></Response>`), entity)
	assert.Equal(t, "test", entity.Name)
}

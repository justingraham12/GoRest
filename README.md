[![Build Status](https://travis-ci.org/Falkenfighter/GoRest.svg?branch=master)](https://travis-ci.org/Falkenfighter/GoRest)  

# GoRest
An immutable rest client for Go

GoRest provides a simple to use wrapper to Golang's net/http package. 

# Feature set:

1) __Fluent URL Path Building__

    GoRest.MakeClient("base_url").Path("param_1", "param_2").Query("key", "value")  
This would build the url: `base_url/param_1/param_2?key=value`  

2) __Response Unmarshalling__

    responseStruct = new(ResponseStruct)
    res, err := GoRest.MakeClient("base_url").Get(responseStruct)  
This takes the response body and unmarshals it into the provided struct. JSON is set as the default unmarshaller but 
can be changed by providing a new MediaType.

    GoRest.MakeClient("base_url").Accept(GoRest.ApplicationXML)
GoRest provides both XML and JSON unmarshalling. If you would like to customize how GoRest unmarshals or add a new type 
you can define a new MediaType. For example if I wanted to unmarshal YAML it would look like this:

    type YAML string
    var ApplicationYAML = "application/yaml"
    
    func (y YAML) String() string {
        return string(y)
    }
    
    func (y YAML) Unmarshal(body []byte, entity interface{}) error {
        // Unmarshal logic goes here converting the body variable into the entity variable
    }  
    
Then to use the new type you provide it to the Accept() function

    GoRest.MakeClient("base_url").Accept(ApplicationYAML)
    
__NOTE:__ GoRest validates the Response "Content-Type" contains the Accept type. So in our example the "Content-Type"
 MUST contain "application/yaml"

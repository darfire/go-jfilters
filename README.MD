# json-filters, a lightweight JSON handling library

**json-filters** is an alternative to the usual Go JSON handling pattern, in which you have to declare Go structures that match your JSON structure. json-filters allows you to access and validate your JSON from a simple *interface{}* object that has been prefilled with the parsed JSON object.


A simple example:
```go
package main

import (
  "encoding/json"
  jfilters "github.com/darfire/go-jfilters"
  "fmt"
)

const S1 = `{
  "z": 123,
  "a": {
    "b": 11,
    "c": ["a", 2, false]
  },
  "d": [0, 1, 2, 3],
  "e": {
    "f": [12, {
      "g": 3.1415
    }]
  }
}`


func main() {
  // parse the string to an agnostic object
	var obj interface{}
	json.Unmarshal([]byte(S1), &obj)

  // create a new filter
	f := jfilters.New()
  // define it; by walking the structure
	fg := f.Key("e").Key("f").Index(1).Key("g")
  // bind it to our parsed object
  b := gb.Bind(obj)
  // v contains the value, err any error that occurs during validation
  v, err := b.GetFloat()
  fmt.Println(v)
}
```

### TODO
- [ ] Test errors
- [ ] Test map and array
- [ ] Document it
- [ ] Add support for user-supplied callbacks and filters
- [ ] Add support for composing filters

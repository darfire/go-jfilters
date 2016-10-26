package json_filters

import (
	"encoding/json"
	"testing"
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

func TestFilters(t *testing.T) {
	var obj map[string]interface{}
	json.Unmarshal([]byte(S1), &obj)

	f := New()
	fa := f.Key("a")
	fz := f.Key("z")
	fb := fa.Key("b")
	fc := fa.Key("c").Index(0)

	fg := f.Key("e").Key("f").Index(1).Key("g")

	if v, err := fz.Bind(obj).GetInt(); err != nil || v != 123 {
		t.Errorf("Expected 123, got err = %v, v = %v", err, v)
	}

	if v, err := fb.Bind(obj).GetInt(); err != nil || v != 11 {
		t.Errorf("Expected 11, got err = %v, v = %v", err, v)
	}

	if v, err := fc.Bind(obj).GetString(); err != nil || v != "a" {
		t.Errorf("Expected \"a\", got err = %v, v = %v", err, v)
	}

	if v, err := fg.Bind(obj).GetFloat(); err != nil || v != 3.1415 {
		t.Errorf("Expected \"a\", got err = %v, v = %v", err, v)
	}
}

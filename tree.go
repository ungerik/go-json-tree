package json

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

func FromBytes(jsonBytes []byte) (result Tree, err error) {
	result.Data = new(interface{})
	err = json.Unmarshal(jsonBytes, result.Data)
	return result, err
}

func FromString(jsonString string) (Tree, error) {
	return FromBytes([]byte(jsonString))
}

func FromFile(filename string) (Tree, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return Tree{}, err
	}
	return FromBytes(data)
}

func FromURL(url string) (Tree, error) {
	if strings.Index(url, "file://") == 0 {
		return FromFile(url[len("file://"):])
	}
	response, err := http.Get(url)
	if err != nil {
		return Tree{}, err
	}
	defer response.Body.Close()
	return FromReader(response.Body)
}

func FromReader(reader io.Reader) (Tree, error) {
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		return Tree{}, err
	}
	return FromBytes(data)
}

// Tree wraps unmarshalled Tree data with a nice interface.
type Tree struct {
	Data *interface{}
}

// String returns a tab indented string representation
func (self Tree) String() string {
	var buf bytes.Buffer
	err := json.Indent(&buf, self.Bytes(), "", "\t")
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Bytes returns a compact UTF-8 string representation
func (self Tree) Bytes() []byte {
	bytes, err := json.Marshal(self.Data)
	if err != nil {
		panic(err)
	}
	return bytes
}

func (self Tree) IsString() bool {
	if self.Data == nil || (*self.Data) == nil {
		return false
	}
	_, ok := (*self.Data).(string)
	return ok
}

func (self Tree) IsNumber() bool {
	if self.Data == nil || (*self.Data) == nil {
		return false
	}
	_, ok := (*self.Data).(float64)
	return ok
}

func (self Tree) IsObject() bool {
	if self.Data == nil || (*self.Data) == nil {
		return false
	}
	_, ok := (*self.Data).(map[string]interface{})
	return ok
}

func (self Tree) IsArray() bool {
	if self.Data == nil || (*self.Data) == nil {
		return false
	}
	_, ok := (*self.Data).([]interface{})
	return ok
}

func (self Tree) IsBool() bool {
	if self.Data == nil || (*self.Data) == nil {
		return false
	}
	_, ok := (*self.Data).(bool)
	return ok
}

func (self Tree) IsNull() bool {
	if self.Data == nil {
		return false
	}
	return *self.Data == nil
}

func (self Tree) IsValid() bool {
	return self.Data != nil
}

func (self Tree) Has(selector string) bool {
	return self.Select(selector).IsValid()
}

func (self Tree) Select(selector string) Tree {
	if self.Data == nil || (*self.Data) == nil {
		return self
	}
	result := self
	for _, s := range strings.Split(selector, ".") {
		switch data := (*result.Data).(type) {
		case map[string]interface{}:
			if value, ok := data[s]; ok {
				// Note that value is a copy and taking its address
				// yields another address than that of the original in the map.
				// But this doesn't matter as long as we are reading only.
				// Also for slices and maps the value copied is a reference.
				result = Tree{&value}
			}

		case []interface{}:
			arrayIndex, err := strconv.ParseUint(s, 10, 32)
			if err != nil {
				return Tree{}
			}
			result = Tree{&data[arrayIndex]}

		default:
			// If it's not a Tree object or array, we can't select a child
			// and thus have to return an invalid instance
			return Tree{}
		}
	}
	return result
}

func (self Tree) GetString(defaultValue ...string) string {
	if self.Data == nil || *self.Data == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		} else {
			return ""
		}
	}
	return (*self.Data).(string)
}

func (self Tree) GetInt(defaultValue ...int) int {
	if self.Data == nil || *self.Data == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		} else {
			return 0
		}
	}
	return int((*self.Data).(float64))
}

func (self Tree) GetFloat(defaultValue ...float64) float64 {
	if self.Data == nil || *self.Data == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		} else {
			return 0
		}
	}
	return (*self.Data).(float64)
}

func (self Tree) GetBool(defaultValue ...bool) bool {
	if self.Data == nil || *self.Data == nil {
		if len(defaultValue) > 0 {
			return defaultValue[0]
		} else {
			return false
		}
	}
	return (*self.Data).(bool)
}

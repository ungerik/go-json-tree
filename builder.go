package json

import (
	"bytes"
	"encoding/json"
	"reflect"
	"strconv"
)

type state int

const (
	atRoot state = iota
	inArray
	inObject
)

type Builder struct {
	stack     []state
	buffer    bytes.Buffer
	afterName bool
}

func (self *Builder) state() state {
	if len(self.stack) == 0 {
		return atRoot
	}
	return self.stack[len(self.stack)-1]
}

// String returns a tab indented string representation
func (self *Builder) String() string {
	var buf bytes.Buffer
	err := json.Indent(&buf, self.Bytes(), "", "\t")
	if err != nil {
		panic(err)
	}
	return buf.String()
}

// Bytes returns a compact UTF-8 string representation
func (self *Builder) Bytes() []byte {
	var buffer bytes.Buffer
	err := json.Compact(&buffer, self.buffer.Bytes())
	if err != nil {
		panic(err)
	}
	return buffer.Bytes()
}

func (self *Builder) Tree() Tree {
	tree, err := FromBytes(self.Bytes())
	if err != nil {
		panic(err)
	}
	return tree
}

func (self *Builder) BeginObject() *Builder {
	if self.state() == inObject {
		panic("Can't BeginObject() in an object, need Name() first")
	}
	self.buffer.WriteByte('{')
	self.stack = append(self.stack, inObject)
	self.afterName = false
	return self
}

func (self *Builder) EndObject() *Builder {
	if self.state() != inObject || self.afterName {
		panic("EndObject() called not at the end of an object")
	}
	self.buffer.WriteString("},")
	self.stack = self.stack[:len(self.stack)-1]
	return self
}

func (self *Builder) BeginArray() *Builder {
	if self.state() == inObject {
		panic("Can't BeginArray() in an object, need Name() first")
	}
	self.buffer.WriteByte('[')
	self.stack = append(self.stack, inObject)
	self.afterName = false
	return self
}

func (self *Builder) EndArray() *Builder {
	if self.state() != inArray || self.afterName {
		panic("EndArray() called not at the end of an array")
	}
	self.buffer.WriteString("],")
	self.stack = self.stack[:len(self.stack)-1]
	return self
}

func (self *Builder) Name(name string) *Builder {
	if self.state() != inObject {
		panic("Name() must be called in an object")
	}
	self.buffer.WriteByte('"')
	self.buffer.WriteString(name)
	self.buffer.WriteString(`":`)
	self.afterName = true
	return self
}

func (self *Builder) Value(value interface{}) *Builder {
	if value == nil {
		self.buffer.WriteString("null")
	}
	v := reflect.ValueOf(value)
	switch v.Kind() {
	case reflect.String:
		self.buffer.WriteByte('"')
		json.HTMLEscape(&self.buffer, []byte(v.String()))
		self.buffer.WriteByte('"')

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		self.buffer.WriteString(strconv.FormatInt(v.Int(), 10))

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		self.buffer.WriteString(strconv.FormatUint(v.Uint(), 10))

	case reflect.Float32, reflect.Float64:
		self.buffer.WriteString(strconv.FormatFloat(v.Float(), 'f', -1, 64))

	case reflect.Bool:
		if v.Bool() {
			self.buffer.WriteString("true")
		} else {
			self.buffer.WriteString("false")
		}

	default:
		panic("Type not supported as JSON value")
	}

	self.buffer.WriteByte(',')
	self.afterName = false
	return self
}

// Package codec 演示 encoding/json 与 encoding/xml 的序列化和反序列化。
//
// import 路径：just-go/stage-1-syntax/06-stdlib-essentials/codec
package codec

import (
	"encoding/json"
	"encoding/xml"
)

// Lesson 是用于 JSON/XML round trip 的示例结构体。
type Lesson struct {
	XMLName xml.Name `json:"-" xml:"lesson"`
	ID      int      `json:"id" xml:"id,attr"`
	Title   string   `json:"title" xml:"title"`
	Done    bool     `json:"done" xml:"done"`
}

// JSONRoundTrip 将结构体编码为 JSON 再解码回来。
func JSONRoundTrip(lesson Lesson) (string, Lesson, error) {
	data, err := json.Marshal(lesson)
	if err != nil {
		return "", Lesson{}, err
	}
	var decoded Lesson
	if err := json.Unmarshal(data, &decoded); err != nil {
		return "", Lesson{}, err
	}
	return string(data), decoded, nil
}

// XMLRoundTrip 将结构体编码为 XML 再解码回来。
func XMLRoundTrip(lesson Lesson) (string, Lesson, error) {
	data, err := xml.Marshal(lesson)
	if err != nil {
		return "", Lesson{}, err
	}
	var decoded Lesson
	if err := xml.Unmarshal(data, &decoded); err != nil {
		return "", Lesson{}, err
	}
	return string(data), decoded, nil
}

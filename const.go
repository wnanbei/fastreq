package fastreq

import jsoniter "github.com/json-iterator/go"

var jsonMarshal = jsoniter.ConfigCompatibleWithStandardLibrary.Marshal
var jsonUnmarshal = jsoniter.ConfigCompatibleWithStandardLibrary.Unmarshal

// SetJsonMarshal can set json marshal function.
// Default Marshal is github.com/json-iterator/go
func SetJsonMarshal(f func(any) ([]byte, error)) {
	jsonMarshal = f
}

// SetJsonUnmarshal can set json unmarshal function.
// Default Unarshal is github.com/json-iterator/go
func SetJsonUnmarshal(f func([]byte, any) error) {
	jsonUnmarshal = f
}

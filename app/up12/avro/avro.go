// Copyright 2023 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package avro provides API for handling avro
package avro

import (
	"log"
	"os"

	"github.com/linkedin/goavro/v2"
)

// NewCodedecFromFile creates codedec from avro schema file
func NewCodedecFromFile(path string) *goavro.Codec {
	schema, err := os.ReadFile(path)
	if err != nil {
		log.Panicln(err)
	}
	codec, err := goavro.NewCodec(string(schema))
	if err != nil {
		log.Panicln(err)
	}
	return codec
}

// EncodeToJSON encodes data to JSON using given avro codedec
func EncodeToJSON(codedec *goavro.Codec, data map[string]interface{}) ([]byte, error) {
	json, err := codedec.TextualFromNative(nil, data)
	if err != nil {
		log.Println("fail to encode data=", data, "err=", err)
	}
	return json, err
}

// DecodeFromJSON decodes json using given avro codedec
func DecodeFromJSON(codedec *goavro.Codec, json []byte) (map[string]interface{}, error) {
	native, _, err := codedec.NativeFromTextual(json)
	if err != nil {
		log.Println("fail to decode json=", json, "err=", err)
		return nil, err
	}
	data, ok := native.(map[string]interface{})
	if !ok {
		log.Println("fail to decode json=", json, "err=", err)
		return nil, err
	}
	return data, err

}

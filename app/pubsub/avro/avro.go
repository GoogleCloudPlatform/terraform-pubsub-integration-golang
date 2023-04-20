// Package avro provide utility of avro
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
		log.Println("Fail to encode data=", data, "err=", err)
	}
	return json, err
}

// DecodeFromJSON decodes json using given avro codedec
func DecodeFromJSON(codedec *goavro.Codec, json []byte) (map[string]interface{}, error) {
	native, _, err := codedec.NativeFromTextual([]byte(json))
	if err != nil {
		log.Println("Fail to decode json=", json, "err=", err)
		return nil, err
	}
	data, ok := native.(map[string]interface{})
	if !ok {
		log.Println("Fail to decode json=", json, "err=", err)
		return nil, err
	}
	return data, err

}

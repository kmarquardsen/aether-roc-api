// Code generated by oapi-codegen. DO NOT EDIT.
// Package server provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.7.0 DO NOT EDIT.
package server

import (
	"bytes"
	"compress/gzip"
	"encoding/base64"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

// Base64 encoded, gzipped, json marshaled Swagger object
var swaggerSpec = []string{

	"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",
	"dKvDWLnDLXp10np/f/+/AQAA//91mU1ItRcBAA==",
}

// GetSwagger returns the content of the embedded swagger specification file
// or error if failed to decode
func decodeSpec() ([]byte, error) {
	zipped, err := base64.StdEncoding.DecodeString(strings.Join(swaggerSpec, ""))
	if err != nil {
		return nil, fmt.Errorf("error base64 decoding spec: %s", err)
	}
	zr, err := gzip.NewReader(bytes.NewReader(zipped))
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}
	var buf bytes.Buffer
	_, err = buf.ReadFrom(zr)
	if err != nil {
		return nil, fmt.Errorf("error decompressing spec: %s", err)
	}

	return buf.Bytes(), nil
}

var rawSpec = decodeSpecCached()

// a naive cached of a decoded swagger spec
func decodeSpecCached() func() ([]byte, error) {
	data, err := decodeSpec()
	return func() ([]byte, error) {
		return data, err
	}
}

// Constructs a synthetic filesystem for resolving external references when loading openapi specifications.
func PathToRawSpec(pathToFile string) map[string]func() ([]byte, error) {
	var res = make(map[string]func() ([]byte, error))
	if len(pathToFile) > 0 {
		res[pathToFile] = rawSpec
	}

	return res
}

// GetSwagger returns the Swagger specification corresponding to the generated code
// in this file. The external references of Swagger specification are resolved.
// The logic of resolving external references is tightly connected to "import-mapping" feature.
// Externally referenced files must be embedded in the corresponding golang packages.
// Urls can be supported but this task was out of the scope.
func GetSwagger() (swagger *openapi3.T, err error) {
	var resolvePath = PathToRawSpec("")

	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true
	loader.ReadFromURIFunc = func(loader *openapi3.Loader, url *url.URL) ([]byte, error) {
		var pathToFile = url.String()
		pathToFile = path.Clean(pathToFile)
		getSpec, ok := resolvePath[pathToFile]
		if !ok {
			err1 := fmt.Errorf("path not found: %s", pathToFile)
			return nil, err1
		}
		return getSpec()
	}
	var specData []byte
	specData, err = rawSpec()
	if err != nil {
		return
	}
	swagger, err = loader.LoadFromData(specData)
	if err != nil {
		return
	}
	return
}
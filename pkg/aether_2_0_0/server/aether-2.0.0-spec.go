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
	"mmVN+fE/AQAA//91OGZdKkIBAA==",
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

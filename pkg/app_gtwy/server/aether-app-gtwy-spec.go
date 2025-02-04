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

	"H4sIAAAAAAAC/+xXUW/bNhD+KwSzhxaQLSfrgMLDgHltZhgL2gHrnpogYKSTfI1EsuTJjmDovw9Hy7Ys",
	"K0gGNNtLnyyRx+/ueN99Om9kYkprNGjycrqRVjlVAoELb6RcDsRPidEEOjwSPFBsC4X6Z5EslfNAv1SU",
	"jd7ypk+WUKpgVluQU+nJoc5l0zSRTMEnDi2h0XLagotXKawwAYFaGG38KDE6w/y1jCSylVW0lJHUqoT9",
	"GRlJB18rdJDKKbkKGHzrOYQ9s/Z2Tuv69n2ADnk5Y8ERQjBQRCpZ8ulNL6j9TiRL9XAFOqelnL6dRLJE",
	"vXs9jyQhFXBs30u4l27fUWF0LjpLIkMoem7PJxdvHnPcRRv0zZmPcCDFK9T3gozYmsgu5O7Mo4C3uTOV",
	"9aegn5YgCvQkTNYCi9Y0kkhQ+gFKnHjen2gNlXOqlpGsNH6tYLGFCfWOZIreFqoebYnRD6fdFbzLuVY+",
	"EGz+90IYJ95dLZ5b3yM3A9eCJeCp+7D6TA+t7SmyHcC1z0a1Q5gey1tMEnxR4jfd7vwsQRM469AHZt3s",
	"zc3dF0iIo+q16yPsImNFASsohLJW5IpgreqWa16wPinU4GTU6/X0MdCrlq4DcF3S/uAgk1N5Fh+EMm61",
	"Ju4LTYfSzNb+9tPEjuTD6B5qzxd3aMebpgM8s3bEoKP3+1h7N8rWqDOz022VBN2GUmHB1NCZ+dVY0Bpo",
	"bdw96nxsXH4Q2Y8WtPiw3xS/m0qnqlWayjHGksj6aRwPwJzo/NmZ+E0l99zcOr3W15qrOQNaghMzawtM",
	"AraYtzV4NZu/FrM/F8I6s8IUvFCi0pghpGE5M07cGVqGmjtTCKVTUQI5TDzXkw7wtlCUGVeOW6dzgX4n",
	"T3xKddxnJqk8pJEAre4KzpuBuhYps4+pxYqCpS2gBE0CNUFRYM7Pu5BMtgvhzTz+ac7rGhLCFVItPLhQ",
	"uUisl1iAUHeenEqIfXpTgijMerSlOmpOSyUIIbUANr7mOhSYgPZB+tqyzSw35OhiPDkq0zSO1+v1WIVd",
	"rlDcHvXx1eLd5Ye/LsORLsMerY2M5Aqc39b1fDzZnmMSKItyKn8MS1H4aIfmiZW1Oa3reHUeb7Yf7yY+",
	"KIKPN0fy0MQeiVf5J7x3uredRLgAIahFKqdyDnToAwfeGs6M7S4mk97g0qll/MVvP8qHeeVf9LkfGmfm",
	"l5/Exz8Euw3jSFWWytXtxpM5D6QcSVJ5UAJ500RHk9nn4XAPJnE7KDXRi4xu91CLXhbD89qx+J+MbS8Y",
	"3e46h+Nqd08j4pv+VqSNN3sJb54m8H/I329P36ey/87l/4HLLxjRUWkHYur+lTjpsKb5JwAA//98N/k0",
	"dg4AAA==",
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

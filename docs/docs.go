// Package docs GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag
package docs

import (
	"bytes"
	"encoding/json"
	"strings"
	"text/template"

	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/ohlc/data": {
            "get": {
                "description": "Show OHLC data matching the given criteria",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ohlc"
                ],
                "summary": "Search OHLC data",
                "operationId": "get_ohlc",
                "parameters": [
                    {
                        "type": "integer",
                        "description": "1644719700000",
                        "name": "unix",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "BTCUSDT",
                        "name": "symbol",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "42123.29000000",
                        "name": "open",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "42148.32000000",
                        "name": "high",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "42120.82000000",
                        "name": "low",
                        "in": "query"
                    },
                    {
                        "type": "number",
                        "description": "42146.06000000",
                        "name": "close",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "1",
                        "name": "page",
                        "in": "query"
                    },
                    {
                        "type": "integer",
                        "description": "20",
                        "name": "itemsPerPage",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.OHLCSearchResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            },
            "post": {
                "description": "Upload OHLC data to database",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "ohlc"
                ],
                "summary": "Upload OHLC data",
                "operationId": "post_ohlc",
                "parameters": [
                    {
                        "type": "file",
                        "description": "csv file",
                        "name": "csv",
                        "in": "formData",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/entity.OHLCSaveResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/v1.response"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "entity.OHLC": {
            "type": "object",
            "properties": {
                "close": {
                    "type": "number",
                    "example": 42146.06
                },
                "high": {
                    "type": "number",
                    "example": 42148.32
                },
                "id": {
                    "type": "integer",
                    "example": 1
                },
                "low": {
                    "type": "number",
                    "example": 42120.82
                },
                "open": {
                    "type": "number",
                    "example": 42123.29
                },
                "symbol": {
                    "type": "string",
                    "example": "BTCUSDT"
                },
                "unix": {
                    "type": "integer",
                    "example": 1644719700000
                }
            }
        },
        "entity.OHLCSaveResponse": {
            "type": "object",
            "properties": {
                "status": {
                    "type": "string"
                }
            }
        },
        "entity.OHLCSearchResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/entity.OHLC"
                    }
                },
                "failedRecords": {
                    "type": "integer"
                },
                "limit": {
                    "type": "integer"
                },
                "message": {
                    "type": "string"
                },
                "offset": {
                    "type": "integer"
                },
                "status": {
                    "type": "string"
                },
                "total": {
                    "type": "integer"
                }
            }
        },
        "v1.response": {
            "type": "object",
            "properties": {
                "error": {
                    "type": "string",
                    "example": "message"
                }
            }
        }
    }
}`

type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{
	Version:     "1.0",
	Host:        "localhost:8080",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "OHLC Handler",
	Description: "Use for import/export OHLC data.",
}

type s struct{}

func (s *s) ReadDoc() string {
	sInfo := SwaggerInfo
	sInfo.Description = strings.Replace(sInfo.Description, "\n", "\\n", -1)

	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
		"escape": func(v interface{}) string {
			// escape tabs
			str := strings.Replace(v.(string), "\t", "\\t", -1)
			// replace " with \", and if that results in \\", replace that with \\\"
			str = strings.Replace(str, "\"", "\\\"", -1)
			return strings.Replace(str, "\\\\\"", "\\\\\\\"", -1)
		},
	}).Parse(doc)
	if err != nil {
		return doc
	}

	var tpl bytes.Buffer
	if err := t.Execute(&tpl, sInfo); err != nil {
		return doc
	}

	return tpl.String()
}

func init() {
	swag.Register("swagger", &s{})
}
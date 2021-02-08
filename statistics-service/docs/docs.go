// GENERATED BY THE COMMAND ABOVE; DO NOT EDIT
// This file was generated by swaggo/swag

package docs

import (
	"bytes"
	"encoding/json"
	"strings"

	"github.com/alecthomas/template"
	"github.com/swaggo/swag"
)

var doc = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{.Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Keptn Team",
            "url": "http://www.keptn.sh"
        },
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/event": {
            "post": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "Handle incoming cloud event",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Events"
                ],
                "summary": "Handle event",
                "parameters": [
                    {
                        "description": "Event type",
                        "name": "event",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/operations.Event"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok"
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/operations.Error"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "$ref": "#/definitions/operations.Error"
                        }
                    }
                }
            }
        },
        "/statistics": {
            "get": {
                "security": [
                    {
                        "ApiKeyAuth": []
                    }
                ],
                "description": "get statistics about Keptn installation",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "Statistics"
                ],
                "summary": "Get statistics",
                "parameters": [
                    {
                        "type": "string",
                        "description": "From",
                        "name": "from",
                        "in": "query"
                    },
                    {
                        "type": "string",
                        "description": "To",
                        "name": "to",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "ok",
                        "schema": {
                            "$ref": "#/definitions/operations.Statistics"
                        }
                    },
                    "400": {
                        "description": "Invalid payload",
                        "schema": {
                            "$ref": "#/definitions/operations.Error"
                        }
                    },
                    "500": {
                        "description": "Internal error",
                        "schema": {
                            "$ref": "#/definitions/operations.Error"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "operations.Error": {
            "type": "object",
            "properties": {
                "errorCode": {
                    "type": "string"
                },
                "message": {
                    "type": "string"
                }
            }
        },
        "operations.Event": {
            "type": "object",
            "properties": {
                "contenttype": {
                    "type": "string"
                },
                "data": {
                    "type": "object"
                },
                "extensions": {
                    "type": "object"
                },
                "id": {
                    "type": "string"
                },
                "shkeptncontext": {
                    "type": "string"
                },
                "source": {
                    "type": "string"
                },
                "specversion": {
                    "type": "string"
                },
                "time": {
                    "type": "string"
                },
                "triggeredid": {
                    "type": "string"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "operations.Project": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "services": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/operations.Service"
                    }
                }
            }
        },
        "operations.Service": {
            "type": "object",
            "properties": {
                "events": {
                    "type": "object",
                    "additionalProperties": {
                        "type": "integer"
                    }
                },
                "executedSequences": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "operations.Statistics": {
            "type": "object",
            "properties": {
                "from": {
                    "type": "string"
                },
                "projects": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/operations.Project"
                    }
                },
                "to": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiKeyAuth": {
            "type": "apiKey",
            "name": "x-token",
            "in": "header"
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
	Host:        "",
	BasePath:    "/v1",
	Schemes:     []string{},
	Title:       "Statistics Service API",
	Description: "This is the API documentation of the Statistics Service.",
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
	swag.Register(swag.Name, &s{})
}

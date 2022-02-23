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
        "/availableHours/": {
            "get": {
                "description": "List all the hours remaining for creaiting an entrie on the schedule",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "titulacion de las horas a obtener",
                        "name": "titulacion",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "curso de las horas a obtener",
                        "name": "curso",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "string",
                        "description": "grupo de las horas a obtener",
                        "name": "grupo",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.AvailableHours"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    }
                }
            }
        },
        "/getEntries/": {
            "get": {
                "description": "List all the entries of the  schedule",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "titulacion de las horas a obtener",
                        "name": "degree",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "curso de las horas a obtener",
                        "name": "year",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "grupo de las horas a obtener",
                        "name": "group",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/domain.AvailableHours"
                            }
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    }
                }
            }
        },
        "/getICS/": {
            "get": {
                "description": "Get the schedule in ics format",
                "produces": [
                    "application/json"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "titulacion de las horas a obtener",
                        "name": "degree",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "curso de las horas a obtener",
                        "name": "year",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "grupo de las horas a obtener",
                        "name": "group",
                        "in": "query",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "string"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    }
                }
            }
        },
        "/listDegrees/": {
            "get": {
                "description": "List all degrees' descriptions avaiable, it do not require any parameter",
                "produces": [
                    "application/json"
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.ListDegreesDTO"
                            }
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    }
                }
            }
        },
        "/ping/": {
            "get": {
                "description": "Response \"pong\" if the server is currrently available",
                "produces": [
                    "text/plain"
                ],
                "responses": {
                    "200": {
                        "description": "Returns \"pong\" "
                    }
                }
            }
        },
        "/updateByCSV/": {
            "post": {
                "description": "The request will update the database creating degrees, subjects, years, groups and hours",
                "produces": [
                    "application/json"
                ],
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
                            "type": "boolean"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    },
                    "404": {
                        "description": "Not Found",
                        "schema": {
                            "$ref": "#/definitions/handlers.ErrorHttp"
                        }
                    }
                }
            }
        },
        "/updateScheduler/": {
            "post": {
                "description": "The request will erase the current scheduler an create one new with\nthe requested entries for the scheduler. The entry will be definied by the initial hour\nand the ending hour, adintional info must be indicated depending of the kind of hours\nthe kinds of subject hours are:\n- Theorical = 1\n- Practices = 2\n- Exercises = 3",
                "produces": [
                    "text/plain"
                ],
                "parameters": [
                    {
                        "type": "string",
                        "description": "titulacion de las horas a obtener",
                        "name": "degree",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "curso de las horas a obtener",
                        "name": "year",
                        "in": "query",
                        "required": true
                    },
                    {
                        "type": "integer",
                        "description": "grupo de las horas a obtener",
                        "name": "group",
                        "in": "query",
                        "required": true
                    },
                    {
                        "description": "Entry to create",
                        "name": "entry",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "type": "array",
                            "items": {
                                "$ref": "#/definitions/handlers.EntryDTO"
                            }
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Receive the date of the latests entry modification with format dd/mm/aaaa"
                    }
                }
            }
        }
    },
    "definitions": {
        "domain.AvailableHours": {
            "type": "object",
            "properties": {
                "maxHours": {
                    "type": "integer"
                },
                "maxMin": {
                    "type": "integer"
                },
                "remainingHours": {
                    "type": "integer"
                },
                "remainingMin": {
                    "type": "integer"
                },
                "subject": {
                    "$ref": "#/definitions/domain.Subject"
                }
            }
        },
        "domain.DegreeDescription": {
            "type": "object",
            "properties": {
                "name": {
                    "type": "string"
                },
                "years": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.YearDescription"
                    }
                }
            }
        },
        "domain.Subject": {
            "type": "object",
            "properties": {
                "kind": {
                    "type": "integer"
                },
                "name": {
                    "type": "string"
                }
            }
        },
        "domain.YearDescription": {
            "type": "object",
            "properties": {
                "groups": {
                    "type": "array",
                    "items": {
                        "type": "string"
                    }
                },
                "name": {
                    "type": "integer"
                }
            }
        },
        "handlers.EntryDTO": {
            "type": "object",
            "properties": {
                "endHour": {
                    "type": "integer"
                },
                "endMin": {
                    "type": "integer"
                },
                "grupo": {
                    "type": "string"
                },
                "initHour": {
                    "type": "integer"
                },
                "initMin": {
                    "type": "integer"
                },
                "kind": {
                    "type": "integer"
                },
                "room": {
                    "type": "string"
                },
                "semana": {
                    "type": "string"
                },
                "subject": {
                    "type": "string"
                },
                "weekday": {
                    "type": "integer"
                }
            }
        },
        "handlers.ErrorHttp": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "handlers.ListDegreesDTO": {
            "type": "object",
            "properties": {
                "list": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/domain.DegreeDescription"
                    }
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
	Version:     "",
	Host:        "",
	BasePath:    "",
	Schemes:     []string{},
	Title:       "",
	Description: "",
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

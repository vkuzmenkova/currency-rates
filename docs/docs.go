// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {
            "name": "Valentina Kuzmenkova",
            "email": "valentinakuzmenkova@gmail.com"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/rates": {
            "get": {
                "description": "Gets rate by UUID",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rate"
                ],
                "summary": "Gets rate by UUID",
                "parameters": [
                    {
                        "type": "string",
                        "format": "string",
                        "description": "uuid of update",
                        "name": "uuid",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.CurrencyRate"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/currencyrates.NoUUIDFoundError"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "type": "string"
                        }
                    }
                }
            }
        },
        "/rates/{code}": {
            "get": {
                "description": "Gets rate from the last update",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rate"
                ],
                "summary": "Gets rate from the last update",
                "parameters": [
                    {
                        "type": "string",
                        "format": "string",
                        "description": "currency base, f.e. EUR",
                        "name": "code",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "format": "string",
                        "description": "currency base, default=USD",
                        "name": "base",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.CurrencyRate"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.BaseAndCodeAreEqual"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.BaseAndCodeAreEqual"
                        }
                    }
                }
            }
        },
        "/rates/{code}/update": {
            "get": {
                "description": "Initiates updating rate",
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "rate"
                ],
                "summary": "Initiates updating rate",
                "parameters": [
                    {
                        "type": "string",
                        "format": "string",
                        "description": "currency base, f.e. EUR",
                        "name": "code",
                        "in": "path"
                    },
                    {
                        "type": "string",
                        "format": "string",
                        "description": "currency base, default=USD",
                        "name": "base",
                        "in": "query"
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/models.CurrencyUpdateUUID"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/controller.BaseAndCodeAreEqual"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/controller.BaseAndCodeAreEqual"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "controller.BaseAndCodeAreEqual": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "controller.InvalidUUIDError": {
            "type": "object"
        },
        "controller.UnavailableCurrencyError": {
            "type": "object",
            "properties": {
                "currencyList": {
                    "type": "string"
                }
            }
        },
        "currencyrates.NoUUIDFoundError": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "models.CurrencyRate": {
            "type": "object",
            "properties": {
                "base": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                },
                "value": {
                    "type": "number"
                }
            }
        },
        "models.CurrencyUpdateUUID": {
            "type": "object",
            "properties": {
                "base": {
                    "type": "string"
                },
                "currency": {
                    "type": "string"
                },
                "uuid": {
                    "type": "string"
                }
            }
        }
    },
    "externalDocs": {
        "description": "OpenAPI",
        "url": "https://swagger.io/resources/open-api/"
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Swagger Currency Rates API",
	Description:      "This is a currency rates service.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

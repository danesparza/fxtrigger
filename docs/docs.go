// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "Apache 2.0",
            "url": "http://www.apache.org/licenses/LICENSE-2.0.html"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/trigger/fire/{id}": {
            "post": {
                "description": "Fires a trigger in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "triggers"
                ],
                "summary": "Fires a trigger in the system",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The trigger id to fire",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/triggers": {
            "get": {
                "description": "List all triggers in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "triggers"
                ],
                "summary": "List all triggers in the system",
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a trigger",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "triggers"
                ],
                "summary": "Update a trigger",
                "parameters": [
                    {
                        "description": "The trigger to update.  Must include trigger.id",
                        "name": "trigger",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.UpdateTriggerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            },
            "post": {
                "description": "Create a new trigger",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "triggers"
                ],
                "summary": "Create a new trigger",
                "parameters": [
                    {
                        "description": "The trigger to create",
                        "name": "trigger",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/api.CreateTriggerRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/triggers/{id}": {
            "delete": {
                "description": "Deletes a trigger in the system",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "triggers"
                ],
                "summary": "Deletes a trigger in the system",
                "parameters": [
                    {
                        "type": "string",
                        "description": "The trigger id to delete",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/api.SystemResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    },
                    "503": {
                        "description": "Service Unavailable",
                        "schema": {
                            "$ref": "#/definitions/api.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "api.CreateTriggerRequest": {
            "type": "object",
            "properties": {
                "description": {
                    "description": "Additional information about the trigger",
                    "type": "string"
                },
                "gpiopin": {
                    "description": "The GPIO pin the sensor or button is on",
                    "type": "integer"
                },
                "minimumsecondsbeforeretrigger": {
                    "description": "Minimum time (in seconds) before a retrigger",
                    "type": "integer"
                },
                "name": {
                    "description": "The trigger name",
                    "type": "string"
                },
                "webhooks": {
                    "description": "The webhooks to send when triggered",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/data.WebHook"
                    }
                }
            }
        },
        "api.ErrorResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        },
        "api.SystemResponse": {
            "type": "object",
            "properties": {
                "data": {},
                "message": {
                    "type": "string"
                }
            }
        },
        "api.UpdateTriggerRequest": {
            "type": "object",
            "properties": {
                "description": {
                    "description": "Additional information about the trigger",
                    "type": "string"
                },
                "enabled": {
                    "description": "Trigger enabled or not",
                    "type": "boolean"
                },
                "gpiopin": {
                    "description": "The GPIO pin the sensor or button is on",
                    "type": "integer"
                },
                "id": {
                    "description": "Unique Trigger ID",
                    "type": "string"
                },
                "minimumsecondsbeforeretrigger": {
                    "description": "Minimum time (in seconds) before a retrigger",
                    "type": "integer"
                },
                "name": {
                    "description": "The trigger name",
                    "type": "string"
                },
                "webhooks": {
                    "description": "The webhooks to send when triggered",
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/data.WebHook"
                    }
                }
            }
        },
        "data.WebHook": {
            "type": "object",
            "properties": {
                "body": {
                    "description": "The HTTP body to send.  This can be empty",
                    "type": "array",
                    "items": {
                        "type": "integer"
                    }
                },
                "headers": {
                    "description": "The HTTP headers to send",
                    "type": "object",
                    "additionalProperties": {
                        "type": "string"
                    }
                },
                "url": {
                    "description": "The URL to connect to",
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0",
	Host:             "",
	BasePath:         "/v1",
	Schemes:          []string{},
	Title:            "fxTrigger",
	Description:      "fxTrigger REST based management for GPIO/Sensor -> endpoint triggers (on Raspberry Pi)",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

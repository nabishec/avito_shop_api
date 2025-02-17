// Package docs Code generated by swaggo/swag. DO NOT EDIT
package docs

import "github.com/swaggo/swag"

const docTemplate = `{
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
        "/api/auth": {
            "post": {
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Аутентификация и получение JWT-токена.",
                "parameters": [
                    {
                        "description": "Данные для аутентификации",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.AuthRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешная аутентификация.",
                        "schema": {
                            "$ref": "#/definitions/model.AuthResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Неавторизован.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/buy/{item}": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Купить предмет за монеты.",
                "responses": {
                    "200": {
                        "description": "Успешный ответ."
                    },
                    "400": {
                        "description": "Неверный запрос.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Неавторизован.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/info": {
            "get": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Получить информацию о монетах, инвентаре и истории транзакций.",
                "responses": {
                    "200": {
                        "description": "Успешный ответ.",
                        "schema": {
                            "$ref": "#/definitions/model.InfoResponse"
                        }
                    },
                    "400": {
                        "description": "Неверный запрос.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Неавторизован.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        },
        "/api/sendCoin": {
            "post": {
                "security": [
                    {
                        "BearerAuth": []
                    }
                ],
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "summary": "Отправить монеты другому пользователю.",
                "parameters": [
                    {
                        "description": "Данные для отправки монет",
                        "name": "body",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/model.SendCoinRequest"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "Успешный ответ."
                    },
                    "400": {
                        "description": "Неверный запрос.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "401": {
                        "description": "Неавторизован.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    },
                    "500": {
                        "description": "Внутренняя ошибка сервера.",
                        "schema": {
                            "$ref": "#/definitions/model.ErrorResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "model.AuthRequest": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "model.AuthResponse": {
            "type": "object",
            "properties": {
                "token": {
                    "type": "string"
                }
            }
        },
        "model.CoinHistory": {
            "type": "object",
            "properties": {
                "received": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Received"
                    }
                },
                "sent": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Sent"
                    }
                }
            }
        },
        "model.ErrorResponse": {
            "type": "object",
            "properties": {
                "errors": {
                    "type": "string"
                }
            }
        },
        "model.InfoResponse": {
            "type": "object",
            "properties": {
                "coinHistory": {
                    "$ref": "#/definitions/model.CoinHistory"
                },
                "coins": {
                    "type": "integer"
                },
                "inventory": {
                    "type": "array",
                    "items": {
                        "$ref": "#/definitions/model.Item"
                    }
                }
            }
        },
        "model.Item": {
            "type": "object",
            "properties": {
                "quantity": {
                    "type": "integer"
                },
                "type": {
                    "type": "string"
                }
            }
        },
        "model.Received": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "fromUser": {
                    "type": "string"
                }
            }
        },
        "model.SendCoinRequest": {
            "type": "object",
            "required": [
                "amount",
                "toUser"
            ],
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "toUser": {
                    "type": "string"
                }
            }
        },
        "model.Sent": {
            "type": "object",
            "properties": {
                "amount": {
                    "type": "integer"
                },
                "toUser": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "BearerAuth": {
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "localhost:8080",
	BasePath:         "/",
	Schemes:          []string{"http"},
	Title:            "API Avito shop",
	Description:      "",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

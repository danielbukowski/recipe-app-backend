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
        "/api/v1/recipes": {
            "post": {
                "description": "Insert a new recipe by providing a request body with a title and a content for the recipe.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Create a recipe",
                "parameters": [
                    {
                        "description": "Request body with title and content",
                        "name": "NewRecipeRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/recipe.NewRecipeRequest"
                        }
                    }
                ],
                "responses": {
                    "201": {
                        "description": "Created",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "408": {
                        "description": "Request Timeout",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    }
                }
            }
        },
        "/api/v1/recipes/{id}": {
            "get": {
                "description": "Get a recipe by ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Get a recipe",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID for a recipe",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "200": {
                        "description": "OK",
                        "schema": {
                            "$ref": "#/definitions/github_com_danielbukowski_recipe-app-backend_internal_shared.DataResponse-recipe_RecipeResponse"
                        }
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "408": {
                        "description": "Request Timeout",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    }
                }
            },
            "put": {
                "description": "Update a title or a content of a recipe by ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Update a recipe",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID for a recipe resource",
                        "name": "id",
                        "in": "path",
                        "required": true
                    },
                    {
                        "description": "Request body for updating title and content fields of a recipe",
                        "name": "UpdateRecipeRequest",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/recipe.UpdateRecipeRequest"
                        }
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "408": {
                        "description": "Request Timeout",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    }
                }
            },
            "delete": {
                "description": "Delete a recipe by ID.",
                "consumes": [
                    "application/json"
                ],
                "produces": [
                    "application/json"
                ],
                "tags": [
                    "recipes"
                ],
                "summary": "Delete a recipe",
                "parameters": [
                    {
                        "type": "string",
                        "description": "UUID for a recipe",
                        "name": "id",
                        "in": "path",
                        "required": true
                    }
                ],
                "responses": {
                    "204": {
                        "description": "No Content"
                    },
                    "400": {
                        "description": "Bad Request",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "408": {
                        "description": "Request Timeout",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    },
                    "500": {
                        "description": "Internal Server Error",
                        "schema": {
                            "$ref": "#/definitions/shared.CommonResponse"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "github_com_danielbukowski_recipe-app-backend_internal_shared.DataResponse-recipe_RecipeResponse": {
            "type": "object",
            "properties": {
                "data": {
                    "$ref": "#/definitions/recipe.RecipeResponse"
                }
            }
        },
        "recipe.NewRecipeRequest": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "recipe.RecipeResponse": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "created_at": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                },
                "updated_at": {
                    "type": "string"
                }
            }
        },
        "recipe.UpdateRecipeRequest": {
            "type": "object",
            "properties": {
                "content": {
                    "type": "string"
                },
                "title": {
                    "type": "string"
                }
            }
        },
        "shared.CommonResponse": {
            "type": "object",
            "properties": {
                "message": {
                    "type": "string"
                }
            }
        }
    }
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "0.1",
	Host:             "localhost:8080",
	BasePath:         "/api/v1",
	Schemes:          []string{},
	Title:            "Recipe API",
	Description:      "A sample of API to recipe backend.",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}

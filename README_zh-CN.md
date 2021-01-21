# swag

🌍 *[English](README.md) ∙ [简体中文](README_zh-CN.md)*

<img align="right" width="180px" src="https://raw.githubusercontent.com/swaggo/swag/master/assets/swaggo.png">

[![Travis Status](https://img.shields.io/travis/swaggo/swag/master.svg)](https://travis-ci.org/swaggo/swag)
[![Coverage Status](https://img.shields.io/codecov/c/github/swaggo/swag/master.svg)](https://codecov.io/gh/swaggo/swag)
[![Go Report Card](https://goreportcard.com/badge/github.com/swaggo/swag)](https://goreportcard.com/report/github.com/swaggo/swag)
[![codebeat badge](https://codebeat.co/badges/71e2f5e5-9e6b-405d-baf9-7cc8b5037330)](https://codebeat.co/projects/github-com-swaggo-swag-master)
[![Go Doc](https://godoc.org/github.com/swaggo/swagg?status.svg)](https://godoc.org/github.com/swaggo/swag)
[![Backers on Open Collective](https://opencollective.com/swag/backers/badge.svg)](#backers) 
[![Sponsors on Open Collective](https://opencollective.com/swag/sponsors/badge.svg)](#sponsors) [![FOSSA Status](https://app.fossa.io/api/projects/git%2Bgithub.com%2Fswaggo%2Fswag.svg?type=shield)](https://app.fossa.io/projects/git%2Bgithub.com%2Fswaggo%2Fswag?ref=badge_shield)
[![Release](https://img.shields.io/github/release/swaggo/swag.svg?style=flat-square)](https://github.com/swaggo/swag/releases)

Swag将Go的注释转换为Swagger2.0文档。我们为流行的 [Go Web Framework](#支持的Web框架) 创建了各种插件，这样可以与现有Go项目快速集成（使用Swagger UI）。

## 目录

- [快速开始](#快速开始)
- [支持的Web框架](#支持的web框架)
- [如何与Gin集成](#如何与gin集成)
- [格式化说明](#格式化说明)
- [开发现状](#开发现状)
- [声明式注释格式](#声明式注释格式)
    - [通用API信息](#通用api信息)
    - [API操作](#api操作)
    - [安全性](#安全性)
- [样例](#样例)
    - [多行的描述](#多行的描述)
    - [用户自定义的具有数组类型的结构](#用户自定义的具有数组类型的结构)
    - [响应对象中的模型组合](#响应对象中的模型组合)
    - [在响应中增加头字段](#在响应中增加头字段)
    - [使用多路径参数](#使用多路径参数)
    - [结构体的示例值](#结构体的示例值)
    - [结构体描述](#结构体描述)
    - [使用`swaggertype`标签更改字段类型](#使用`swaggertype`标签更改字段类型)
    - [使用`swaggerignore`标签排除字段](#使用swaggerignore标签排除字段)
    - [将扩展信息添加到结构字段](#将扩展信息添加到结构字段)
    - [对展示的模型重命名](#对展示的模型重命名)
    - [如何使用安全性注释](#如何使用安全性注释)
- [项目相关](#项目相关)

## 快速开始

1. 将注释添加到API源代码中，请参阅声明性注释格式。
2. 使用如下命令下载swag：

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

从源码开始构建的话，需要有Go环境（1.16及以上版本）。

或者从github的release页面下载预编译好的二进制文件。

3. 在包含`main.go`文件的项目根目录运行`swag init`。这将会解析注释并生成需要的文件（`docs`文件夹和`docs/docs.go`）。

```bash
swag init
```

确保导入了生成的`docs/docs.go`文件，这样特定的配置文件才会被初始化。如果通用API注释没有写在`main.go`中，可以使用`-g`标识符来告知swag。

```bash
swag init -g http/api.go
```

4. (可选) 使用`fmt`格式化 SWAG 注释。(请先升级到最新版本)

```bash
swag fmt
```

## swag cli

```bash
swag init -h
NAME:
   swag init - Create docs.go

USAGE:
   swag init [command options] [arguments...]

OPTIONS:
   --generalInfo value, -g value          API通用信息所在的go源文件路径，如果是相对路径则基于API解析目录 (默认: "main.go")
   --dir value, -d value                  API解析目录 (默认: "./")
   --exclude value                        解析扫描时排除的目录，多个目录可用逗号分隔（默认：空）
   --propertyStrategy value, -p value     结构体字段命名规则，三种：snakecase,camelcase,pascalcase (默认: "camelcase")
   --output value, -o value               文件(swagger.json, swagger.yaml and doc.go)输出目录 (默认: "./docs")
   --parseVendor                          是否解析vendor目录里的go源文件，默认不
   --parseDependency                      是否解析依赖目录中的go源文件，默认不
   --markdownFiles value, --md value      指定API的描述信息所使用的markdown文件所在的目录
   --generatedTime                        是否输出时间到输出文件docs.go的顶部，默认是
   --codeExampleFiles value, --cef value  解析包含用于 x-codeSamples 扩展的代码示例文件的文件夹，默认禁用
   --parseInternal                        解析 internal 包中的go文件，默认禁用
   --parseDepth value                     依赖解析深度 (默认: 100)
   --instanceName value                   设置文档实例名 (默认: "swagger")
```

```bash
swag fmt -h
NAME:
   swag fmt - format swag comments

USAGE:
   swag fmt [command options] [arguments...]

OPTIONS:
   --dir value, -d value          API解析目录 (默认: "./")
   --exclude value                解析扫描时排除的目录，多个目录可用逗号分隔（默认：空）
   --generalInfo value, -g value  API通用信息所在的go源文件路径，如果是相对路径则基于API解析目录 (默认: "main.go")
   --help, -h                     show help (default: false)

```

## 支持的Web框架

- [gin](http://github.com/swaggo/gin-swagger)
- [echo](http://github.com/swaggo/echo-swagger)
- [buffalo](https://github.com/swaggo/buffalo-swagger)
- [net/http](https://github.com/swaggo/http-swagger)
- [net/http](https://github.com/swaggo/http-swagger)
- [gorilla/mux](https://github.com/swaggo/http-swagger)
- [go-chi/chi](https://github.com/swaggo/http-swagger)
- [flamingo](https://github.com/i-love-flamingo/swagger)
- [fiber](https://github.com/gofiber/swagger)
- [atreugo](https://github.com/Nerzal/atreugo-swagger)
- [hertz](https://github.com/hertz-contrib/swagger)

## 如何与Gin集成

[点击此处](https://github.com/swaggo/swag/tree/master/example/celler)查看示例源代码。

1. 使用`swag init`生成Swagger2.0文档后，导入如下代码包：

```go
import "github.com/swaggo/gin-swagger" // gin-swagger middleware
import "github.com/swaggo/files" // swagger embed files
```

2. 在`main.go`源代码中添加通用的API注释：

```go
// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
    r := gin.Default()

    c := controller.NewController()

    v1 := r.Group("/api/v1")
    {
        accounts := v1.Group("/accounts")
        {
            accounts.GET(":id", c.ShowAccount)
            accounts.GET("", c.ListAccounts)
            accounts.POST("", c.AddAccount)
            accounts.DELETE(":id", c.DeleteAccount)
            accounts.PATCH(":id", c.UpdateAccount)
            accounts.POST(":id/images", c.UploadAccountImage)
        }
    //...
    }
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
    r.Run(":8080")
}
//...
```

此外，可以动态设置一些通用的API信息。生成的代码包`docs`导出`SwaggerInfo`变量，使用该变量可以通过编码的方式设置标题、描述、版本、主机和基础路径。使用Gin的示例：

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/swaggo/files"
    "github.com/swaggo/gin-swagger"

    "./docs" // docs is generated by Swag CLI, you have to import it.
)

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

    // programatically set swagger info
    docs.SwaggerInfo.Title = "Swagger Example API"
    docs.SwaggerInfo.Description = "This is a sample server Petstore server."
    docs.SwaggerInfo.Version = "1.0"
    docs.SwaggerInfo.Host = "petstore.swagger.io"
    docs.SwaggerInfo.BasePath = "/v2"
    docs.SwaggerInfo.Schemes = []string{"http", "https"}

    r := gin.New()

    // use ginSwagger middleware to serve the API docs
    r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

    r.Run()
}
```

3. 在`controller`代码中添加API操作注释：

```go
package controller

import (
    "fmt"
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/swaggo/swag/example/celler/httputil"
    "github.com/swaggo/swag/example/celler/model"
)

// ShowAccount godoc
// @Summary      Show an account
// @Description  get string by ID
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Account ID"
// @Success      200  {object}  model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {
  id := ctx.Param("id")
  aid, err := strconv.Atoi(id)
  if err != nil {
    httputil.NewError(ctx, http.StatusBadRequest, err)
    return
  }
  account, err := model.AccountOne(aid)
  if err != nil {
    httputil.NewError(ctx, http.StatusNotFound, err)
    return
  }
  ctx.JSON(http.StatusOK, account)
}

// ListAccounts godoc
// @Summary      List accounts
// @Description  get accounts
// @Tags         accounts
// @Accept       json
// @Produce      json
// @Param        q    query     string  false  "name search by q"  Format(email)
// @Success      200  {array}   model.Account
// @Failure      400  {object}  httputil.HTTPError
// @Failure      404  {object}  httputil.HTTPError
// @Failure      500  {object}  httputil.HTTPError
// @Router       /accounts [get]
func (c *Controller) ListAccounts(ctx *gin.Context) {
  q := ctx.Request.URL.Query().Get("q")
  accounts, err := model.AccountsAll(q)
  if err != nil {
    httputil.NewError(ctx, http.StatusNotFound, err)
    return
  }
  ctx.JSON(http.StatusOK, accounts)
}
//...
```

```bash
swag init
```

4. 运行程序，然后在浏览器中访问 http://localhost:8080/swagger/index.html 。将看到Swagger 2.0 Api文档，如下所示：

![swagger_index.html](https://raw.githubusercontent.com/swaggo/swag/master/assets/swagger-image.png)

## 格式化说明

可以针对Swag的注释自动格式化，就像`go fmt`。   
此处查看格式化结果 [here](https://github.com/swaggo/swag/tree/master/example/celler).

示例：
```shell
swag fmt
```

排除目录（不扫描）示例：
```shell
swag fmt -d ./ --exclude ./internal
```

## 开发现状

[Swagger 2.0 文档](https://swagger.io/docs/specification/2-0/basic-structure/)

- [x] Basic Structure
- [x] API Host and Base Path
- [x] Paths and Operations
- [x] Describing Parameters
- [x] Describing Request Body
- [x] Describing Responses
- [x] MIME Types
- [x] Authentication
  - [x] Basic Authentication
  - [x] API Keys
- [x] Adding Examples
- [x] File Upload
- [x] Enums
- [x] Grouping Operations With Tags
- [ ] Swagger Extensions

## 声明式注释格式

## 通用API信息

**示例** [`celler/main.go`](https://github.com/swaggo/swag/blob/master/example/celler/main.go)

| 注释                    | 说明                                                                                            | 示例                                                            |
| ----------------------- | ----------------------------------------------------------------------------------------------- | --------------------------------------------------------------- |
| title                   | **必填** 应用程序的名称。                                                                       | // @title Swagger Example API                                   |
| version                 | **必填** 提供应用程序API的版本。                                                                | // @version 1.0                                                 |
| description             | 应用程序的简短描述。                                                                            | // @description This is a sample server celler server.          |
| tag.name                | 标签的名称。                                                                                    | // @tag.name This is the name of the tag                        |
| tag.description         | 标签的描述。                                                                                    | // @tag.description Cool Description                            |
| tag.docs.url            | 标签的外部文档的URL。                                                                           | // @tag.docs.url https://example.com                            |
| tag.docs.description    | 标签的外部文档说明。                                                                            | // @tag.docs.description Best example documentation             |
| termsOfService          | API的服务条款。                                                                                 | // @termsOfService http://swagger.io/terms/                     |
| contact.name            | 公开的API的联系信息。                                                                           | // @contact.name API Support                                    |
| contact.url             | 联系信息的URL。 必须采用网址格式。                                                              | // @contact.url http://www.swagger.io/support                   |
| contact.email           | 联系人/组织的电子邮件地址。 必须采用电子邮件地址的格式。                                        | // @contact.email support@swagger.io                            |
| license.name            | **必填** 用于API的许可证名称。                                                                  | // @license.name Apache 2.0                                     |
| license.url             | 用于API的许可证的URL。 必须采用网址格式。                                                       | // @license.url http://www.apache.org/licenses/LICENSE-2.0.html |
| host                    | 运行API的主机（主机名或IP地址）。                                                               | // @host localhost:8080                                         |
| BasePath                | 运行API的基本路径。                                                                             | // @BasePath /api/v1                                            |
| accept                  | API 可以使用的 MIME 类型列表。 请注意，Accept 仅影响具有请求正文的操作，例如 POST、PUT 和 PATCH。 值必须如“[Mime类型](#mime类型)”中所述。                                  | // @accept json |
| produce                 | API可以生成的MIME类型的列表。值必须如“[Mime类型](#mime类型)”中所述。                                  | // @produce json |
| query.collection.format | 请求URI query里数组参数的默认格式：csv，multi，pipes，tsv，ssv。 如果未设置，则默认为csv。 | // @query.collection.format multi                               |
| schemes                 | 用空格分隔的请求的传输协议。                                                                    | // @schemes http https                                          |
| externalDocs.description | Description of the external document. | // @externalDocs.description OpenAPI |
| externalDocs.url         | URL of the external document. | // @externalDocs.url https://swagger.io/resources/open-api/ |
| x-name                  | 扩展的键必须以x-开头，并且只能使用json值                                                        | // @x-example-key {"key": "value"}                              |

### 使用Markdown描述

如果文档中的短字符串不足以完整表达，或者需要展示图片，代码示例等类似的内容，则可能需要使用Markdown描述。要使用Markdown描述，请使用一下注释。

| 注释                     | 说明                                                                                 | 示例                                                                              |
| ------------------------ | ------------------------------------------------------------------------------------ | --------------------------------------------------------------------------------- |
| title                    | **必填** 应用程序的名称。                                                            | // @title Swagger Example API                                                     |
| version                  | **必填** 提供应用程序API的版本。                                                     | // @version 1.0                                                                   |
| description.markdown     | 应用程序的简短描述。 从`api.md`文件中解析。 这是`@description`的替代用法。           | // @description.markdown No value needed, this parses the description from api.md |
| tag.name                 | 标签的名称。                                                                         | // @tag.name This is the name of the tag                                          |
| tag.description.markdown | 标签说明，这是`tag.description`的替代用法。 该描述将从名为`tagname.md的`文件中读取。 | // @tag.description.markdown                                                      |

## API操作

Example [celler/controller](https://github.com/swaggo/swag/tree/master/example/celler/controller)

| 注释                 | 描述                                                                                                    |
| -------------------- | ------------------------------------------------------------------------------------------------------- |
| description          | 操作行为的详细说明。                                                                                    |
| description.markdown | 应用程序的简短描述。该描述将从名为`endpointname.md`的文件中读取。                                       |
| id                   | 用于标识操作的唯一字符串。在所有API操作中必须唯一。                                                     |
| tags                 | 每个API操作的标签列表，以逗号分隔。                                                                     |
| summary              | 该操作的简短摘要。                                                                                      |
| accept               | API 可以使用的 MIME 类型列表。 请注意，Accept 仅影响具有请求正文的操作，例如 POST、PUT 和 PATCH。 值必须如“[Mime类型](#mime类型)”中所述。                                  |
| produce              | API可以生成的MIME类型的列表。值必须如“[Mime类型](#mime类型)”中所述。                                  |
| param                | 用空格分隔的参数。`param name`,`param type`,`data type`,`is mandatory?`,`comment` `attribute(optional)` |
| security             | 每个API操作的[安全性](#安全性)。                                                                      |
| success        
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
   --generatedTime  
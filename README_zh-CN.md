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

## 快速开
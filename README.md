# Laverna

[![Version](https://img.shields.io/github/tag/mrwormhole/laverna.svg)](https://github.com/mrwormhole/laverna/tags)
[![CI Build](https://github.com/mrwormhole/laverna/actions/workflows/tests.yaml/badge.svg)](https://github.com/mrwormhole/laverna/actions/workflows/tests.yaml)
[![GoDoc](https://godoc.org/github.com/mrwormhole/laverna?status.svg)](https://godoc.org/github.com/mrwormhole/laverna)
[![Report Card](https://goreportcard.com/badge/github.com/mrwormhole/laverna)](https://goreportcard.com/report/github.com/mrwormhole/laverna)
[![License](https://img.shields.io/github/license/mrwormhole/laverna)](https://github.com/mrwormhole/laverna/blob/main/LICENSE)
[![Coverage Status](https://coveralls.io/repos/github/mrwormhole/laverna/badge.svg?branch=main)](https://coveralls.io/github/mrwormhole/laverna?branch=main)

the goddess of the thieves, helps you steal translation speeches from the G-daddy monopoly.

<img src="https://github.com/user-attachments/assets/d1d344c9-f36b-4cf7-af70-f162f93ea9f0" width="400" alt="squash-goose">

### Install Via Go

```shell
  go install github.com/lingua-sensei/laverna@latest
```

### Grab Binaries

You can find binaries through GitHub releases.

### Sample Usage

Let's create example YAML.

```yaml
- speed: normal
  voice: th
  text: "สวัสดีครับ"
- speed: slower
  voice: en
  text: "Hello there"
- speed: slowest
  voice: ja
  text: "こんにちは~"
```

Running below command will generate audios in the same directory.

```shell
  laverna -file example.yaml 
```

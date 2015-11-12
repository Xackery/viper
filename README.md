# viper
A Tumblr API for Golang

[![godoc badge](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/Xackery/viper/tumblr) [![Go report](http://goreportcard.com/badge/xackery/discord)](http://goreportcard.com/report/xackery/viper) [![Build Status](https://travis-ci.org/Xackery/discord.svg)](https://travis-ci.org/Xackery/viper.svg?branch=master)

Features
---

* To do


Installation
----
To install, simply `go get github.com/xackery/viper/tumblr` in a command line, followed by the usage example


Usage Example
---


```
	tumblr.SetConsumerKey("YourKey")
	tumblr.SetConsumerSecret("YourSecret")
	api := tumblr.NewAPI("YourTokenKey", "YourTokenSecret")

```

Options
---

* Enable/Disable Throttling with the function calls
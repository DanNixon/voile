# Voile
[![Build Status](https://travis-ci.org/DanNixon/voile.svg?branch=master)](https://travis-ci.org/DanNixon/voile)
[![Go Report Card](https://goreportcard.com/badge/github.com/dannixon/voile)](https://goreportcard.com/report/github.com/dannixon/voile)

Command line bookmark management tool focused on simplicity and speed.

## Installation

```
go get github.com/DanNixon/voile
```

## Features

- Plain text (JSON) library
- Query by a combination of tags, name, URL and description
- Text editor based entry manipulation
- Open bookmarks in browser
- Copy bookmarks to and from clipboard
- Integration with Git if bookmarks are stored in a Git repository
- Integration with [Newsboat's](https://newsboat.org/) [bookmark plugin architecture](https://newsboat.org/releases/2.19/docs/newsboat.html#_bookmarking)
- Helper to prune old bookmarks/keep bookmarks up to date

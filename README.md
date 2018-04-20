# Voile [![Build Status](https://travis-ci.org/DanNixon/voile.svg?branch=master)](https://travis-ci.org/DanNixon/voile)

Command line bookmark management tool focused on simplicity and speed.

## Installation

```
go get github.com/DanNixon/voile
```

## Features

- Plain text (JSON) library (compatible with [Buku](https://github.com/jarun/buku) JSON format)
- Query by a combination of tags, name, URL and description
- Text editor based entry manipulation
- Open bookmarks in browser
- Copy bookmarks to and from clipboard

## To Do

- Add existing tags as commented line in editor (for autocompletion)
- Option to auto-set title from page title
- "Quick add" command (copies URL from clipboard, gets name from page title, takes tags from arguments)
- Editing via flags

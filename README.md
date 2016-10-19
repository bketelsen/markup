# ml
[![Build Status](https://travis-ci.org/murlokswarm/ml.svg?branch=master)](https://travis-ci.org/murlokswarm/ml)
[![Go Report Card](https://goreportcard.com/badge/github.com/murlokswarm/ml)](https://goreportcard.com/report/github.com/murlokswarm/ml)
[![Coverage Status](https://coveralls.io/repos/github/murlokswarm/ml/badge.svg?branch=master)](https://coveralls.io/github/murlokswarm/ml?branch=master)
[![GoDoc](https://godoc.org/github.com/murlokswarm/ml?status.svg)](https://godoc.org/github.com/murlokswarm/ml)

Package ml implements a markup language to build user interfaces.

Markups are based on HTML. They must be declared in the Render method when
implementing the Componer interface.
A markup must follow these rules:
- Regular HTML elements must be in lowercase.
- Component elements must have their first letter capitalized.
- Component root element must be a standard HTML tag.
- Each element must have a closing tag.
- HTML event handlers should start with '@'.
- Template must follow the rules of https://golang.org/pkg/text/template.

Example:
```go 
type Hello struct {
	Name string
}

func (c *Hello) OnInputChange(v string) string {
	c.Name = v
	app.Render(c)
}

func (c *Hello) Render() string {
	return `
 <p>
  	Hello,
 	<input onchange="@OnInputChange" />
 	<World Name="{{.Name}}" />
 </p>
 	`
}
```
```go 
type World struct {
	Name string
}

func (c *World) Render() string {
	return `
 <span>
 	{{if len .Name}}
    	{{.Name}}
  	{{else}}
      	World
  	{{end}}
 </span>
 	`
}

```
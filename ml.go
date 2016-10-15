// Package ml implements a markup language that allow to build user interfaces.
//
// Markups are based on HTML. They must be declared in the Render method when
// implementing the Componer interface.
//
// A markup must follow these rules:
// - Regular HTML elements must be in lowercase.
// - Component elements must have their first letter capitalized.
// - Component root element must be a standard HTML tag.
// - HTML event handlers should start with '@'.
// - Template must follow the rules of https://golang.org/pkg/text/template.
//
// Example:
// // type Hello struct {
// //	Name string
// // }
// //
// //
// //  func (c *Hello) OnInputChange(v string) string {
// //  	c.Name = v
// //  }
// //
// //  func (c *Hello) Render() string {
// //  	return `
// //  <p>
// //   Hello,
// //  <input onchange="@OnInputChange" />
// //  <World Name="{{.Name}}" />
// //  </p>
// //  	`
// //  }
// //
// //  type World struct {
// // 	Name string
// //  }
// //
// //  func (c *World) Render() string {
// //  	return `
// //  <span>
// //   {{if len .Name}}
// //       {{.Name}}
// //   {{else}}
// //       World
// //   {{end}}
// //  </span>
// //  	`
// //  }
package ml

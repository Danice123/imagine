package main

import (
	"html/template"
)

type ImageData struct {
	Url         string
	Mimetype    string
	IsVideo     bool
	HasNext     bool
	Next        string
	HasPrevious bool
	Previous    string
	RandomState string
}

var imageTemplateHtml = `
<html>
	<head>
		<meta name="viewport" content="width=device-width">
		<style>
			video {
				position: absolute;
				top: 0px;
				right: 0px;
				bottom: 0px;
				left: 0px;
				max-height: 100%;
				max-width: 100%;
				margin: auto;
			}

			img {
				position: absolute;
				top: 0px;
				right: 0px;
				bottom: 0px;
				left: 0px;
				max-height: 100%;
				max-width: 100%;
				margin: auto;
			}

			#random {
				display: none;
			}

			#next:hover #random {
				display: inline;
			}
		</style>
	</head>
	<body style="background-color: black; text-align: center;">
		{{if .IsVideo}}
		<video autoplay loop>
			<source src="{{.Url}}" type="{{.Mimetype}}">
		</video>
		{{else}}
		<img src="{{.Url}}" />
		{{end}}
		{{if .HasPrevious}}
		<div style="position: fixed; bottom: 0px; left: 0px; z-index: 2147483600;">
			<a href="{{.Previous}}">
				<img title="Previous" src="chrome-extension://lhlckkgdiojkapplglfeomlkjllphilo/img/arrow_left.png" width="77" style="float:left;position:relative;cursor:pointer;vertical-align: bottom;">
			</a>
		</div>
		{{end}}
		{{if .HasNext}}
		<div id="next" style="position: fixed; bottom: 0px; right: 0px; z-index: 2147483600;">
			<a href="{{.Next}}">
				<img title="Next" src="chrome-extension://lhlckkgdiojkapplglfeomlkjllphilo/img/arrow_right.png" width="77" style="float:right;position:relative;cursor:pointer;vertical-align: bottom;">
			</a>
			<a href="/api/random">
				<img id="random" title="Toggle Random ({{.RandomState}})" src="chrome-extension://lhlckkgdiojkapplglfeomlkjllphilo/img/r_left.png" width="28" style="cursor: pointer; position: relative; vertical-align: bottom;">
			</a>
		</div>
		{{end}}
	</body>
</html>
`

var imageTemplate = template.New("Image")

func init() {
	if _, err := imageTemplate.Parse(imageTemplateHtml); err != nil {
		panic(err.Error())
	}
}

<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<link rel="icon" type="image/vnd.microsoft.icon" href="{{.Config.Favicon}}">
	<style type="text/css">
	h1 {
  		color: {{.Config.Colors.Title}};
	}
	p,h1,h2,h3,body,div,span,form,button,input,label {
  		font-family: Calibri, Candara, Segoe, Segoe UI, Optima, Arial, sans-serif;
	}
	button {
		font-size: inherit;
		color: {{.Config.Colors.ButtonFG}};
		background: {{.Config.Colors.ButtonBG}};
		margin-right: 1em;
		border: none;
	}
	button:hover {
		font-weight: bold;
	}
	input {
		font-size: inherit;
	}
	</style>

	<title>{{template "pagetitle" .}}</title>

	{{template "jsdata" .}}
</head>
<body class="container-fluid" style="padding-top: 15px;">
	{{template "yield" .}}
</body>
</html>
{{define "pagetitle"}}{{end}}
{{define "jsdata"}}{{end}}
{{define "yield"}}{{end}}


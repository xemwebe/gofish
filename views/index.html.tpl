{{define "pagetitle"}}{{.Config.Title}}{{end}}

<h1>{{.Config.Title}}</h1>
<h2>Dateien und Verzeichnisse</h2>
<p>Aktueller Pfad: <a href="/home">home</a>{{range .LocalPath}} &raquo; <a href="/home/{{.}}">{{lastStub .}}</a>{{end}}
</p>
{{if isAdmin}}
  {{template "_uploadForm" .}}
{{else}}
  {{template "_downloadForm" .}}
{{end}}
<br>
<p>Bei Fragen kontaktieren Sie bitte {{.Config.Author}} per
<a href="mailto:{{.Config.EMail}}?Subject={{.Config.EMailSubject}}">eMail</a>.</p>


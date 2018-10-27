{{$fullList := lastElement .LocalPath }}
{{$fullpath := join "/home" $fullList }}
{{$servePath := join "/serve" $fullList }}
{{range .Dirs}}
  <img src="/images/folder.png" alt="Folder" style="height:1em;"> <a href="{{$fullpath}}{{.}}">{{.}}</a>
  <br>
{{end}}
{{range .Files}}
  <a href="{{$servePath}}{{.}}">{{.}}</a>
  <br>
{{end}}

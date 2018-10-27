<style>
    label.upLabel input[type="file"] {
        position: absolute;
        top: -1000px;
        background: green;
    }
    .upLabel {
        border: 1px solid #AAA;
        padding: 2px 5px;
        margin: 2px;
        display: inline-block;
    }
    .upLabel:hover {
        background: #CCC;
    }
    .upLabel:active {
        background: #CCF;
    }
    .upLabel:invalid+span {
        color: #4A4;
    }
    .upLabel:valid+span {
        color: #4A4;
    }
</style>
<script type="text/javascript">
    function updateFileName(obj) {
        var label = obj.parentNode;

        var fileName = '';
		if( obj.files && obj.files.length > 1 )
			fileName = ( "{count} Dateien ausgewählt." ).replace( '{count}', obj.files.length );
		else
			fileName = obj.value.split( '\\' ).pop();

		if( fileName )
			label.querySelector( 'span' ).innerHTML = fileName;
		else
			label.innerHTML = "Wähle Datei zum Hochladen...";
    }
</script>
<form method="post">
  <div class="form-check">
    {{$fullList := lastElement .LocalPath }}
    {{$fullpath := join "/home" $fullList }}
    {{$servePath := join "/serve" $fullList }}
    {{range .Dirs}}
      <input name="{{.}}" type="checkbox" class="form-check-input" id="id_{{.}}">
      <label class="form-check-label" for="id_{{.}}"><img src="/images/folder.png" alt="Folder" style="height:1em;"> <a href="{{$fullpath}}{{.}}">{{.}}</a></label>
      <br>
    {{end}}
    {{range .Files}}
      <input name="{{.}}" type="checkbox" class="form-check-input" id="id_{{.}}">
      <label class="form-check-label" for="id_{{.}}"><a href="{{$servePath}}{{.}}">{{.}}</a></label>
      <br>
    {{end}}
  </div>
  <h2>Administration</h2>
   <button type="submit" class="btn btn-primary">Ausgewählte Dateien löschen</button>
</form>
<br>
<form action="/newDir" method="post">
    <label for="dirName">Neues Verzeichnis</label>
    <input name="fullpath" value="{{$fullpath}}" hidden=true>
    <input name="dirName" id="dirName" type="text">
    <button type="submit" class="btn btn-warning">erzeugen</button>
</form>
<br>
<form method="post" enctype="multipart/form-data" action="/upload">
    <input name="fullpath" value="{{$fullpath}}" hidden="true">
    <label class="upLabel"><input type="file" name="uploadfile" required onchange="updateFileName(this);" multiple/>
        <span>Wähle Datei zum Hochladen...</span>
    </label>
    <button type="submit">Hochladen</button>
</form>

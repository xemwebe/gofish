package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aarondl/tpl"
	"github.com/golang/glog"
	"golang.org/x/crypto/bcrypt"
)

type ColorConf struct {
	Title    string `json:"title"`
	ButtonFG string `json:"buttonfg"`
	ButtonBG string `json:"buttonbg"`
}

type Configuration struct {
	AllowAdmin   bool      `json:"allow_admin"`
	FilePath     string    `json:"file_path"`
	IpAddress    string    `json:"ip_address"`
	Port         string    `json:"port"`
	Title        string    `json:"title"`
	Favicon      string    `json:"favicon"`
	Realm        string    `json:"realm"`
	Author       string    `json:"author"`
	EMail        string    `json:"email"`
	EMailSubject string    `json:"subject"`
	UserPWHash   string    `json:"userpwhash`
	AdminPWHash  string    `json:"adminpwhash`
	Colors       ColorConf `json:"colors"`
}

func isAdmin() bool {
	return adminFlag
}

func joinPath(a string, b ...string) string {
	js := a
	for _, s := range b {
		if s == "/" {
			s = ""
		}
		js = js + "/" + s
	}
	if len(js) == 0 || js[len(js)-1:] != "/" {
		js += "/"
	}
	glog.Infof("JoinPath: a:%v, b:%v, js:%v", a, b, js)
	return js
}

func splitString(s string, sep rune) []string {
	return strings.FieldsFunc(s, func(c rune) bool { return c == sep })
}

var (
	templates = tpl.Must(tpl.Load("views", "views/partials", "layout.html.tpl", funcs))
	funcs     = template.FuncMap{
		"yield":    func() string { return "" },
		"join":     func(a, b string) string { return joinPath(a, b) },
		"joinList": func(a []string) string { return joinPath(a[0], a[1:]...) },
		"lastStub": func(a string) string {
			s := splitString(a, '/')
			glog.Infof("lastStub: s=%v", s)
			i := len(s) - 1
			for ; i > 0; i-- {
				if s[i] != "" {
					break
				}
			}
			return s[i]
		},
		"lastElement": func(a []string) string {
			if len(a) == 0 {
				return ""
			}
			return a[len(a)-1]
		},
		"isAdmin": isAdmin,
	}
	Conf = Configuration{
		AllowAdmin:  false,
		FilePath:    "/var/www/My_File_Storage",
		IpAddress:   "127.0.0.1",
		Port:        "7356",
		Title:       "My Private Web File Sharing Site",
		Favicon:     "./favicon.ico",
		UserPWHash:  "",
		AdminPWHash: "",
		Colors: ColorConf{
			Title:    "Navy",
			ButtonFG: "OrangeRed",
			ButtonBG: "Gold"},
	}
	adminFlag = false
)

func main() {
	genPW := flag.Bool("gen-pwd", false, "Generate Password Hash")
	flag.StringVar(&Conf.FilePath, "path", Conf.FilePath, "Path to files folder")
	flag.BoolVar(&Conf.AllowAdmin, "admin", Conf.AllowAdmin, "Setup server for administration")
	flag.StringVar(&Conf.IpAddress, "ipaddress", Conf.IpAddress, "IP address server to listen at")
	flag.StringVar(&Conf.Port, "port", Conf.Port, "File server port")
	flag.StringVar(&Conf.Realm, "realm", Conf.Realm, "Name of realm")
	configFile := flag.String("config", "", "Name of config file, if given, this overrides all flags")
	flag.Parse()

	if *genPW {
		pass := flag.Arg(0)
		passHash, err := bcrypt.GenerateFromPassword([]byte(pass+"!GoFi"), bcrypt.DefaultCost)
		if err != nil {
			fmt.Errorf("Hash generation failed: %v\n", err)
		}
		fmt.Println(string(passHash))
		return
	}
	if *configFile != "" {
		confFile, err := ioutil.ReadFile(*configFile)
		if err != nil {
			fmt.Printf("Could not open config file %q: %v\n", *configFile, err)
			glog.Errorf("Could not open config file %q: %v\n", *configFile, err)
			return
		}
		err = json.Unmarshal(confFile, &Conf)
		if err != nil {
			fmt.Printf("Config file %q has unexpected format: %v\n", *configFile, err)
			glog.Errorf("Config file %q has unexpected format: %v\n", *configFile, err)
			return
		}
	}

	// Set server
	server := http.Server{
		Addr: Conf.IpAddress + ":" + Conf.Port,
	}
	http.HandleFunc("/home/", authhandler(index))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/home/", http.StatusFound)
	})
	http.HandleFunc("/newDir", authhandler(makeNewDir))
	http.HandleFunc("/serve/", authhandler(forDownload))
	http.HandleFunc("/upload", authhandler(upload))
	http.HandleFunc("/images/", authhandler(images))

	defer glog.Flush()
	glog.Infoln(server.ListenAndServe())
}

func mustRender(w http.ResponseWriter, r *http.Request, name string, data interface{}) {
	err := templates.Render(w, name, data)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Error occurred rendering template:", err)
}

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request:", err)

	return true
}

func badRequestHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/home/", http.StatusFound)
}

func index(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimPrefix(r.URL.Path, "/home/")
	fullPath := joinPath(Conf.FilePath, urlPath)
	glog.Infof("Insert index handler with path=%s", fullPath)

	if r.Method == "POST" {
		glog.Infof("Post request.")
		r.ParseForm()
		if glog.V(2) {
			glog.Infof("PostForm: %v", r.PostForm)
		}
		files := make([]string, len(r.PostForm))
		for key, _ := range r.PostForm {
			deleteFile(fullPath + key)
		}
		glog.Infof("Files: %v", files)
	}

	var di DirInfo
	di.Config = Conf
	stubs := splitString(urlPath, '/')
	di.LocalPath = make([]string, len(stubs))
	for i, _ := range stubs {
		if i == 0 {
			di.LocalPath[i] = stubs[0]
		} else {
			di.LocalPath[i] = joinPath(stubs[0], stubs[1:i+1]...)
		}
	}
	glog.Infof("stubs: %v, LocalPath: %v", stubs, di.LocalPath)
	di.getDir(fullPath)
	mustRender(w, r, "index", di)
}

type DirInfo struct {
	Config    Configuration
	LocalPath []string
	Dirs      []string
	Files     []string
}

func (di *DirInfo) getDir(path string) {
	files, err := ioutil.ReadDir(path)
	path = sanitize(path)
	if err != nil {
		glog.Errorf("Could not read directory: %v", err)
		return
	}

	for _, f := range files {
		name := f.Name()
		if name[0] == '.' {
			continue
		}
		if f.IsDir() {
			di.Dirs = append(di.Dirs, f.Name())
		} else if f.Mode().IsRegular() {
			di.Files = append(di.Files, f.Name())
		}
	}
}

func makeNewDir(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	dirName := r.PostFormValue("dirName")
	urlPath := strings.TrimPrefix(r.PostFormValue("fullpath"), "/home")
	newDirPath := Conf.FilePath + sanitize(urlPath+"/"+dirName)
	err := os.Mkdir(newDirPath, 0770)
	if err != nil {
		glog.Errorf("makeNewDir failed: %v", err)
	}
	glog.Infof("makeNewDir: urlPath='%s', newDirPath='%s'", urlPath, newDirPath)
	http.Redirect(w, r, "/home"+urlPath, http.StatusFound)
}

func deleteFile(path string) {
	path = sanitize(path)
	glog.Infof("Delete File: %s", path)
	os.RemoveAll(path)
}

func forDownload(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimPrefix(r.URL.Path, "/serve/")
	totalPath := Conf.FilePath + "/" + sanitize(urlPath)
	glog.Infof("forDownload: urlPath='%v'", totalPath)
	http.ServeFile(w, r, totalPath)
}

func images(w http.ResponseWriter, r *http.Request) {
	urlPath := strings.TrimPrefix(r.URL.Path, "/images/")
	totalPath := "images" + "/" + sanitize(urlPath)
	glog.Infof("Image: urlPath='%v'", totalPath)
	http.ServeFile(w, r, totalPath)
}

func upload(w http.ResponseWriter, r *http.Request) {
	if glog.V(2) {
		glog.Infof("upload called: mehtod:%s %v", r.Method, *r)
	}
	if r.Method != "POST" {
		http.Redirect(w, r, "/home", http.StatusFound)
		return
	}

	r.ParseMultipartForm(32 << 20)
	relPath := r.PostFormValue("fullpath")
	absPath := Conf.FilePath + sanitize(strings.TrimPrefix(relPath, "/home"))
	fhs := r.MultipartForm.File["uploadfile"]
	for _, fh := range fhs {
		file, err := fh.Open()
		if err != nil {
			glog.Errorf("Upload failed: %v", err)
			return
		}
		defer file.Close()
		if glog.V(2) {
			glog.Infof("Upload Header: %v", fh.Header)
		}
		newFileName := absPath + sanitize(fh.Filename)
		glog.Infof("uplad file with name '%s'", newFileName)
		f, err := os.OpenFile(newFileName, os.O_WRONLY|os.O_CREATE, 0660)
		if err != nil {
			glog.Errorln(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)
	}

	http.Redirect(w, r, relPath, http.StatusFound)
}

func sanitize(s string) string {
	reg, err := regexp.Compile(`[][$&%?@:;#"'´{}]+` + "0x60")
	if err != nil {
		glog.Fatalf("Parsing regexp failed: %v", err)
	}
	r := reg.ReplaceAllString(s, "")

	return r
}

func authhandler(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if ok {
			if user == "" {
				// Readonly user
				err := bcrypt.CompareHashAndPassword([]byte(Conf.UserPWHash), []byte(pass+"!GoFi"))
				if err == nil {
					adminFlag = false
					handler(w, r)
					return
				}
			}
			if user == "admin" {
				// admin user may upload files
				err := bcrypt.CompareHashAndPassword([]byte(Conf.AdminPWHash), []byte(pass+"!GoFi"))
				if err == nil {
					adminFlag = true
					handler(w, r)
					return
				}
			}
		}
		w.Header().Set("www-authenticate", `basic realm="`+Conf.Realm+`"`)
		http.Error(w, "unauthorized.", http.StatusUnauthorized)
	}
}

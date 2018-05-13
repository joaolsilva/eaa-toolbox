package toolbox /* import "r2discover.com/go/eaa-toolbox/pkg/toolbox" */

import (
	"encoding/json"
	"errors"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type actionFromRemote struct {
	Action  string `json:"action"`
	TreePos string `json:"tree_pos"`
}

type webMenuEntry struct {
	TreePos string
	Text    string
}

type templateData struct {
	Position      string
	GoTo          string
	MenuEntries   []webMenuEntry
	HasBackButton bool
}

type restrictedFilesystem struct {
	fs http.FileSystem
}

func (fs restrictedFilesystem) Open(name string) (http.File, error) {
	//log.Printf("restrictedFilesystem: name %s", name)
	if strings.Index(name, "..") != -1 || strings.HasPrefix(name, "/.") || strings.Index(name, "~") != -1 {
		return nil, errors.New("403")
	}
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return restrictedReaddirFile{f}, nil
}

type restrictedReaddirFile struct {
	http.File
}

// Readdir restricts reading directories
func (f restrictedReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

func (toolbox *Toolbox) serve(w http.ResponseWriter, r *http.Request) {
	t := template.New("index.html")
	t, err := t.ParseFiles(filepath.Join(expandHomeDir(toolbox.appConfig.Paths.Web), "template/index.html"))
	if err != nil {
		log.Printf("toolbox.serve: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	data := templateData{}
	data.Position = toolbox.screen.Position
	data.GoTo = toolbox.screen.GoTo

	err = t.Execute(w, data)
	if err != nil {
		log.Printf("toolbox.serve: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (toolbox *Toolbox) serverAction(w http.ResponseWriter, r *http.Request) {
	postData := actionFromRemote{}
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&postData)
}

func (toolbox *Toolbox) startServer() {
	toolbox.serverStarted = true

	r := mux.NewRouter()
	r.HandleFunc("/", toolbox.serve)
	r.HandleFunc("/index.html", toolbox.serve)
	r.HandleFunc("/index.htm", toolbox.serve)
	r.HandleFunc("/actions", toolbox.serverAction).Methods("POST")
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(restrictedFilesystem{http.Dir(filepath.Join(expandHomeDir(toolbox.appConfig.Paths.Web), "static"))})))

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:32243",
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	srv.ListenAndServe()
}

func (toolbox *Toolbox) stopServer() {

}

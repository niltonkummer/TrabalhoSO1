package main

import (
	"fmt"
	"hash/crc32"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"text/template"
)

var Tokens = make(map[uint32]bool)

func setToken(hash uint32, value bool) {
	Tokens[hash] = value
}
func getToken(hash uint32) bool {
	if _, ok := Tokens[hash]; ok {
		return true
	}
	return false
}

const (
	DEFAULT_JQUERY = `<script src="http://code.jquery.com/jquery-1.11.0.min.js"></script>`
)

type HtmlTemplate struct {
	Head, Body string
}

func serve() {
	http.HandleFunc("/admin/", indexAdmin)
	http.HandleFunc("/admin/login/", adminLogin)
	http.HandleFunc("/admin/login/valida/", adminValidaLogin)
	http.HandleFunc("/admin/cadastrar-aluno/", cadastraAluno)
	http.HandleFunc("/admin/cadastrar-disciplina/", cadastraDisciplina)
	http.HandleFunc("/admin/cadastrar-professor/", cadastraProfessor)
	http.HandleFunc("/busca-aluno/", buscaAluno)
	http.ListenAndServe(":8080", nil)
}

func parseDefaultTemplate(w http.ResponseWriter, reader io.Reader) {
	tmpl, err := template.ParseFiles("./files/html/backbone.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := ioutil.ReadAll(reader)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	html_tmpl := HtmlTemplate{DEFAULT_JQUERY, string(data)}
	tmpl.Execute(w, html_tmpl)
}

func adminLogin(w http.ResponseWriter, r *http.Request) {
	fh, err := os.Open("./files/html/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	parseDefaultTemplate(w, fh)
}

func validaLogin(user, password string) bool {
	if user == "root" && password == "root" {
		return true
	}
	return false
}

func adminValidaLogin(w http.ResponseWriter, r *http.Request) {
	user := r.PostFormValue("usuario")
	password := r.PostFormValue("senha")
	ok := validaLogin(user, password)
	if !ok {
		http.Error(w, "Login invalido", http.StatusForbidden)
		return
	}
	hash := crc32.ChecksumIEEE(append([]byte(user), []byte(password)...))
	setToken(hash, true)
	http.Redirect(w, r, "/admin/?c="+fmt.Sprintf("%d", hash), http.StatusFound)
	return
}

func indexAdmin(w http.ResponseWriter, r *http.Request) {
	token := r.FormValue("c")
	hash, _ := strconv.ParseUint(token, 10, 32)
	fmt.Println(hash)
	if !getToken(uint32(hash)) {
		http.Error(w, "Não logado", http.StatusForbidden)
		return
	}
	w.Write([]byte("oi, bem vindo"))
}

func cadastraUsuario(w http.ResponseWriter, r *http.Request) {
	//exec.Command("", )
}

func cadastraDisciplina(w http.ResponseWriter, r *http.Request) {}

func cadastraProfessor(w http.ResponseWriter, r *http.Request) {}

func cadastraAluno(w http.ResponseWriter, r *http.Request) {
	filepath_abs, _ := filepath.Abs("./files/scripts/functions.sh")
	nome := r.FormValue("nome")
	data, err := exec.Command(filepath_abs, "cadastra_aluno", nome).Output()
	if err != nil {
		http.Error(w, string(data)+" "+err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(data)
}

func buscaAluno(w http.ResponseWriter, r *http.Request) {
	filepath_abs, _ := filepath.Abs("./files/scripts/functions.sh")
	data, err := exec.Command(filepath_abs, "busca_aluno", "Nilton").Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Write(data)
}

func main() {
	serve()
}

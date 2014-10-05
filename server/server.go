package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os/exec"
	//"path/filepath"
	"strings"
	"text/template"
	"encoding/json"
)

var SHELL_FUNC_PATH = "./"

const (
	DEFAULT_JQUERY  = `<script src="http://code.jquery.com/jquery-1.11.0.min.js"></script>`
	SHELL_FUNCTIONS = `./files/scripts/functions2.sh`
)

// shell functions
const (
	cadastrar_aluno = "cadastrar_aluno" // "$2" "$3" "$4" "$5"
	// disciplina, nome do aluno
	pesquisar_aluno = "pesquisar_aluno" // "$2"

	cadastrar_disciplina = "cadastrar_disciplina" //"$2"

	listar_disciplinas = "listar_disciplinas"

	listar_alunos_por_disciplina = "listar_alunos_por_disciplina" // "$2"

	verificar_disciplina = "verificar_disciplina"

	// params: disciplina
	sa_salvar = "sa_salvar"

	// Listar BKP
	// params: disciplina
	listar_bkp = "listar_bkp" // $2

	// Repoe um arquivo de turma por um escolhido
	// params: nome_disciplina{data}
	sa_repor = "sa_repor"

	// params: nome_disciplina
	sa_apagar = "sa_apagar"

	compactar = "compactar"
)

type HtmlTemplate struct {
	Head, Body string
}

func serve() {
	fmt.Println(http.FileServer(http.Dir("./")))
	http.Handle("/static/", http.FileServer(http.Dir("files")))
	http.HandleFunc("/", paginaInicial)
	http.HandleFunc("/cadastrar-aluno/", cadastrarAlunoView)
	http.HandleFunc("/listar-alunos/", listarAlunosView)
	http.HandleFunc("/pesquisar-aluno/", pesquisarAlunoView)
	http.HandleFunc("/cadastrar-disciplina/", cadastrarDisciplinaView)
	http.HandleFunc("/salvar-turma/", salvarTurmaView)
	http.HandleFunc("/repor-turma/", reporTurmaView)
	http.HandleFunc("/apagar-turma/", apagarTurmaView)

	http.HandleFunc("/cadastrar-aluno-api/", cadastrarAluno)
	http.HandleFunc("/pesquisar-aluno-api/", pesquisarAluno)
	http.HandleFunc("/listar-alunos-api/", listarAlunos)
	http.HandleFunc("/cadastrar-disciplina-api/", cadastrarDisciplina)
	http.HandleFunc("/salvar-turma-api/", salvarTurma)
	http.HandleFunc("/verificar-turma-api/", verificarTurma)
	http.HandleFunc("/repor-turma-api/", reporTurma)
	http.HandleFunc("/apagar-turma-api/", apagarTurma)
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

type A struct {
	Disciplinas []string
}

func paginaInicial(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./files/html/main.html"))
	t.Execute(w, nil)
}

func cadastrarAlunoView(w http.ResponseWriter, r *http.Request) {
	dis := A{listarDisciplinas()}
	t := template.Must(template.ParseFiles("./files/html/cadastroAluno.html"))
	t.Execute(w, dis)
}

func cadastrarDisciplinaView(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./files/html/cadastrarDisciplina.html"))
	t.Execute(w, nil)
}

func listarAlunosView(w http.ResponseWriter, r *http.Request) {
	dis := A{listarDisciplinas()}
	t := template.Must(template.ParseFiles("./files/html/listarAlunos.html"))
	t.Execute(w, dis)
}

func pesquisarAlunoView(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.ParseFiles("./files/html/pesquisarAluno.html"))
	t.Execute(w, nil)
}

type ListaAlunos struct {
	Disciplina string
	Items      []Linha
}
type Linha struct {
	Cor  string
	Data []string
}

func parseAlunosLista(header, lista []string) ListaAlunos {

	lst := ListaAlunos{"", make([]Linha, len(lista)+1)}
	lst.Items[0] = Linha{"", header}
	for i, registro := range lista {

		item := strings.Split(registro, ":")
		color := "#FFF"
		if i%2 == 0 {
			color = "#EEE"
		}
		items := []string{}
		for _, valor := range item {
			items = append(items, valor)
		}
		lst.Items[i+1] = Linha{color, items}
	}

	return lst
}

func listarAlunos(w http.ResponseWriter, r *http.Request) {
	disciplina := r.FormValue("disciplina")
	data, err := exec.Command(SHELL_FUNCTIONS,
		listar_alunos_por_disciplina,
		disciplina,
	).Output()
	if err != nil {
		return
	}
	slc := strings.Split(string(data), "\n")
	lst := parseAlunosLista([]string{"Matricula", "Nome", "Conceito"}, slc[:len(slc)-1])
	lst.Disciplina = disciplina
	t := template.Must(template.ParseFiles("./files/html/listaDeAlunos.html"))
	t.Execute(w, lst)
}

func listarDisciplinas() []string {
	data, err := exec.Command(SHELL_FUNCTIONS,
		listar_disciplinas,
	).Output()
	if err != nil {
		return nil
	}
	slc := strings.Split(string(data), "\n")
	return slc[:len(slc)-1]
}
func cadastrarAluno(w http.ResponseWriter, r *http.Request) {
	matricula := r.FormValue("matricula")
	nome := r.FormValue("nome")
	conceito := r.FormValue("conceito")
	disciplina := r.FormValue("disciplina")

	data, err := exec.Command(SHELL_FUNCTIONS,
		cadastrar_aluno,
		matricula,
		nome,
		conceito,
		disciplina,
	).Output()
	if err != nil {
		http.Error(w, string(data)+" "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(data))
	w.Write([]byte(data))
}

func pesquisarAluno(w http.ResponseWriter, r *http.Request) {
	nome := r.FormValue("nome")
	data, err := exec.Command(
		SHELL_FUNCTIONS,
		pesquisar_aluno,
		nome).Output()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slc := strings.Split(string(data), "\n")
	lst := parseAlunosLista([]string{"Disciplina", "Matricula", "Nome", "Conceito"}, slc[:len(slc)-1])
	t := template.Must(template.ParseFiles("./files/html/listaDeAlunos.html"))
	t.Execute(w, lst)
}

func cadastrarDisciplina(w http.ResponseWriter, r *http.Request) {
	nome := r.FormValue("nome")

	data, err := exec.Command(SHELL_FUNCTIONS,
		cadastrar_disciplina,
		nome).Output()
	if err != nil {
		http.Error(w, string(data)+" "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(data))
	w.Write([]byte(data))
}

func salvarTurmaView(w http.ResponseWriter, r *http.Request) {
	dis := A{listarDisciplinas()}
	t := template.Must(template.ParseFiles("./files/html/salvarTurma.html"))
	t.Execute(w, dis)
}

func salvarTurma(w http.ResponseWriter, r *http.Request) {
	nome := r.FormValue("disciplina")
	data, err := exec.Command(SHELL_FUNCTIONS,
		sa_salvar,
		nome).Output()
	if err != nil {
		http.Error(w, string(data)+" "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(data))
	w.Write([]byte(data))
}

func verificarTurma(w http.ResponseWriter, r *http.Request) {
	nome := r.FormValue("disciplina")
	data, err := exec.Command(SHELL_FUNCTIONS,
		verificar_disciplina,
		nome).Output()
	if err != nil {
		http.Error(w, string(data)+" "+err.Error(), http.StatusBadRequest)
		return
	}
	var ret map[string]bool
	if strings.TrimSpace(string(data)) == "1" {
		ret = map[string]bool{"existe": true}
	} else {
		ret = map[string]bool{"existe": false}
	}
	data, _ = json.Marshal(ret)
	w.Write([]byte(data))
}

func reporTurmaView(w http.ResponseWriter, r *http.Request) {
	dis := A{listarDisciplinas()}
	t := template.Must(template.ParseFiles("./files/html/reporTurma.html"))
	t.Execute(w, dis)
}

func reporTurma(w http.ResponseWriter, r *http.Request) {
	nome := r.FormValue("disciplina")
	data, err := exec.Command(SHELL_FUNCTIONS,
		sa_repor,
		nome).Output()
	if err != nil {
		http.Error(w, string(data)+" "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(data))
	w.Write([]byte(data))
}

func apagarTurmaView(w http.ResponseWriter, r *http.Request) {
	dis := A{listarDisciplinas()}
	t := template.Must(template.ParseFiles("./files/html/apagarTurma.html"))
	t.Execute(w, dis)
}

func apagarTurma(w http.ResponseWriter, r *http.Request) {
	nome := r.FormValue("disciplina")
	data, err := exec.Command(SHELL_FUNCTIONS,
		sa_apagar,
		nome).Output()
	if err != nil {
		http.Error(w, string(data)+" "+err.Error(), http.StatusBadRequest)
		return
	}
	fmt.Println(string(data))
	w.Write([]byte(data))
}


func main() {
	// SHELL_FUNC_PATH, _ = filepath.Abs(SHELL_FUNCTIONS)
	serve()
}

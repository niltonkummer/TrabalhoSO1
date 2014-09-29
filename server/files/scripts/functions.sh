#!/bin/bash

# $0, 
# $1 = função
# $2 = parametros
SEPARADOR=:
pwd=$(pwd)
DIRETORIOSALVO="$pwd/files/db"
DIRETORIO_ALUNOS="$DIRETORIOSALVO/alunos"
DIRETORIO_DISCIPLINAS="$DIRETORIOSALVO/disciplinas"
DIRETORIO_PROFESSORES="$DIRETORIOSALVO/professores"
DIRETORIO_TURMAS="$DIRETORIOSALVO/turmas"
DIRETORIO_AVALIACAO="$DIRETORIO_TURMAS$DIRETORIOAVALIACAO/"
DIRETORIO_ID_MATRICULA="$DIRETORIOSALVO/matriculas_sequencia"
DIRETORIO_ID_DISCIPLINA="$DIRETORIOSALVO/disciplinas_sequencia"
DIRETORIO_ID_TURMA="$DIRETORIOSALVO/turmas_sequencia"
DIRETORIO_PROFESSOR_TURMA="$DIRETORIOSALVO/professor_turma"

touch "$DIRETORIO_ID_MATRICULA"
touch "$DIRETORIO_ID_DISCIPLINA"
touch "$DIRETORIO_DISCIPLINAS"
touch "$DIRETORIO_ID_TURMA"

# id_disciplina/turmas/alunos
# id_disciplina/turmas/professor
# Espanhol/A/alunos_matriculados professor aval_1 aval_1

# Um professor é responsavel por uma ou varias turmas
# Uma turma possui apenas um professor
# Uma disciplina possui varios professores
# Um professor possui varias disciplinas

# alunos_matriculados
# 1
# 2
# 3

# aval_1
# 1:B
# 2:C
#


# TODO Criar diretorios que não existem
function verifica_id_aluno() {
	e=$(cat "$DIRETORIO_ALUNOS" | cut -d $SEPARADOR -f 1 | grep -w $1 | wc -l)
	return $e
}

function pega_id_aluno() {
	e=$(echo $1 | cut -d $SEPARADOR -f 1)
	return $e
}

function gerar_id() {
	id_atual=$(cat "$1" | wc -l)
	if test $id_atual -eq 0
	then
		echo 100 > "$1"
		# Usuário já cadastrado
		return 100
	fi
	
	identificacao=$(cat "$1")
	identificacao=$(expr $identificacao + 1)
	echo $identificacao > "$1"
	return $identificacao
}

# Verifica se existe o aluno e o insere no arquivo
# retorna 0 para erro ou 1 para inserido
# params: nome_do_aluno
function cadastrar_aluno() {
	gerar_id "$DIRETORIO_ID_MATRICULA"
	id_aluno=$?
	# insere o usuario no arquivo de alunos
	$(echo "$id_aluno:$1" >> "$DIRETORIO_ALUNOS")
	return 0
}

# Busca o usuario pelo nome
function buscar_alunos() {
	linhas=$(cat "$DIRETORIO_ALUNOS" | cut -d $SEPARADOR -f 2 | grep -n "$1" | cut -d $SEPARADOR -f 1)
	for linha in $linhas
	do
		# Le a linha N e somente isso
		head -n $linha "$DIRETORIO_ALUNOS" | tail -1 
	done
}

# Verifica pelo nome da disciplina se a mesma já existe
function verificar_disciplina() {
	e=$(cat "$DIRETORIO_DISCIPLINAS" | cut -d $SEPARADOR -f 2 | grep -w $1 | wc -l)
	return $e
}

# cria estrutura disciplina
function cadastrar_disciplina() {
	verificar_disciplina $1
	ret=$?
	if test $ret -eq 1
	then
		# Disciplina já cadastrada		
		return 1
	fi
	gerar_id "$DIRETORIO_ID_DISCIPLINA"
	id_disciplina=$?
	# insere a disciplina no arquivo de disciplina
	$(echo "$id_disciplina:$1" >> "$DIRETORIO_DISCIPLINAS")
	return 0
}

## FUNCOES DE TURMA
# cadastra uma nova turma
# params: nome_disciplina, id_professor
function cadastrar_turma(){
	# gerar turma
	if ! test $2  
	then
		return 1
	fi
	validar_professor $2
	valido=$?
	if test $valido -eq 0  
	then
		return 1
	fi
	gerar_id "$DIRETORIO_ID_TURMA"
	nome_turma=$?
	mkdir -p "$DIRETORIO_TURMAS/$1/$nome_turma"
	echo "$nome_turma:$2" >> "$DIRETORIO_PROFESSOR_TURMA"
}

function gerar_turma() {
	echo ""
}

# Matricula um aluno em uma turma
# params: id_aluno,id_disciplina
function matricular_aluno() {
	# ID_matricula
	echo ""
}

function validar_professor() {
	if ! test -f "$DIRETORIO_PROFESSORES" 
	then
		# Professor não existe
		return 0
	fi
	
	id=$(cat "$DIRETORIO_PROFESSORES" |  cut -d $SEPARADOR -f 1 | grep -w $1)
	
	if  [ "$id" != "" ]  && test $id -eq $1 
	then
		return 1
	fi
	return 0
	
}

function verifica_turma() {
	grep "$DIRETORIO_TURMAS"
}

# Seleciona a professora pelo login

# salva o arquivo, caso o comando seja executado
# várias vezes com o mesmo arquivo serão salvas várias versões
# do arquivo
function sa_salvar() {
	echo "salvar"
}

# copia a versão mais recente do arquivo para o
# diretório corrente, pede confirmação caso já exista o arquivo no
# diretório corrente.
function sa_repor() {
	echo "repor"
}

# [-f] arquivo: apagar o arquivo e todas as versões
# existentes, caso a opção -f for passada pede confirmação
function sa_apagar() {
	echo "apagar";
}

# cria um arquivo de compactação com arquivo
# e apaga todas as versões do arquivo. É perguntado para o usuário
# qual o tipo de compactação deseja-se fazer. 
function sa_compactar() {
	echo "compactar";
}

function cadastrar_professores_padrao() {
	echo "1:Roberto" >> "$DIRETORIO_PROFESSORES"
	echo "2:Jose" >> "$DIRETORIO_PROFESSORES"
	echo "3:Alfredo" >> "$DIRETORIO_PROFESSORES"
}

case $1 in 
	"cadastrar_aluno")
	cadastrar_aluno $2
	;;
	"buscar_alunos")
	buscar_alunos $2
	;;
	"cadastrar_disciplina")
	cadastrar_disciplina $2
	;;
	"cadastrar_professores_p")
	cadastrar_professores_padrao
	;;
	"cadastrar_turma")
	cadastrar_turma $2 $3
	;;
	"salvar")
	sa_salvar
	;;
	"repor")
	;;
	"apagar")
	;;
	"compactar")
	;;
esac
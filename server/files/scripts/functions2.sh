#!/bin/bash

# $0, 
# $1 = função
# $2 = parametros
SEPARADOR=:
pwd=$(pwd)

DIRETORIOSALVO="./files/db"
DIRETORIO_DISCIPLINAS="$DIRETORIOSALVO/disciplinas"
# Disciplinas são pastas
# Alunos são registros em arquivos dentro de Disciplinas
function verifica_id_aluno() {
	e=$(cat "$DIRETORIO_DISCIPLINAS/$2/turma" | cut -d $SEPARADOR -f 1 | grep -w $1 | wc -l)
	return $e
}

function cadastrar_aluno() {
	if ! test -d "$DIRETORIO_DISCIPLINAS/$4/"
	then
		# disciplina não cadastrada
		return 1
	fi

	verifica_id_aluno "$1" "$4" 
	cadastrado=$?
	if  test $cadastrado -eq 1
	then
		# disciplina não cadastrada
		echo "Aluno já está cadastrado"
		return 0
	fi
	echo "$1$SEPARADOR$2$SEPARADOR$3"
	echo "$1$SEPARADOR$2$SEPARADOR$3" >> "$DIRETORIO_DISCIPLINAS/$4/turma"
	echo "Aluno cadastrado"
}

function pesquisar_aluno() {
	# Busca os arquivos que contem o aluno com o nome procurado
	dir_path="$DIRETORIO_DISCIPLINAS/*/turma"
	IFS=$(echo -en "\n\b")
	dir_turmas=$(grep $1 -l $dir_path)
	for dir_turma in $dir_turmas
	do
		if ! test -f "$dir_turma"
		then
			# disciplina não cadastrada
			continue
		fi
		aluno=$(cat "$dir_turma" | grep $1)
		turma=$(echo "$dir_turma" | cut -d "/" -f 5)
		echo $turma:$aluno 
	done 
	return

	linhas=$(cat "$DIRETORIO_DISCIPLINAS/$1/turma" | cut -d $SEPARADOR -f 2 | grep -n $2 | cut -d $SEPARADOR -f 1)
	for linha in $linhas
	do
		echo $linha
	done 
}

function listar_alunos_por_disciplina() {
	if ! test -f "$DIRETORIO_DISCIPLINAS/$1/turma"
	then
		# disciplina não cadastrada
		return 1
	fi
	cat "$DIRETORIO_DISCIPLINAS/$1/turma"
}

function cadastrar_disciplina() {
	if test -d "$DIRETORIO_DISCIPLINAS/$1"
	then 
		echo "Disciplina já cadastrada"
		return
	fi	
	mkdir -p "$DIRETORIO_DISCIPLINAS/$1"
	echo "Disciplina cadastrada"
}

function listar_disciplinas() {
	#echo "$DIRETORIO_DISCIPLINAS/"
	if ! test -d "$DIRETORIO_DISCIPLINAS/"
	then 
		#echo "Nenhuma disciplina"
		return
	fi
	cd "$DIRETORIO_DISCIPLINAS"
	for dir in *
	do
		echo "$dir"
	done 
}


# params: disciplina
function sa_salvar() {
	if ! test -d "$DIRETORIO_DISCIPLINAS/$1/bkp"
	then
		mkdir -p "$DIRETORIO_DISCIPLINAS/$1/bkp"
	fi
	cp "$DIRETORIO_DISCIPLINAS/$1/turma" "$DIRETORIO_DISCIPLINAS/$1/bkp/"
	time=$(date +%Y-%m-%d:%H:%M:%S)
	mv "$DIRETORIO_DISCIPLINAS/$1/bkp/turma" "$DIRETORIO_DISCIPLINAS/$1/bkp/turma_$time"
}

# Verifica se o arquivo de disciplina existe no diretorio
function verificar_disciplina() {
	if ! test -f "$DIRETORIO_DISCIPLINAS/$1/turma"
	then
		echo 0
		return 0
	fi
	echo 1
}

# Repoe um arquivo de uma determinada disciplina
# params: disciplina
function sa_repor() {
	file=$(ls -t "$DIRETORIO_DISCIPLINAS/$1/bkp" | head -1)
	cp -f "$DIRETORIO_DISCIPLINAS/$1/bkp/$file" "$DIRETORIO_DISCIPLINAS/$1/turma"
}

function listar_disciplina_bkp() {
	#echo "$DIRETORIO_DISCIPLINAS/"
	if ! test -d "$DIRETORIO_DISCIPLINAS/$1/bkp"
	then 
		#echo "Nenhuma disciplina"
		return
	fi
	cd "$DIRETORIO_DISCIPLINAS/$1/bkp"
	IFS=$(echo -en "\n\b")
	for dir in `ls -t`
	do
		echo "$dir"
	done
}

# params: disciplina, arquivo
function sa_recuperar() {
	cp -f "$DIRETORIO_DISCIPLINAS/$1/bkp/$file" "$DIRETORIO_DISCIPLINAS/$1/turma"
}

function sa_apagar() {
	rm -rf "$DIRETORIO_DISCIPLINAS/$1/"
}

# Verifica a função que vai executar
case $1 in 
	"cadastrar_aluno")
 	# Mat, Nome, Conceito, Disciplina
	cadastrar_aluno "$2" "$3" "$4" "$5"
	;;
	# disciplina, nome do aluno
	"pesquisar_aluno")
	pesquisar_aluno "$2"
	;;
	"cadastrar_disciplina")
	cadastrar_disciplina "$2"
	;;
	"listar_disciplinas")
	listar_disciplinas 
	;;
	"listar_disciplina_bkp")
	listar_disciplina_bkp "$2"
	;;
	"listar_alunos_por_disciplina")
	listar_alunos_por_disciplina "$2"
	;;
	"verificar_disciplina")
	verificar_disciplina "$2"
	;;
	"sa_salvar")
	sa_salvar "$2"
	;;
	"sa_repor")
	sa_repor "$2"
	;;
	"sa_apagar")
	sa_apagar "$2"
	;;
	"compactar")
	;;
esac
#!/usr/bin/env bash

function showAll() {
    git branch -a
}

function createAndTrack() {
    git checkout -B $1 --track $1
}

function changeBranch() {
    git reset --hard
    git checkout -B $1
    git pull
}

function generatorData() {
    cd tools/gen_const_sheet
    ./gen_const_sheet.sh
    cd -
}

function commit() {
    git status
    git commit -a -m "gen plist data"
}

function push() {
    git push --all
    exit 0
}

function error() {
    echo "Usage: git.sh {git-path} {all|checkout|generator|commit|push} {name}"
    exit
}

if [[ -z "$1" ]] || [[ -z "$2" ]]; then
	error
fi
cd "$1"
case "$2" in
	"all")
        showAll
		;;
	"checkout")
	    if [[ -z "$3" ]]; then
	        error
	    fi
		changeBranch $3
		;;
	"generator")
        generatorData
	    ;;
    "commit")
        commit
	    ;;
    "push")
        push
	    ;;
	*)
		error
		;;
esac




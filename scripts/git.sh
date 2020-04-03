#!/usr/bin/env bash

function fetch() {
    git fetch --all
    git fetch -p
}

function showAll() {
    git branch -a | grep -v HEAD
}

function revert() {
    git add --all
    git checkout -f
    git reset --hard
}

function checkout() {
    revert
    git checkout -B $1 --track $2
    exit 0
}

function generateData() {
    git pull
    cd tools/gen_const_sheet
    ./gen_const_sheet.sh
    cd -
    git add --all
}

function commit() {
    git status
    git commit -a -m "gen plist data by go-gpt"
}

function push() {
    git pull
    git push --all
}

function update() {
    ./deploy_dev.sh
    exit 0
}

function svnSync() {
    ./deploy_svn.sh $1 $2
    exit 0
}

function compress() {
    ./zip.sh "$1" "$2" "$3"
    ls | grep '.zip\|.txt'
    exit 0
}

function error() {
    echo "Usage: git.sh {git-path} {fetch|revert|showAll|checkout|generate|commit|push|update|svnSync|compress} {name}"
    exit
}

if [[ -z "$1" ]] || [[ -z "$2" ]]; then
	error
fi
cd "$1"
case "$2" in
    "fetch")
        fetch
        ;;
     "revert")
        revert
        ;;
    "showAll")
        showAll
        ;;
    "checkout")
        if [[ -z "$3" ]] || [[ -z "$4" ]]; then
            error
        fi
        checkout $3 $4
        ;;
    "generate")
        generateData
        ;;
    "commit")
        commit
        ;;
    "push")
        push
        ;;
    "update")
        update
        ;;
    "svnSync")
        if [[ -z "$3" ]] || [[ -z "$4" ]]; then
            error
        fi
        svnSync $3 $4
        ;;
    "compress")
        if [[ -z "$3" ]] || [[ -z "$4" ]] || [[ -z "$5" ]]; then
            error
        fi
        compress "$3" "$4" "$5"
        ;;
    *)
        error
        ;;
esac




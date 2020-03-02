#!/usr/bin/env bash

function checkout() {
    svn --username $1 --password $2 checkout $3
}

function version() {
    svn --username $1 --password $2
}

function add() {
    svn --username $1 --password $2 add $1
}

function status() {
    svn --username $1 --password $2 status
}

function commit() {
    svn --username $1 --password $2 commit --message "${3}commit by ${1}@go-gpt"
}

function update() {
    svn --username $1 --password $2 update
}

function info() {
    svn --username $1 --password $2 info
}

function log() {
    svn --username $1 --password $2 log -l 1 -v
}

function lock() {
    svn --username $1 --password $2 lock
}

function unlock() {
    svn --username $1 --password $2 unlock
}

function error() {
    echo "Usage: svn.sh {username} {password} {git-path} {checkout|add|status|commit|update|log} {message}"
    exit
}

if [[ -z "$1" ]] || [[ -z "$2" ]] || [[ -z "$3" ]] || [[ -z "$4" ]] ; then
	error
fi
cd "$3"
case "$4" in
    "checkout")
        checkout $1 $2
        ;;
    "add")
        add $1 $2
        ;;
    "status")
        status $1 $2
        ;;
    "commit")
        if [[ -z "$5" ]]; then
            error
        fi
        commit $1 $2 $5
        ;;
    "update")
        update $1 $2
        ;;
    "log")
        log $1 $2
        ;;
    *)
        error
        ;;
esac
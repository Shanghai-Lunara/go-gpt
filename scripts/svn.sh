#!/usr/bin/env bash

function checkout() {
    svn --username $1 --password $2 checkout $3
}

function add() {
    svn --username $1 --password $2 add $3
}

function addAll() {
    svn --username $1 --password $2 status | grep ? | awk '{print $2}' | xargs svn --username $1 --password $2 add
}

function status() {
    svn --username $1 --password $2 status
}

function revert() {
    svn --username $1 --password $2 revert $3
}

function revertAll() {
    svn --username $1 --password $2 status | awk '{print $2}' | xargs svn --username $1 --password $2 revert --depth infinity
}

function removeAll() {
    svn --username $1 --password $2 status | grep ? | awk '{print $2}' | xargs rm -rf
}

function clean() {
    revertAll $1 $2
    removeAll $1 $2
}

function commit() {
    addAll $1 $2
    svn --username $1 --password $2 commit --message "${3} committed by ${1}@go-gpt"
}

function update() {
    clean $1 $2
    svn --username $1 --password $2 update
}

function info() {
    svn --username $1 --password $2 info --xml
}

function log() {
    svn --username $1 --password $2 log -l $3 -v --xml
}

function lock() {
    svn --username $1 --password $2 lock
}

function unlock() {
    svn --username $1 --password $2 unlock
}

function error() {
  cat <<EOF
Usage: $(basename "$0") <username> <password> <svn-path> <svn_dir> <command> ...

  <username>          the svn account's username
  <password>          the svn account's password
  <svn-path>          the svn remote url
  <svn_dir>           the directory of local
  <command>           the commands (e.g. checkout|addAll||status|revertAll|removeAll|clean|commit|update)
  <params>            the external param for commands

Examples:
  $(basename "$0") {username} {password} {svn-path} {svn_dir} {checkout} {checkout-url}"
EOF
    exit
}

if [[ -z "$1" ]] || [[ -z "$2" ]] || [[ -z "$3" ]] || [[ -z "$4" ]] || [[ -z "$5" ]]; then
	error
fi
cd "$3"
case "$5" in
    "checkout")
        if [[ -z "$6" ]]; then
            error
        fi
        checkout $1 $2 $6
        ;;
    "addAll")
        if [[ -z "$6" ]]; then
            error
        fi
        cd $4
        addAll $1 $2
        ;;
    "status")
        cd $4
        status $1 $2
        ;;
    "clean")
        cd $4
        clean $1 $2
        ;;
    "commit")
        if [[ -z "$6" ]]; then
            error
        fi
        cd $4
        commit $1 $2 $6
        ;;
    "update")
        cd $4
        update $1 $2
        ;;
    "log")
        if [[ -z "$6" ]]; then
            error
        fi
        cd $4
        log $1 $2 $6
        ;;
    *)
        error
        ;;
esac
#!/usr/bin/env bash

function showAll() {
    git branch -a
}

function changeBranch() {
    git checkout $1
    git pull
}

if [[ -z "$1" ]] || [[ -z "$2" ]]; then
	echo "Usage: git.sh {all|change} version"
	exit
fi
cd "$1"
case "$2" in
	"all")
        showAll
		;;
	"change")
		changeBranch $3
		;;
	*)
		echo "Usage: git.sh {all|change} version"
		exit
		;;
esac




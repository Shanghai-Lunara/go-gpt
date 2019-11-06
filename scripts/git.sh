#!/usr/bin/env bash

function showAll() {
    git branch -a
}

function showActive() {
    git branch -a | grep \*
}

if [[ -z "$1" ]] || [[ -z "$2" ]]; then
	echo "Usage: git.sh {all|act} version"
	exit
fi
cd "$1"
case "$2" in
	"all")
        showAll $1 $2
		;;
	"act")
		showActive
		;;
	"check")
	    ;;
	*)
		echo "Usage: git.sh {all|act} version"
		exit
		;;
esac




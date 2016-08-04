#!/usr/bin/env bash

build_bin=false
timer=false
comperssion=false
keep_bin=false

Usage() {
    echo "${0##\*/} [-tbck] Runs all example puzzles in succession"
    echo "      -t : Include timer output"
    echo "      -b : Compile source first (use go build, then run from binary)"
    echo "      -c : When compiling source, use compression flag (-ldflags '-s')"
    echo "      -k : Keep binary"
    echo "      -h : Show this help message"
    exit "${1-0}"
}

while getopts :tbckh opt; do
    case $opt in
        t)
            timer=true
            ;;
        b)
            build_bin=true
            ;;
        c)
            compression=true
            ;;
        k)
            keep_bin=true
            ;;
        h)
            Usage 0
            ;;
        \?)
            echo "Unknown flag -${OPTARG}" >&2
            Usage 1
            ;;
    esac
done
if [ "$build_bin" = true ]; then
    if [ "$compression" = true ]; then
        go build -ldflags "-s" sudoku.go
    else
        go build sudoku.go
    fi
    for f in $(ls *.input); do
        echo "input file: ${f}"
        if [ "$timer" = true ]; then
            ./sudoku -time -path "$f"
        else
            ./sudoku -path "$f"
        fi
        echo
    done
    if [ "$keep_bin" = false ]; then
        rm sudoku -f
    fi
else
    if [ "$timer" = true ]; then
        ls *.input | xargs -n 1 go run sudoku.go -time -path
    else
       ls *.input | xargs -n 1 go run sudoku.go -path
    fi
fi



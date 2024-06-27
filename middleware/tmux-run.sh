#!/bin/env bash

# run env-setter.sh, if there's error in this script, notify the user
source env-setter.sh || { echo "Error in env-setter.sh"; exit 1; }

# start new tmux session
tmux new-session -d -s asmr

# split window into two horizontal panes then two vertical panes
tmux split-window -v
tmux split-window -h

# run simulator in first pane
tmux send-keys -t 0 'cd sim' C-m
tmux send-keys -t 0 'go run main.go' C-m

# run rest server in second pane
tmux send-keys -t 1 'cd rest_server' C-m
tmux send-keys -t 1 'go run main.go' C-m

# run rule engine in third pane
tmux send-keys -t 2 'cd rule_engine' C-m
tmux send-keys -t 2 'go run main.go' C-m

# attach to tmux session
tmux attach-session -t asmr

if [ -z "$IN_NIX_SHELL" ]; then
    exec nix-shell --run "$0"
fi

if tmux has-session -t vox 2>/dev/null; then
    tmux attach -t vox
else
    tmux new-session -ds vox -n make
    tmux new-window -t vox -n app
    tmux send-keys -t vox:app "cd app/; clear" C-m
    tmux attach -t vox
fi

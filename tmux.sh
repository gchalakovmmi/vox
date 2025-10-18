if tmux has-session -t vox 2>/dev/null; then
    echo "Session 'vox' already exists. Attaching..."
    tmux attach -t vox
else
    echo "Creating new session 'vox'..."
    tmux new-session -ds vox -n root
    tmux send-keys -t vox:root "echo 'Root dir of the project'" C-m
    tmux new-window -t vox -n app
    tmux set-option -t vox status-style fg=white,bg=black
    tmux attach -t vox
fi

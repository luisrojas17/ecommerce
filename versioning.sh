#!/bin/bash

# This file wraps all operations like git add, git commit, git push in order to send 
# all local changes from local repository to remote repository.

git add .
git commit -m $1
git push
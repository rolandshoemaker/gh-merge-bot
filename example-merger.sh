#!/bin/bash

# DO SOMETHING MERGE LIKE

# assumes git config has already been setup for the user executing
# this script (name, email, signing key)
local_user="rolandshoemaker"
merge_flags="--no-ff -S"

git fetch origin
git checkout master

if ["$MERGE_USER" = "$local_user"]; then
    git merge $merge_flags $MERGE_BRANCH
else
    git checkout -b $MERGE_USER-$MERGE_BRANCH master
    git pull https://github.com/$MERGE_USER/repo.git $MERGE_BRANCH
    git checkout master
    git merge $merge_flags $MERGE_USER-$MERGE_BRANCH
fi
git push origin master

# `merge-bot`

This is a pretty simple GitHub bot that can be called upon to merge pull requests
in strange fashions. Before executing the merge it will check that the PR was reviewed
the required number of times, has all successful commit statuses, and that it is up to
date with master.

Instead of trying to merge things itself (because eww) it will run a configured command
(e.g. `bash example-merger.sh`) which can preform the merge. `merge-bot` passes state to
this command via environment variables (TBD) so it knows what it is actually trying to
merge. The configured command should exit with either 0 indicating success or 1 indicating
failure (in the later case the command should also output a message on STDERR that can be
logged or used in a comment on the PR) so `merge-bot` actually knows what happened during
the merge.

# Work flow

Basically how `merge-bot` works

1. Someone makes a comment on a PR/Issue
2. `merge-bot` receives the issue webhook payload, verifies it is from GH, and
   1. checks the comment is on a PR
   2. checks that the body of the comment is '@merge-bot-username merge'
   3. checks the commenter is a configured reviewer
   4. checks the required number of review labels are present (incl. checking per-label overrides)
   5. checks the PR is against master and that it is up to date
   6. checks the PR has no conflicts
   7. checks the combined PR statuses == 'successful'
3. if any of the above checks fail log (and if any are important/addressable comment on PR)
4. merge script is executed with env vars set
   1. if merge script was successful do nothing(?)
   2. if merge script fails comment on PR with STDERR(?)

#!/usr/bin/env bash
# Credit: Skycoin Skywire project
# Reference: https://github.com/skycoin/skywire/blob/master/static/script/unix/check

cd ${GOPATH}/src/github.com/BigOokie/skywire-wing-commander
git checkout master > /dev/null 2>&1
[ $(git rev-parse HEAD) = $(git ls-remote $(git rev-parse --abbrev-ref @{u} | \
sed 's/\// /g') | cut -f1) ] && echo "false" || echo "true"
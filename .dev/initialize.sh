#!/usr/bin/env bash

echo "🎒 check available golangci-lint..."
if ! command -v golangci-lint &> /dev/null
then
    echo "command golangci-lint not exist, process to install"
    curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.39.0
fi
echo "👍 golangci-lint [OK]"

echo "📮 check available go test..."
if ! command -v gotest &> /dev/null
then
    go get -u github.com/rakyll/gotest
fi
echo "👍 gotest [OK]"

echo "🎓 check available gqlgen..."
if ! command -v gqlgen &> /dev/null
then
    go install github.com/99designs/gqlgen
fi
echo "👍 gqlgen [OK]"

echo "📡 initialize GIT hooks..."
chmod +x ./.github/pre-commit.sh
cp ./.github/pre-commit.sh ./.git/hooks/pre-commit
chmod +x ./.github/pre-push.sh
cp ./.github/pre-push.sh ./.git/hooks/pre-push
git config commit.template ./.github/git-commit-template
echo "👍 GIT pre commit/push and template configured [OK]"

echo "🎉 Great! you all set, next step run this command \"make services-up\". Happy Coding 😆"




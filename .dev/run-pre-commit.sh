#!/usr/bin/env sh

if [ -n "$SKIP_PRE_COMMIT" ]; then
  echo "✅ Skipping pre-commit because env var SKIP_PRE_COMMIT exists and not-empty"
  exit 0
fi

CHANGED_GO_FILES=$(git diff HEAD --name-only | egrep '\.go$')

if [ -z "$CHANGED_GO_FILES" ]; then

  echo "✅ No golang files changed"

else
  echo "🔎 Check Buildable"
  if ! make check-buildable; then
    echo "⛔ Code not buildable"
    exit 1
  fi
  echo "✅ Build OK"

  echo "🔎 Linting"
  if ! make lint; then
    echo "⛔ Code linting failed"
    exit 1
  fi
  echo "✅ Lint OK"

  echo "🔎 Golang Imports"
  if ! make check-imports-newline; then
    echo "⛔ Found extra new lines in golang imports"
    exit 1
  fi
  echo "✅ Golang imports OK"

  echo "🔎 Run Unit Test"
  if ! make test; then
    echo "⛔ Unit test failed"
    exit 1
  fi
  echo "✅ Unit Test OK"

  echo "🔎 Run Save Unit Test Coverage"
  if ! make test-coverage; then
    echo "⛔ Unit test failed"
    exit 1
  fi
  echo "🔎 add converage.out in commit"
  git add coverage.out
  echo "✅ Test Coverage Saved"

  echo "🔎 Run Unit Test Coverage"
  coverage=$(go tool cover -func coverage.out | grep total | awk '{print substr($3, 1, length($3)-1)}')
  min=65.0
  if (( ${coverage%%.*} < ${min%%.*} || ( ${coverage%%.*} == ${min%%.*} && ${coverage##*.} < ${min##*.} ) )) ; then
    echo "⛔ Unit Test Coverage ${coverage} < ${min}"
    exit 1
  fi
  echo "✅ Unit Test Coverage OK ${coverage} > ${min}"

fi
echo "✅ Pre-Commit OK"
exit 0

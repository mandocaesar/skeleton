#!/usr/bin/env sh

has_error=0
for fn in $( git diff HEAD --name-only | grep -P '.+\.go$$' )
do
  stat "$fn" 1> /dev/null 2>/dev/null
  if [ $? -eq 0 ]; then
    .dev/check-newline-in-imports -f="$fn"
    if [ $? -ne 0 ]; then
        has_error=1
    fi
  fi
done
exit $has_error
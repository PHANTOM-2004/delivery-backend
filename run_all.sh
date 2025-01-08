#!/bin/sh

run() {
  cd "$1" && air
}

root_path="$(pwd)"

echo "root: $root_path"

for file in ./service/*; do
  if [ -e "$file/main.go" ]; then
    echo "find main.go in $file, it is a service"
    cd "$file" && air &
    cd "$root_path" || exit
  fi
done

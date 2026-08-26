#!/usr/bin/env bash
set -euo pipefail

if [[ $# -ne 1 || -z "$1" ]]; then
  echo "usage: bash scripts/rename.sh <new-module-name>"
  exit 1
fi

if [[ ! -f go.mod ]]; then
  echo "error: run this script from the project root"
  exit 1
fi

NEW_MODULE=$1
if [[ ! "$NEW_MODULE" =~ ^[[:alnum:].~_-]+(/[[:alnum:].~_-]+)*$ ]]; then
  echo "error: invalid Go module name: $NEW_MODULE"
  exit 1
fi

OLD_MODULE=$(go list -m)
if [[ "$OLD_MODULE" == "$NEW_MODULE" ]]; then
  echo "Module is already named $NEW_MODULE"
  exit 0
fi

echo "Renaming module: $OLD_MODULE -> $NEW_MODULE"

OLD_PATTERN=${OLD_MODULE//./\\.}

while IFS= read -r -d '' file; do
  if grep -qF "$OLD_MODULE/" "$file"; then
    sed "s|$OLD_PATTERN/|$NEW_MODULE/|g" "$file" > "$file.rename-tmp"
    mv "$file.rename-tmp" "$file"
  fi
done < <(find . -type f -name '*.go' ! -path './.git/*' ! -path './vendor/*' -print0)

go mod edit -module="$NEW_MODULE"
go fmt ./...
go mod tidy
go test ./...
go build ./...

echo "Module renamed successfully."

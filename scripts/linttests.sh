#!/bin/bash

#!/bin/bash

# Проверяем наличие параметра -type
echo "Usage: $0 -type $2 $3"

if [ "$2" = "staticlint" ]
then
    go run ./cmd/staticlint ./...
elif [ "$2" = "compare" ]
then
    go run -ldflags "-X main.buildVersion=v1.0.1" ./cmd/server ./...
fi




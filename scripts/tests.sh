#!/bin/bash

# Проверяем наличие параметра -type
echo "Usage: $0 -type $2 $3"
# if [ -z "$1" ]; then
#   echo "Usage: $0 -type $2"
#   exit 1
# fi

if [ "$2" = "generate" ]
then
    curl http://127.0.0.1:8080/debug/pprof/heap?seconds=120 -o ./profiles/mem_out.pprof
    curl http://127.0.0.1:8080/debug/pprof/profile?seconds=30 -o ./profiles/cpu_out.pprof
elif [ "$2" = "compare" ]
then
    pprof -top -diff_base=profiles/base_cpu_out.pprof profiles/cpu_out.pprof
elif [ "$2" = "tests" ]
then
    find ${WORK_DIR}/internal/server -type d -exec sh -c 'go test -bench=. -cpuprofile=profiles/cpu_$(basename {}).out {}' \;
elif [ "$2" = "trun" ]
then
    go tool pprof -http=":9090" profiles/server/"$3".test profiles/server/cpu_"$3".out
elif [ "$2" = "fmt" ]
then
    go install golang.org/x/tools/cmd/goimports@latest
    find . -name "*.go" | while read file; do
        echo "Formatting $file"
        gofmt -w "$file"
        goimports -local "github.com/Arcadian-Sky/musthave-metrics" -w "$file"
    done
fi



# curl http://127.0.0.1:8080/debug/pprof/heap?seconds=120 -o ./profiles/mem_out.pprof
# curl http://127.0.0.1:8080/debug/pprof/profile?seconds=30 -o ./profiles/cpu_out.pprof

# go tool pprof -svg -alloc_objects server ./profiles/mem_out.pprof > ./profiles/mem_ao.svg
# go tool pprof -svg server ./profiles/cpu_out.pprof > ./profiles/cpu.svg
# pprof -top ./profiles/cpu_out.pprof
# pprof -http=":9090" ./profiles/cpu_out.pprof
# pprof -top -diff_base=profiles/base_cpu_out.pprof profiles/cpu_out.pprof



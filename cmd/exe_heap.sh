rm -f main
rm -f heap.prof

go build -o main

./main

go tool pprof ./main heap.prof

rm -f main
rm -f heap

go build -o main

./main

go tool pprof ./main heap

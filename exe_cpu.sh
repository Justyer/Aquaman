rm -f main
rm -f cpu

go build -o main

./main

go tool pprof ./main cpu

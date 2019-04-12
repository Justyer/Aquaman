rm -f main
rm -f cpu.prof

go build -o main

./main

go tool pprof ./main cpu.prof

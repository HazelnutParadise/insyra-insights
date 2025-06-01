.PHONY: all clean

all: insyra-insights

insyra-insights: cmd/loader/core.exe
	go build -o insyra-insights ./cmd/loader

cmd/loader/core.exe:
	go build -o cmd/loader/core.exe ./cmd/core

clean:
	rm -f cmd/loader/core.exe insyra-insights

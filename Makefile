VERSION=v0.0.1

.PHONY: build deck win
build:
	go build -ldflags "-X main.version=$(VERSION) -X main.buildtime=`date +%Y-%m-%d@%H:%M:%S`" -o bin/

clean:
	rm -rf bin/*

deck:
	cd build && ./build_4_deck.sh $(VERSION) && cd ..

run: build
	./bin/synt-go-thic

win:
	env GOOS=windows GOARCH=amd64   go build -ldflags "-X main.version=$(VERSION) -X main.buildtime=`date +%Y-%m-%d@%H:%M:%S`" -o bin/steam_synth-go-thic

upload: deck
	scp bin/steam_synth-go-thic deck@steamdeck.local:~/Downloads/
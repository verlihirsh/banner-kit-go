BIN=banner-gen

.PHONY: build
build:
	go build -o $(BIN)

.PHONY: test
test: build
	./$(BIN) ~/work/banner-kit/test-projects/nerd-icon-example light center
	./$(BIN) ~/work/banner-kit/test-projects/emoji-test dark left
	./$(BIN) ~/work/banner-kit/test-projects/vibetunnel muted right

.PHONY: clean
clean:
	rm -f $(BIN)
	find ~/work/banner-kit/test-projects -name "banner.svg" -delete
	find ~/work/banner-kit/test-projects -name "banner.png" -delete

.PHONY: install
install: build
	cp $(BIN) ~/bin/ || cp $(BIN) /usr/local/bin/

.PHONY: run
run: build
	./$(BIN)

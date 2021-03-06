.PHONY: build
build:
	cd authservice && go build -o bin/auth -v ./cmd/
	cd imageservice && go build -o bin/image -v ./cmd/
	cd postservice && go build -o bin/post -v ./cmd/
	cd userservice && go build -o bin/user -v ./cmd/

.PHONY: run
run:
	cd authservice && ./bin/auth &> logauth.txt&
	cd imageservice && ./bin/image &> logimage.txt&
	cd postservice && ./bin/post &> logpost.txt&
	cd userservice && ./bin/user &> loguser.txt&

.PHONY: stop
stop:
	killall -9 auth
	killall -9 image
	killall -9 post
	killall -9 user

.PHONY: test
test:
	go test ./...

.PHONY: cover
cover:
	go test -v -coverpkg=./... -coverprofile=profile.cov ./...
	go tool cover -html=profile.cov -o cover.html

.PHONY: db
db:
	psql -U postgres -d postgres -a -f userservice/userschema.sql
	psql -U postgres -d postgres -a -f postservice/postschema.sql

.PHONY: tar
tar:
	killall -9 tarantool
	cd db && rm *.snap && rm *.xlog && tarantool fl_store.lua
.PHONY: test-cover-html
PACKAGES = $(shell find ./ -type d -not -path '*/\.*')

test-cover-html:
	echo "mode: count" > coverage-all.out
	$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)
	go tool cover -html=coverage-all.out

.DEFAULT_GOAL := build
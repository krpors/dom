all:
	go install

# Generate coverage outfile + html.
coverage:
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out -o=coverage.html

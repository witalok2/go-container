# Export enviroment variables to commands
export

# Variables
go_cover_file=coverage.out

test:: ## Do the tests in go
	@ go test -race -coverprofile $(go_cover_file) ./...

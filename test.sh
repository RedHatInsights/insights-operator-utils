go test -race -covermode=atomic -coverprofile coverage.out $(go list ./... | grep -v tests) $@
exit_code=$?

echo "Test coverage:"
go tool cover -func=coverage.out

exit $exit_code

PKG_LIST=$(go list ./... | grep -v /vendor/)
for package in ${PKG_LIST}; do
    go test -covermode=count -coverprofile "cover/${package##*/}.cov" "$package" ;
done
tail -q -n +2 cover/*.cov >> cover/coverage.cov
go tool cover -html=cover/coverage.cov -o coverage.html

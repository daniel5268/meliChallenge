echo "migrate -database $DB_CONNECTION -path ./migrations up"
migrate -database $DB_CONNECTION -path ./migrations up
go test ./...
if [ $? -ne 0 ]; then
  echo "Unit tests failed"
  exit 1
fi;

# cloud-final-services

## Project Hierarchy
All microservices are stored in separate directory. For go services, they all use the same `go.mod` file, 
but separate `main.go` files and `Dockerfile`s.

## Running Go services
1. Go to the service's directory (e.g. `/api`)
2. Rename `config/config.example.yml` to `config/config.yml` and fill in missing configurations.
3. Run  
```shell
   $ go run main.go
   ```

## Migration
Don't forget to copy `./migrations/migrate.sh` -> `./migrations/migrate.local.sh` and edit the 
credentials to match your database.

### Create migration
```shell
$  migrate create -ext sql -dir migrations -seq <migration_name>
```
### Apply migration (up)
Linux:
```shell
$  ./migrations/migrate.local.sh
```

Windows:
```shell
$  ./migrations/migrate.local.ps1
```
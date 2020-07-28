# Run

- Locally

      ```
      $ go run server.go
      ```

- Docker

      ```
      $ docker build --force-rm -t meli .
      ```

      ```
      $ docker run -it -p 8000:8000 --restart=always --name=meli meli
      ```
 
 
# Migrations
    Script:
        ```
            ./scripts/run_migrations.sh
        ```
# Environment Variables

    `export $(cat configs/.env | grep -v ^# | xargs)`

# Tests

    Generate mocks:
        mockgen -source=app/item/item_controller.go -destination=internal/mocks/mock_item_controller.go
        
    Generate a test suite
        ginkgo bootstrap
        
    Generate a test file
        ginkgo generate
    
    Run tests
        export $(cat config/.env.testing | grep -v ^# | xargs) && go test ./... -coverprofile=coverage.out
    
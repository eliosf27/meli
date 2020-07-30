# Run

- Locally

    $ go run server.go

- Docker

    $ docker build --env-file configs/.env --force-rm -t meli .
      
    $ docker run -it --env-file config/.env -p 8000:8000 --restart=always --name=meli meli
      

# Update dependencies vendor dir (Try to not use this command unless is strictly necessary)
    
    go mod vendor
     
 
# Migrations
    
    Script:

        ./scripts/run_migrations.sh

# Environment Variables

    export $(cat configs/.env | grep -v ^# | xargs)

# Tests

    Generate mocks:
        mockgen -source=app/item/item_controller.go -destination=internal/mocks/mock_item_controller.go
        
    Generate a test suite
        ginkgo bootstrap
        
    Generate a test file
        ginkgo generate
    
    Run tests
        export $(cat config/.env.testing | grep -v ^# | xargs) && go test ./... -coverprofile=coverage.out
    
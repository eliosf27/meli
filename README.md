# Run

- Locally

    $ go run server.go

- Docker

    $ docker build --env-file config/.env --force-rm -t meli .
      
    $ docker run -it --env-file config/.env -p 8000:8000 --restart=always --name=meli meli
      

# Update dependencies vendor dir (Try to not use this command unless is strictly necessary)
    
    go mod vendor
     
 
# Migrations
    Docker:
        
            Build:
            
                docker build -f Dockerfile.migrations --force-rm -t meli-migrations .    
                
            Run:

                docker run --network="host" --env-file config/.env -it meli-migrations

    Script:

        ./scripts/run_migrations.sh

# Environment Variables

    export $(cat config/.env | grep -v ^# | xargs)

# Tests

    Generate mocks:
    
        mockgen -source=app/item/item_controller.go -destination=internal/mocks/mock_item_controller.go
        
    Generate a test suite
    
        ginkgo bootstrap
        
    Generate a test file
    
        ginkgo generate
    
    Run tests
    
        export $(cat config/.env.testing | grep -v ^# | xargs) && go test ./... -coverprofile=coverage.out
    
# docker-compose

    Run project
    
        docker-compose -f docker-compose.db.yml -f docker-compose.yml up --build
         
    Run migrations
        
        docker-compose -f docker-compose.db.yml -f docker-compose.migrations.yml up
        
    Scripts
        - Enter in the scripts directory
            
            cd scripts
        
        - Run and build the project
        
            make build 
            
        - Run the project
        
            make run
            
        - Run the migrations
            
            make run
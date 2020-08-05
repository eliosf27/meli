# Items

This project is a test made it to a mercadolibre interview

This project was build keeping in mind the most important
things about [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html), are these:

  1. **Independent of Frameworks**:

  The architecture does not depend on the existence of some library of feature laden software. This allows
  you to use such frameworks as tools, rather than having to cram your system into their limited constraints.

  2. **Testable**:

  The business rules can be tested without the UI, Database, Web Server, or any other external element.

  3. **Independent of UI**:

  The UI can change easily, without changing the rest of the system. A Web UI could be replaced with a
  console UI, for example, without changing the business rules.

  4. **Independent of Database**:

  You can swap out Oracle or SQL Server, for Mongo, BigTable, CouchDB, or something else. Your business rules
  are not bound to the database.

# Configurations

### Run

- Locally

    $ go run server.go

- Docker

    $ docker build --env-file config/.env --force-rm -t meli .
      
    $ docker run -it --env-file config/.env -p 8000:8000 --restart=always --name=meli meli
      

### Update dependencies vendor dir (Try to not use this command unless is strictly necessary)
    
    go mod vendor
     
 
### Migrations
    Docker:
        
            Build:
            
                docker build -f Dockerfile.migrations --force-rm -t meli-migrations .    
                
            Run:

                docker run --network="host" --env-file config/.env -it meli-migrations


### Environment Variables

    export $(cat config/.env | grep -v ^# | xargs)

### Tests

    Generate mocks:
    
        mockgen -source=internal/app/item/item_controller.go -destination=internal/mocks/mock_item_controller.go
        
    Generate a test suite
    
        ginkgo bootstrap
        
    Generate a test file
    
        ginkgo generate
    
    Run tests
    
        export $(cat config/.env.testing | grep -v ^# | xargs) && go test ./... -coverprofile=coverage.out
    
### docker-compose

    Run project
    
        docker-compose -f docker-compose.db.yml -f docker-compose.yml up --build
         
    Run migrations
        
        docker-compose -f docker-compose.db.yml -f docker-compose.migrations.yml up
        
### Scripts
    - Enter in the scripts directory
            
        cd scripts
        
    - Run and build the project
        
        make build 
            
    - Run the project
        
        make run
            
    - Run the migrations
            
        make migrations
  
# Technologies

- [echo golang framework](https://echo.labstack.com/)
- [redis](https://redis.io/)
- [postgres](https://www.postgresql.org/docs/12/index.html)
- [docker](https://www.docker.com/)
- [docker-compose](https://docs.docker.com/compose/)
- [gomock](https://github.com/golang/mock)
- [gomega](https://onsi.github.io/gomega/)
- [test-containers](https://www.testcontainers.org/)
- [ginkgo](https://github.com/onsi/ginkgo)
- [sling-http-client](https://github.com/dghubble/sling)

# Documentation

- You can find the documentation and diagrams in the documentation directory in the root of the project

/documentation
        
# Improves/Comments

- I use queues and redis to improve the performance of the metrics calculation to avoid fetching data from postgres or any other 
storage like mongo. Only one metric is calculated at the time to keep the consistency of the metrics data.

- I use redis and/or postgres to keep an item data in a cache to not fetching data from the mercadolibre api every single time.

- I use toggle features to activate/deactivate a storage strategy for the cache. I create two endpoints to configure
and check the configuration of the storage used to cache the item data. The default storage is postgres.

## License
[MIT](https://choosealicense.com/licenses/mit/)
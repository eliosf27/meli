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
    Docker:
    
        Build:
            ```
                docker build -f Dockerfile.migrations --force-rm -t meli-migrations .
            ```      
            
        Run:
            ```
                docker run --network="host" --env-file configs/.env -it meli-migrations
            ```
    
    Script:
        ```
            ./scripts/run_migrations.sh
        ```
# Environment Variables

    `export $(cat configs/.env | grep -v ^# | xargs)`

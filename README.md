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

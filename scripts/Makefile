build:
	$(info Make: Building and running containers.)
	docker-compose -f ../docker-compose.db.yml -f ../docker-compose.yml up --build

run:
	$(info Make: Running containers.)
	docker-compose -f ../docker-compose.db.yml -f ../docker-compose.yml up

migrations:
	$(info Make: Running containers.)
	docker-compose -f ../docker-compose.db.yml -f ../docker-compose.migrations.yml up
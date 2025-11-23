NAME := vox
APP_IMAGE := vox-app
DATABASE_IMAGE := vox-database
DREG := dreg.buldev.com
DOMAIN := vox.buldev.com
LINUX_USER := gchalakov
SSH_KEY := ~/.ssh/ChaluSRV

.PHONY: clear app app-logs database database-logs redis redis-logs redis-connect all up down restart

clear:
	@clear

app:
	@echo "=== App ==="
	@echo "Compiling source code..."
	@cd app && make generate compile
	@echo "Building image..."
	@docker compose up -d --build app

push: app database
	@echo "=== Push ==="
	@docker tag $(IMAGE) $(DREG)/$(IMAGE):latest
	@docker push $(DREG)/$(APP_IMAGE):latest
	@docker push $(DREG)/$(DATABASE_IMAGE):latest
	@ssh -i $(SSH_KEY) $(LINUX_USER)@$(DOMAIN)

app-logs:
	@docker logs --follow $(NAME)-app-1

database:
	@echo "=== Database ==="
	@echo "Building image..."
	@docker compose up -d --build database

database-logs:
	@docker logs --follow $(NAME)-database-1

database-connect:
	@docker exec -it $(DATABASE_IMAGE)-1 psql \
		-h localhost \
		-U $$(grep 'POSTGRES_USER' .env | cut -d '=' -f2 | cut -c 2- | rev | cut -c 2- | rev) \
		-d $$(grep 'POSTGRES_DB' .env | cut -d '=' -f2 | cut -c 2- | rev | cut -c 2- | rev)

redis:
	@echo "=== Redis ==="
	@echo "Building image..."
	@docker compose up -d --build redis

redis-logs:
	@docker logs --follow $(NAME)-redis-1

redis-connect:
	@docker exec -it $(NAME)-redis-1 redis-cli

all: app database redis

up:
	docker compose up -d

down:
	docker compose down

restart: down up

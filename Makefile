NAME := vox
APP_IMAGE := vox-app
DATABASE_IMAGE := vox-database
DREG := dreg.buldev.com
DOMAIN := vox.buldev.com
LINUX_USER := gchalakov
SSH_KEY := ~/.ssh/ChaluSRV

.PHONY: clear app app-logs database database-logs all up down restart

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

all: app database

up:
		docker compose up -d

down:
		docker compose down

restart: down up

run:
	@docker compose --env-file ./app/.env up -d --build
stop:
	@docker compose down
watch:
	@air

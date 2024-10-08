SCALE = 2

run:
	@docker compose --env-file ./app/.env up -d --build --scale app=$(SCALE)
stop:
	@docker compose down
watch:
	@air

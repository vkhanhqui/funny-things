compose:
	docker-compose -f docker-compose.yml up -d --build && \
	docker-compose -f docker-compose.yml run -T wait-for-infra

compose-all:
	docker-compose -f docker-compose.yml up -d --build && \
	docker-compose -f docker-compose.yml run -T wait-for-infra && \
	docker-compose -f docker-compose.svc.yml up -d --build && \
	docker-compose -f docker-compose.svc.yml run -T wait-for-svc

compose-down:
	@docker-compose -f docker-compose.yml -f docker-compose.svc.yml down -t 10 --remove-orphans
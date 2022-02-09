update:
	docker-compose -f ../docker-compose.yml  down --rmi local
	docker-compose -f ../docker-compose.yml up  -d
run:
	docker-compose -f ../docker-compose.yml  down --rmi local
	docker-compose -f ../docker-compose.yml up  -d
	docker logs -f backend | grcat conf.dockerinfo
.PHONY: update run

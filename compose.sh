docker-compose down
go test ./...
docker build -t postgres_db ./tests/postgres/.
docker build -t karate ./tests/karate/.
docker build -t pubsub-config tests/pubsub/.
echo "Defining env variable for Grafana"
docker-compose build
docker-compose up -d --force-recreate

# new docker compose version
# docker compose down
# go test ./...
# docker build -t postgres_db ./tests/postgres/.
# docker build -t karate ./tests/karate/.
# docker build -t pubsub-config tests/pubsub/.
# echo "Defining env variable for Grafana"
# docker compose build
# docker compose up -d --force-recreate

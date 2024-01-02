#!/bin/bash

docker network create compliance_network
sh compose.sh

export PATH=$PATH:$PWD
chmod +x docker-compose
echo "Waiting for all systems go up..."
sleep 10


docker run --network compliance_network --rm --env COMPLIANCE_API_URL=http://temis-compliance-api:8193/compliance-int --env MOCKS_URL=http://mocks:9093 --env PUBSUB_PROJECT_HOST=pubsub:8681 --env METRICS_API_URL=http://temis-compliance-api:7777 --env METRICS_EVENTS_URL=http://compliance:7777 karate -t ~@ignore -T 4 /karate/features/apis
if [ $? -ne 0 ]; then
   echo "Error running integration tests"
   exit 1
fi

docker run --network compliance_network --rm --env COMPLIANCE_API_URL=http://temis-compliance-api:8193/compliance-int --env MOCKS_URL=http://mocks:9093 --env PUBSUB_PROJECT_HOST=pubsub:8681 --env METRICS_API_URL=http://temis-compliance-api:7777 --env METRICS_EVENTS_URL=http://compliance:7777 karate -t ~@ignore -T 4 /karate/features/engines
if [ $? -ne 0 ]; then
   echo "Error running integration tests"
   exit 1
fi

docker run --network compliance_network --rm --env COMPLIANCE_API_URL=http://temis-compliance-api:8193/compliance-int --env MOCKS_URL=http://mocks:9093 --env PUBSUB_PROJECT_HOST=pubsub:8681 --env METRICS_API_URL=http://temis-compliance-api:7777 --env METRICS_EVENTS_URL=http://compliance:7777 karate -t ~@ignore -T 4 /karate/features/events
if [ $? -ne 0 ]; then
   echo "Error running integration tests"
   exit 1
fi

docker run --network compliance_network --rm --env COMPLIANCE_API_URL=http://temis-compliance-api:8193/compliance-int --env MOCKS_URL=http://mocks:9093 --env PUBSUB_PROJECT_HOST=pubsub:8681 --env METRICS_API_URL=http://temis-compliance-api:7777 --env METRICS_EVENTS_URL=http://compliance:7777 karate -t ~@ignore -T 4 /karate/features/metrics
if [ $? -ne 0 ]; then
   echo "Error running integration tests"
   exit 1
fi

docker run --network compliance_network --rm --env COMPLIANCE_API_URL=http://temis-compliance-api:8193/compliance-int --env MOCKS_URL=http://mocks:9093 --env PUBSUB_PROJECT_HOST=pubsub:8681 --env METRICS_API_URL=http://temis-compliance-api:7777 --env METRICS_EVENTS_URL=http://compliance:7777 karate -t ~@ignore -T 4 /karate/features/rules
if [ $? -ne 0 ]; then
   echo "Error running integration tests features"
   exit 1
fi

docker-compose stop
docker-compose down -v

rm -f docker-compose

exit 0

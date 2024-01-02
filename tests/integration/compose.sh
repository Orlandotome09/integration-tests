docker-compose down

cp -r features ../karate/features
cp -r karate-config.js ../karate
docker build -t karate ../karate/.
rm -rf ../karate/features
rm -rf ../karate/karate-config.js

docker build -t pubsub-config ../pubsub/.

docker build -t postgres_db ../postgres --no-cache

# If it is running locally, define ENV variable
if [ -z "$SHORT_COMMIT" ]; then
  echo "Defining env variable SHORT_COMMIT"
  export SHORT_COMMIT=local
  docker build -t temis-compliance:${SHORT_COMMIT} ../../.
  if [ $? -ne 0 ]; then
      echo "Error building Temis Compliance"
      exit 1
  fi
fi

# If docker-compose is not installed, download it
docker-compose ps
if [ $? -ne 0 ]; then

  curl -L "https://github.com/docker/compose/releases/download/1.27.3/docker-compose-$(uname -s)-$(uname -m)" -o docker-compose
  if [ $? -ne 0 ]; then
      echo "Error getting docker compose"
      exit 1
  fi
  export PATH=$PATH:$PWD
  chmod +x docker-compose
fi

docker-compose up --force-recreate -d --build --remove-orphans
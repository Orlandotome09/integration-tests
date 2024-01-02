#!/bin/sh

# If pubsub host is not defined use default pubsub
if [ -z "$PUBSUB_HOST" ]
then
  echo "Defining default pubsub host"
  export PUBSUB_HOST=pubsub
fi


/karate/wait-for -t 30 ${PUBSUB_HOST}:8681 -- java -Xms2g -Xmx6g -XX:MaxPermSize=450m -Dlogback.configurationFile=log-config.xml  -cp karate.jar:./features:KarateUtils.jar com.intuit.karate.Main "$@"

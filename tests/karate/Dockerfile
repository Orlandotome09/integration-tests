FROM openjdk:alpine
WORKDIR /karate
RUN mkdir /features
ADD . /karate

RUN chmod +x /karate/run.sh

# Current directory and ./features folder are in class path. Karate config
ENTRYPOINT ["/karate/run.sh"]
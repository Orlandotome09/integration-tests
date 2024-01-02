FROM openjdk:alpine
WORKDIR /karate
RUN mkdir /cases
ADD . /karate

# Current directory and ./cases folder are in class path. Karate config
ENTRYPOINT ["java","-Dlogback.configurationFile=log-config.xml","-cp","karate.jar:./cases","com.intuit.karate.Main","$@"]
# Installing Kafka

## Installing the Confluent Platform (Kafka, Zookeeper, and all the Confluent tools)

I chose to install Kafka using the Confluent Platform, since that's what I'll be using at work.  I also chose to install everything as docker containers using docker-compose because I want to be awesome.

```bash
cd ~/[where you keep your code]

git clone ssh://git@github.com/confluentinc/examples.git
cd examples/cp-all-in-one

docker-compose up -d --build
```

I followed the instructions [here](https://docs.confluent.io/current/quickstart/ce-docker-quickstart.html) to verify that my environment was running, and also to learn about the Control Center.

## Installing the Kafka CLI tools

Kafka comes packaged with some really useful command-line tools, but I couldn't figure out how to access them in one of the docker containers, so I decided to download the CLI tools separately.

### Installing Java

In order to run these tools, even though my Kafka is in a docker container, I had to install the Java JDK

```bash
brew tap caskroom/cask
brew cask install java
```

### Installing the CLI Tools
```bash
cd /tmp
curl https://mirrors.gigenet.com/apache/kafka/2.4.1/kafka-2.4.1-src.tgz —output /tmp/kafka-2.4.1-src.tgz

tar -xzf kafka_2.12-2.4.1.tgz

mv kafka_2.12-2.4.1 ~/kafka
# The tools are located in ~/kafka/bin
```

### Verifying the Kafka environment using the command line tools

And finally I was able to run a few commands from the command line to verify that things where working

```bash
cd ~/kafka/bin

# Create a new 'test' topic
./kafka-topics.sh --create --zookeeper localhost:2181 --replication-factor 1 --partitions 1 --topic test

# Show the details of the 'test' topic
./kafka-topics.sh --zookeeper localhost:2181 --describe --topic test

# List all the topics in your Kafka
./kafka-topics.sh --list --zookeeper localhost:2181

# Produce messages to the 'test' topic
./kafka-console-producer.sh --broker-list localhost:9092 --topic test
Test Message 1
Test Message 2
^D
# Press ctrl+d to stop entering messages

# Consume messages from the 'test' topic
./kafka-console-consumer.sh --bootstrap-server localhost:9092 --topic test --from-beginning
# you should see your messages now .. press ctrl+c to stop consuming
```


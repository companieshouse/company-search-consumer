# company-search-consumer
A service that consumes messages from the stream-company-profile topic. The messages are then unmarhsalled and a RESTful request is sent to an external API with the intention of adding or updating documents present in an elastic search index.

Requirements
=============

- [Go](https://golang.org/doc/install)
- [Git](https://git-scm.com/downloads)
- [Kafka](https://kafka.apache.org/) - this is installed as part of the chs local build (within the chs-kafka VM)
- [Kafka 2](https://companieshouse.atlassian.net/wiki/spaces/DEV/pages/563806512/Setting+up+Kafka+2+on+Docker+for+Mac+for+use+on+chs-dev) - use this within Docker

Getting started
--------------
To build the service:
 1. Clone the repository
 2. Build the executable using:
 ```shell
 make build
 ```
 3. Run the service in foreground with the start script:
 ```shell
 ./start.sh
 ```
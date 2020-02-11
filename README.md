# company-search-consumer
A service that consumes messages from the stream-company-profile topic and then sends a RESTful POST request to the search.api.ch.gov.uk service. The data sent through the POST request will be used in search.api.ch.gov.uk to maintain and update the elastic search alpha_index.

Requirements (local setup)
=============

- [Go](https://golang.org/doc/install)
- [Git](https://git-scm.com/downloads)
- [Kafka](https://kafka.apache.org/) - this is installed as part of the chs local build (within the chs-kafka VM)
- [Kafka 2]

Getting started
--------------
To build the service:
 1. Clone the repository into your GOPATH under `src/github.com/companieshouse`
 2. Build the executable using:
 ```shell
 make build
 ```
 3. Run the service in foreground with the start script:
 ```shell
 ./start.sh
 ```
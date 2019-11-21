package main

import (
	"github.com/companieshouse/chs.go/log"
	"github.com/companieshouse/company-search-consumer/config"
)

func main() {
	cfg := config.Get()

	log.Namespace = cfg.Namespace()
}

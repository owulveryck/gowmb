// Copyright 2016 Olivier Wulveryck
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/kelseyhightower/envconfig"
)

type configuration struct {
	Debug       bool
	Scheme      string
	Port        int
	Address     string
	PrivateKey  string
	Certificate string
}

var config configuration

func main() {

	var help = flag.Bool("help", false, "show help message")
	flag.Parse()
	// Default values
	config.Port = 443
	config.Scheme = "https"
	config.Address = "0.0.0.0"
	config.Debug = false
	config.PrivateKey = "ssl/server.key"
	config.Certificate = "ssl/server.pem"
	defaultConf := config
	err := envconfig.Process("DEMO", &config)
	if err != nil {
		log.Fatal(err.Error())
	}
	log.Printf("==> DEMO_PORT: %v (default: %v)", config.Port, defaultConf.Port)
	log.Printf("==> DEMO_SCHEME: %v (default: %v)", config.Scheme, defaultConf.Scheme)
	log.Printf("==> DEMO_ADDRESS: %v (default: %v)", config.Address, defaultConf.Address)
	log.Printf("==> DEMO_DEBUG: %v (default: %v)", config.Debug, defaultConf.Debug)
	log.Printf("==> DEMO_PRIVATEKEY: %v (default: %v)", config.PrivateKey, defaultConf.PrivateKey)
	log.Printf("==> DEMO_CERTIFICATE: %v (default: %v)", config.Certificate, defaultConf.Certificate)
	if *help {
		os.Exit(0)
	}

	router := NewRouter()

	//	go server.AllHubs.Run()
	addr := fmt.Sprintf("%v:%v", config.Address, config.Port)
	if config.Scheme == "https" {
		log.Fatal(http.ListenAndServeTLS(addr, config.Certificate, config.PrivateKey, router))

	} else {
		log.Fatal(http.ListenAndServe(addr, router))

	}
}

package config

import "flag"

type flagValues struct {
	config   string
	verbose  bool
	grpcAddr string
	httpAddr string
}

func readFlags(values flagValues) flagValues {
	flag.StringVar(&values.config, "config", values.config, "path to config file")

	flag.StringVar(&values.grpcAddr, "grpc_addr", values.grpcAddr, "grpc server address")
	flag.StringVar(&values.httpAddr, "http_addr", values.httpAddr, "http server address")
	flag.BoolVar(&values.verbose, "verbose", values.verbose, "enable verbose logging")

	flag.Parse()

	return values
}

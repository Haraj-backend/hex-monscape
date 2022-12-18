package main

type config struct {
	LocalDeployment          localDeploymentConfig `yaml:"local_deployment" cfg:"local_deployment"`
	Dynamo                   dynamoConfig          `yaml:"ddb" cfg:"ddb"`
	OtelExporterOTLPEndPoint string                `yaml:"otel_exporter_otlp_endpoint" cfg:"otel_exporter_otlp_endpoint"`
}

type localDeploymentConfig struct {
	Enabled  bool   `cfg:"enabled"`
	Endpoint string `cfg:"endpoint"`
	Port     int    `cfg:"port" cfgDefault:"9186"`
}

type dynamoConfig struct {
	BattleTable  string `cfg:"battle_table" cfgDefault:"Battles"`
	GameTable    string `cfg:"game_table" cfgDefault:"Games"`
	PokemonTable string `cfg:"pokemon_table" cfgDefault:"Pokemons"`
}

package main

type config struct {
	Port    string        `cfg:"port" cfgDefault:"9186"`
	Storage storageConfig `cfg:"storage"`
}

type storageConfig struct {
	Type   storageType         `cfg:"type" cfgDefault:"memory"`
	Memory storageMemoryConfig `cfg:"memory"`
}

type storageMemoryConfig struct {
	MonsterDataPath string `cfg:"monster_data_path" cfgRequired:"true"`
}

type storageType string

const (
	storageTypeMemory   storageType = "memory"
	storageTypeMySQL    storageType = "mysql"
	storageTypeDynamoDB storageType = "dynamodb"
)

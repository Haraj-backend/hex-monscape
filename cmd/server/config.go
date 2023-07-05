package main

type config struct {
	Port    string        `cfg:"port" cfgDefault:"9186"`
	Storage storageConfig `cfg:"storage"`
}

type storageConfig struct {
	Type     storageType           `cfg:"type" cfgDefault:"memory"`
	Memory   storageMemoryConfig   `cfg:"memory"`
	DynamoDB storageDynamoDBConfig `cfg:"dynamodb"`
	MySQL    storageMySQLConfig    `cfg:"mysql"`
}

type storageMemoryConfig struct {
	MonsterDataPath string `cfg:"monster_data_path"`
}

type storageDynamoDBConfig struct {
	LocalstackEndpoint string `cfg:"localstack_endpoint"`
	BattleTableName    string `cfg:"battle_table_name"`
	GameTableName      string `cfg:"game_table_name"`
	MonsterTableName   string `cfg:"monster_table_name"`
}

type storageMySQLConfig struct {
	SQLDSN string `cfg:"sql_dsn"`
}

type storageType string

const (
	storageTypeMemory   storageType = "memory"
	storageTypeMySQL    storageType = "mysql"
	storageTypeDynamoDB storageType = "dynamodb"
)

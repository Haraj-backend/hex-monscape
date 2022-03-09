package shared

var TestConfig = struct {
	EnvKeyLocalstackEndpoint string
	EnvKeyBattleTableName    string
}{
	EnvKeyLocalstackEndpoint: "LOCALSTACK_ENDPOINT",
	EnvKeyBattleTableName:    "DDB_TABLE_BATTLE_NAME",
}

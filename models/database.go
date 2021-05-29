package models

// Configuration model
type Config struct {
	DatabaseEngine string `mapstructure:"database_engine"`
	DatabaseServer string `mapstructure:"database_server"`
	DatabasePort   string `mapstructure:"database_port"`
	DatabaseName   string `mapstructure:"database_name"`
	Proto          string `proto:"proto"`
	ServerPort     string `proto:"server_port"`
}

type Collection struct {
	CollectionName string  `mapstructure:"collection_name"`
	Fields         []Field `mapstructure:"fields"`
}

type Field struct {
	FieldName string `mapstructure:"field_name"`
	Unique    bool   `mapstructure:"unique"`
}

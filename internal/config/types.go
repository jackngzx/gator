package config

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

type State struct {
	Cfg *Config
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	RegisteredCommand map[string]func(*State, Command) error
}

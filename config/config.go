package config

import (
    "io/ioutil"
    "github.com/adrg/xdg"
    "github.com/pelletier/go-toml"
)

const CONFIG_PATH string = "SpaceFarmerBot/config.toml"

/* toml structs */
type Config struct {
    Bot BotConfig
}
type BotConfig struct {
    Token string
}

func GetConfig() (*Config, error) {
    /* get config path */
    config_path, err := xdg.SearchConfigFile(CONFIG_PATH)
    if err != nil {
        return nil, err
    }

    /* read config data */
    toml_data, err := ioutil.ReadFile(config_path)
    if err != nil {
        return nil, err
    }

    /* parse config data */
    var config Config
    if err := toml.Unmarshal(toml_data, &config); err != nil {
        return nil, err
    }
    return &config, nil
}

package db

type Config struct {
    Host    string  `yaml:"host"`
    Name    string  `yaml:"name"`
    User    string  `yaml:"user"`
    Password string `yaml:"password"`
    Port    int     `yaml:"port"`
}

func NewConfig() *Config {
    return &Config{

    }
}

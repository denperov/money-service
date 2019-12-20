package main

type Config struct {
	API struct {
		ListenAddress string `env:"API_LISTEN_ADDRESS" envDefault:"0.0.0.0:8080"`
	}
	Database struct {
		Address  string `env:"DATABASE_ADDRESS,required"`
		Name     string `env:"DATABASE_NAME,required"`
		User     string `env:"DATABASE_USER,required"`
		Password string `env:"DATABASE_PASSWORD,required"`
	}
}

package domain

type Env struct {
	PORT                   int    `envconfig:"PORT"                      required:"true"`
	APP_ENV                string `envconfig:"APP_ENV"                   required:"true"`
	DB_HOST                string `envconfig:"DB_HOST"                   required:"true"`
	DB_PORT                string `envconfig:"DB_PORT"                   required:"true"`
	DB_DATABASE            string `envconfig:"DB_DATABASE"               required:"true"`
	DB_USERNAME            string `envconfig:"DB_USERNAME"               required:"true"`
	DB_PASSWORD            string `envconfig:"DB_PASSWORD"               required:"true"`
	AccessTokenExpiryHour  int    `envconfig:"ACCESS_TOKEN_EXPIRY_HOUR"  required:"true"`
	RefreshTokenExpiryHour int    `envconfig:"REFRESH_TOKEN_EXPIRY_HOUR" required:"true"`
	AccessTokenSecret      string `envconfig:"ACCESS_TOKEN_SECRET"       required:"true"`
	RefreshTokenSecret     string `envconfig:"REFRESH_TOKEN_SECRET"      required:"true"`
}

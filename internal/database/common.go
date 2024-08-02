package database

const DriverName = "mysql"

type Parameters struct {
	Host     string
	Database string
	Username string
	Password string
	Port     int
}

type Config struct {
	Charset                    string
	ParseTime                  bool
	MultiStatements            bool
	EstablishConnectionTimeout string
}

func DefaultConnectionConfig() Config {
	return Config{
		Charset:                    "utf8",
		ParseTime:                  true,
		MultiStatements:            true,
		EstablishConnectionTimeout: "3s",
	}
}

package config

import "fmt"

type Database struct {
	Drive string `mapstructure:"drive" default:"mysql"`
	Host  string `mapstructure:"host" default:"localhost"`
	Port  int    `mapstructure:"port" default:"3306"`

	Charset   string `mapstructure:"charset" default:"utf8mb4"`
	ParseTime bool   `mapstructure:"parseTime" default:"true"`
	Loc       string `mapstructure:"loc" default:"Local"`

	Name     string `mapstructure:"name" default:""`
	User     string `mapstructure:"user" default:""`
	Password string `mapstructure:"password" default:""`

	Refresh bool `mapstructure:"refresh" default:"false"`
}

func (d *Database) Dsn() string {
	var dsn string
	switch d.Drive {
	case "mysql":
		dsn = fmt.Sprintf(
			"%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t&loc=%s",
			d.User, d.Password, d.Host, d.Port, d.Name, d.Charset, d.ParseTime, d.Loc,
		)
	}
	return dsn
}

package bot

type Config struct {
	groupID string
	token   string
	server  string
	key     string
	ts      string
	v       string
	debug   bool
}

type Option func(*Config)

func GroupID(gid string) Option {
	return func(c *Config) {
		c.groupID = gid
	}
}

func Token(token string) Option {
	return func(c *Config) {
		c.token = token
	}
}

func Version(v string) Option {
	return func(c *Config) {
		c.v = v
	}
}

func Debug(debug bool) Option {
	return func(c *Config) {
		c.debug = debug
	}
}

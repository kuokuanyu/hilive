package table

const (
	DEFAULT_DRIVER     = "mysql"
	DEFAULT_CONNECTION = "hilives"
	PRIMARY_NAME       = "id"
	PRIMARY_TYPE       = "INT"
)

// Config 設置資訊
type Config struct {
	driver     string
	canAdd     bool
	canEdit    bool
	canDelete  bool
	primaryKey PrimaryKey
	conn       string
}

// DefaultConfig 預設Config
func DefaultConfig() Config {
	return Config{
		driver:    DEFAULT_DRIVER,
		conn:      DEFAULT_CONNECTION,
		canAdd:    true,
		canEdit:   true,
		canDelete: true,
		primaryKey: PrimaryKey{
			Type: PRIMARY_TYPE,
			Name: PRIMARY_NAME,
		},
	}
}

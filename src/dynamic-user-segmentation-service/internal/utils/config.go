package utils

type DatabaseConfiguration struct {
	Type string `json:"type"`
	Name string `json:"name"`

	User     string `json:"user"`
	Password string `json:"password"`

	Host string `json:"host"`
	Port string `json:"port"`
}

var DbConfig DatabaseConfiguration

func InitConfiguration() {
	DbConfig = DatabaseConfiguration{
		"postgres",
		"dynamic-user-segmentation",
		"postgres",
		"postgres",
		"postgres",
		"5432",
	}
}

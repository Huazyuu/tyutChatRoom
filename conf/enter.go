package conf

type Config struct {
	Mysql      Mysql  `yaml:"mysql"`
	Log        Log    `yaml:"log"`
	System     System `yaml:"system"`
	Email      Email  `yaml:"email"`
	Jwt        Jwt    `yaml:"jwt"`
	Redis      Redis  `yaml:"redis"`
	UploadPath string `yaml:"upload_path"`
}

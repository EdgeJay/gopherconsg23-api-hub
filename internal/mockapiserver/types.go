package mockapiserver

type AppFlags struct {
	InputFile string
	Port      int
}

type MockDataMapping map[string]string

type MockDataExamplesMapping struct {
	Mappings map[string]MockDataExampleConfig `yaml:"mappings"`
}

type MockDataExampleConfig struct {
	Parameters []MockDataExampleConfigParameter `yaml:"parameters"`
}

type MockDataExampleConfigParameter struct {
	In    string `yaml:"in"`
	Name  string `yaml:"name"`
	Match string `yaml:"match"`
}

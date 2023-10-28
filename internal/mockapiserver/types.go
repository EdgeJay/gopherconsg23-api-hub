package mockapiserver

type AppFlags struct {
	InputFile string
	Port      int
}

type MockDataMapping map[string]string

package config

type testService struct {
	config *Config
}

func NewTestService(config *Config) Service {
	return &testService{
		config: config,
	}
}

func (s *testService) Get() *Config {
	return &(*s.config)
}

func (s *testService) Refresh() error {
	return nil
}

func (s *testService) Store(newStored *StoredConfig) {
}

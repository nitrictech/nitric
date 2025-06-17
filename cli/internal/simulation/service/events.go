package service

type ServiceLogWriter struct {
	output  OutputType
	service *ServiceSimulation
}

func (s *ServiceLogWriter) Write(content []byte) (int, error) {
	s.service.events <- ServiceEvent{
		SimulatedService: s.service,
		Output:           &s.output,
		Content:          &content,
	}

	return len(content), nil
}

func newServiceLogWriter(service *ServiceSimulation, output OutputType) *ServiceLogWriter {
	return &ServiceLogWriter{
		output:  output,
		service: service,
	}
}

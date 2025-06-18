package service

import "slices"

type ServiceLogWriter struct {
	output  OutputType
	service *ServiceSimulation
}

func (s *ServiceLogWriter) Write(content []byte) (int, error) {
	// Needed to prevent overwriting the underlying array
	data := slices.Clone(content)

	s.service.events <- ServiceEvent{
		SimulatedService: s.service,
		Output:           &s.output,
		Content:          data,
		PreviousStatus:   s.service.currentStatus,
	}

	return len(content), nil
}

func newServiceLogWriter(service *ServiceSimulation, output OutputType) *ServiceLogWriter {
	return &ServiceLogWriter{
		output:  output,
		service: service,
	}
}

package balancer

import (
	"balancer/internal/balancer/algoritms"
	"balancer/internal/domain/models"
	"balancer/internal/logger"
)

type BalanceAlgorithm interface {
	Next() *models.Backend
}

func GetBalancer(name string, backends []*models.Backend) BalanceAlgorithm {
	switch name {
	case "round-robin":
		logger.Base().Info("using round-robin balancing algorithm")
		return algoritms.NewRoundRobin(backends)
	case "least-connections":
		logger.Base().Info("using least-connections balancing algorithm")
		return algoritms.NewLeastConnections(backends)
	case "random":
		logger.Base().Info("using random balancing algorithm")
		return algoritms.NewRandom(backends)
	default:
		logger.Base().Warn("unknown algorithm name, defaulting to round-robin", "name", name)
		return algoritms.NewRoundRobin(backends)
	}
}

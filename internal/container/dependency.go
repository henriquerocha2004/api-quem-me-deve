package container

import (
	"github.com/henriquerocha2004/quem-me-deve-api/core/client"
	"github.com/henriquerocha2004/quem-me-deve-api/core/debt"
)

type Dependencies struct {
	DebtService   debt.Service
	ClientService client.Service
}

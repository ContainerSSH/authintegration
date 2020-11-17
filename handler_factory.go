package authintegration

import (
	"fmt"

	"github.com/containerssh/auth"
	"github.com/containerssh/sshserver"
)

// New creates a new handler that authenticates the users with passwords and public keys.
//goland:noinspection GoUnusedExportedFunction
func New(
	authClient auth.Client,
	backend sshserver.Handler,
	behavior Behavior,
) (sshserver.Handler, error) {
	if authClient == nil {
		return nil, fmt.Errorf("the authClient parameter to authintegration.New cannot be nil")
	}
	if backend == nil {
		return nil, fmt.Errorf("the backend parameter to authintegration.New cannot be nil")
	}
	if !behavior.validate() {
		return nil, fmt.Errorf("the behavior field contains an invalid value: %d", behavior)
	}
	return &handler{
		authClient: authClient,
		backend:    backend,
		behavior:   behavior,
	}, nil
}

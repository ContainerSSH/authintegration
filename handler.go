package authintegration

import (
	"context"
	"net"

	"github.com/containerssh/auth"
	"github.com/containerssh/log"
	"github.com/containerssh/sshserver"
)

// Behavior dictactes how when the authentication requests are passed to the backends.
type Behavior int

const (
	// BehaviorNoPassthrough means that the authentication integration will never call the backend for authentication.
	BehaviorNoPassthrough Behavior = iota
	// BehaviorPassthroughOnFailure will call the backend if the authentication server returned a failure.
	BehaviorPassthroughOnFailure Behavior = iota
	// BehaviorPassthroughOnSuccess will call the backend if the authentication server returned a success.
	BehaviorPassthroughOnSuccess Behavior = iota
	// BehaviorPassthroughOnUnavailable will call the backend if the authentication server is not available.
	BehaviorPassthroughOnUnavailable Behavior = iota
)

func (behavior Behavior) validate() bool {
	switch behavior {
	case BehaviorNoPassthrough:
		return true
	case BehaviorPassthroughOnFailure:
		return true
	case BehaviorPassthroughOnSuccess:
		return true
	case BehaviorPassthroughOnUnavailable:
		return true
	default:
		return false
	}
}

type handler struct {
	backend    sshserver.Handler
	authClient auth.Client
	behavior   Behavior
}

func (h *handler) OnReady() error {
	if h.backend != nil {
		return h.backend.OnReady()
	}
	return nil
}

func (h *handler) OnShutdown(shutdownContext context.Context) {
	if h.backend != nil {
		h.backend.OnShutdown(shutdownContext)
	}
}

func (h *handler) OnNetworkConnection(client net.TCPAddr, connectionID string) (sshserver.NetworkConnectionHandler, error) {
	var backend sshserver.NetworkConnectionHandler = nil
	var err error
	if h.backend != nil {
		backend, err = h.backend.OnNetworkConnection(client, connectionID)
		if err != nil {
			return nil, err
		}
	}
	return &networkConnectionHandler{
		connectionID: connectionID,
		ip:           client.IP,
		authClient:   h.authClient,
		backend:      backend,
		behavior:     h.behavior,
	}, nil
}

type networkConnectionHandler struct {
	backend      sshserver.NetworkConnectionHandler
	authClient   auth.Client
	ip           net.IP
	connectionID string
	behavior     Behavior
}

func (h *networkConnectionHandler) OnShutdown(shutdownContext context.Context) {
	h.backend.OnShutdown(shutdownContext)
}

func (h *networkConnectionHandler) OnAuthPassword(username string, password []byte) (response sshserver.AuthResponse, reason error) {
	success, err := h.authClient.Password(username, password, h.connectionID, h.ip)
	if !success {
		if err != nil {
			if h.behavior == BehaviorPassthroughOnUnavailable {
				return h.backend.OnAuthPassword(username, password)
			}
			return sshserver.AuthResponseUnavailable, err
		}
		if h.behavior == BehaviorPassthroughOnFailure {
			return h.backend.OnAuthPassword(username, password)
		}
		return sshserver.AuthResponseFailure, nil
	}
	if h.behavior == BehaviorPassthroughOnSuccess {
		return h.backend.OnAuthPassword(username, password)
	} else {
		return sshserver.AuthResponseSuccess, nil
	}
}

func (h *networkConnectionHandler) OnAuthPubKey(username string, pubKey string) (response sshserver.AuthResponse, reason error) {
	success, err := h.authClient.PubKey(username, pubKey, h.connectionID, h.ip)
	if !success {
		if err != nil {
			if h.behavior == BehaviorPassthroughOnUnavailable {
				return h.backend.OnAuthPubKey(username, pubKey)
			}
			return sshserver.AuthResponseUnavailable, err
		}
		if h.behavior == BehaviorPassthroughOnFailure {
			return h.backend.OnAuthPubKey(username, pubKey)
		}
		return sshserver.AuthResponseFailure, nil
	}
	if h.behavior == BehaviorPassthroughOnSuccess {
		return h.backend.OnAuthPubKey(username, pubKey)
	}
	return sshserver.AuthResponseSuccess, nil
}

func (h *networkConnectionHandler) OnAuthKeyboardInteractive(
	username string,
	challenge func(
		instruction string,
		questions sshserver.KeyboardInteractiveQuestions,
	) (answers sshserver.KeyboardInteractiveAnswers, err error),
) (response sshserver.AuthResponse, reason error) {
	if h.behavior == BehaviorPassthroughOnUnavailable {
		return h.backend.OnAuthKeyboardInteractive(username, challenge)
	}
	return sshserver.AuthResponseUnavailable, log.UserMessage(
		EUnsupported,
		"keyboard-interactive authentication not available",
		"Keyboard-interactive authentication is not supported by ContainerSSH.",
	).Label("connectionId", h.connectionID).Label("username", username)
}

func (h *networkConnectionHandler) OnHandshakeFailed(reason error) {
	h.backend.OnHandshakeFailed(reason)
}

func (h *networkConnectionHandler) OnHandshakeSuccess(username string) (connection sshserver.SSHConnectionHandler, failureReason error) {
	return h.backend.OnHandshakeSuccess(username)
}

func (h *networkConnectionHandler) OnDisconnect() {
	h.backend.OnDisconnect()
}

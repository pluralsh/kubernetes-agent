package api

import (
	"crypto/sha256"

	"go.opentelemetry.io/otel/attribute"
)

const (
	// TraceAgentIdAttr is tracing attribute that holds an agent id.
	TraceAgentIdAttr attribute.Key = "agent_id"
)

// AgentToken is agentk's bearer access token.
type AgentToken string

// AgentInfo contains information about an agentk.
type AgentInfo struct {
	// Id is the agent's id in the database.
	Id int64
	// ProjectId is the id of the configuration project of the agent.
	ProjectId int64

	// Name is the agent's name.
	// Can contain only /a-z\d-/
	Name string
	// DefaultBranch is the name of the default branch in the agent's configuration repository.
	DefaultBranch string
}

type ProjectInfo struct {
	ProjectId int64
	// DefaultBranch is the name of the default branch in a repository.
	DefaultBranch string
}

func AgentToken2key(token AgentToken) []byte {
	// We use only the first half of the token as a key. Under the assumption of
	// a randomly generated token of length at least 50, with an alphabet of at least
	//
	// - upper-case characters (26)
	// - lower-case characters (26),
	// - numbers (10),
	//
	// (see https://gitlab.com/gitlab-org/gitlab/blob/master/app/models/clusters/agent_token.rb)
	//
	// we have at least 62^25 different possible token prefixes. Since the token is
	// randomly generated, to obtain the token from this hash, one would have to
	// also guess the second half, and validate it by attempting to log in (kas
	// cannot validate tokens on its own)
	n := len(token) / 2
	tokenHash := sha256.Sum256([]byte(token[:n]))
	return tokenHash[:]
}

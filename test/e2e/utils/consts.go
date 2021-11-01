package utils

const (
	// RoleWorker contains the worker role
	RoleWorker = "worker"
	// DefaultNodeName we rely on kind for our CI
	DefaultNodeName = "kind-worker"
	// LabelRole contains the key for the role label
	LabelRole = "node-role.kubernetes.io"
	// LabelHostname contains the key for the hostname label
	LabelHostname = "kubernetes.io/hostname"
	// TestNodeLabel unique label for the node that
	TestNodeLabel = "rte-e2e-test-node"
)

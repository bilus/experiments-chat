package test

import "testing"

// TestFakeAgent is not really a test. It is an entrypoint for starting an
// ACP server with the fake agent for testing, injected into an ACP connection
// using UseFakeAgent.
func TestFakeAgent(t *testing.T) {
	RunFakeAgent(t)
}

// Package integration provides integration tests for StreamSpace.
// These tests verify component interaction across the API, database,
// and Kubernetes controller.
package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	streamv1alpha1 "github.com/streamspace/streamspace/api/v1alpha1"
)

var (
	testEnv   *envtest.Environment
	k8sClient client.Client
	cfg       *rest.Config
)

// TestMain sets up the test environment for integration tests.
func TestMain(m *testing.M) {
	testEnv = &envtest.Environment{
		CRDDirectoryPaths: []string{
			"../../k8s-controller/config/crd/bases",
		},
		ErrorIfCRDPathMissing: true,
	}

	var err error
	cfg, err = testEnv.Start()
	if err != nil {
		panic(err)
	}

	err = streamv1alpha1.AddToScheme(scheme.Scheme)
	if err != nil {
		panic(err)
	}

	k8sClient, err = client.New(cfg, client.Options{Scheme: scheme.Scheme})
	if err != nil {
		panic(err)
	}

	code := m.Run()

	err = testEnv.Stop()
	if err != nil {
		panic(err)
	}

	os.Exit(code)
}

// Helper functions for integration tests

// waitForCondition waits for a condition to be true with timeout.
func waitForCondition(timeout time.Duration, interval time.Duration, condition func() bool) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		if condition() {
			return true
		}
		time.Sleep(interval)
	}
	return false
}

// getTestContext returns a context with timeout for tests.
func getTestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 60*time.Second)
}

// createTestNamespace creates a namespace for testing.
func createTestNamespace(t *testing.T, name string) {
	t.Helper()
	// Implementation will be added when tests are assigned
}

// cleanupTestNamespace removes a test namespace and all its resources.
func cleanupTestNamespace(t *testing.T, name string) {
	t.Helper()
	// Implementation will be added when tests are assigned
}

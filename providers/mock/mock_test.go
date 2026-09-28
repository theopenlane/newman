package mock

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/theopenlane/newman"
)

// TestEmailSenderImplementation checks if EmailSender implements the newman EmailSender interface
func TestEmailSenderImplementation(t *testing.T) {
	var _ newman.EmailSender = (*EmailSender)(nil)
}

func TestVerify(t *testing.T) {
	sender, err := New("")
	require.NoError(t, err)

	require.NoError(t, sender.Verify(context.Background()))
}

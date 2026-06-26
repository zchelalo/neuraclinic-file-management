package grpc

import (
	"context"
	"testing"

	"github.com/zchelalo/neuraclinic-file-management/internal/modules/files/application"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/appctx"
	"github.com/zchelalo/neuraclinic-file-management/internal/shared/i18n"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestMapErrorLocalizesMessage(t *testing.T) {
	t.Parallel()

	ctx := appctx.WithLanguage(context.Background(), i18n.Spanish)
	err := mapError(ctx, application.ErrInvalidInput)
	st, ok := status.FromError(err)
	if !ok {
		t.Fatalf("expected grpc status, got %v", err)
	}
	if st.Code() != codes.InvalidArgument {
		t.Fatalf("status code = %s, want %s", st.Code(), codes.InvalidArgument)
	}
	if st.Message() != "entrada invalida" {
		t.Fatalf("status message = %q", st.Message())
	}
}

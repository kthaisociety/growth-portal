package database_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/kthaisociety/growth-portal/internal/database"
	"github.com/kthaisociety/growth-portal/internal/database/databasetest"
)

func TestMigrationsApplyToAnEmptyDatabaseAndRerunAsANoOp(t *testing.T) {
	db := databasetest.Open(t)

	require.NoError(t, database.Migrate(context.Background(), db))
}

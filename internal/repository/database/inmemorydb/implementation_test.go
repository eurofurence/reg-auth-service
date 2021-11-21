package inmemorydb

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/eurofurence/reg-auth-service/docs"
	"github.com/eurofurence/reg-auth-service/internal/entity"
	"github.com/stretchr/testify/require"
)

var (
	cut *InMemoryRepository
)

func TestMain(m *testing.M) {
	cut = &InMemoryRepository{}
	code := m.Run()
	os.Exit(code)
}

func tstSetup() {
	cut.Open()
}

func tstShutdown() {
	cut.Close()
}

func TestOpenClose(t *testing.T) {
	docs.Description("low level test for Open() and Close()")
	cut2 := &InMemoryRepository{}
	cut2.Open()
	require.NotNil(t, cut2.authRequests)
	cut2.Close()
	// Since we are not an actual database, closing the connection will only clear the repository.
	require.NotNil(t, cut2.authRequests)
}

func TestAddAuthRequest(t *testing.T) {
	docs.Description("it should be possible to add an auth request and then retrieve it again")
	tstSetup()
	defer tstShutdown()
	state := "test-state"
	ar := &entity.AuthRequest{State: state, ExpiresAt: time.Now().Add(time.Hour)}
	err := cut.AddAuthRequest(context.TODO(), ar)
	require.Nil(t, err, "unexpected error during add")

	ar2, err := cut.GetAuthRequestByState(context.TODO(), state)
	require.Nil(t, err, "unexpected error during get")
	require.EqualValues(t, *ar, *ar2, "comparison failure")
}

func TestAddAuthRequestStateAlreadyPresent(t *testing.T) {
	docs.Description("adding an auth request with a state already present should fail")
	tstSetup()
	defer tstShutdown()
	state := "test-state"
	ar := &entity.AuthRequest{State: state}
	err := cut.AddAuthRequest(context.TODO(), ar)
	require.Nil(t, err, "unexpected error during add")

	err2 := cut.AddAuthRequest(context.TODO(), ar)
	require.Equal(t, fmt.Sprintf("cannot add auth request '%s' - already present", state), err2.Error(), "unexpected error message")
}

func TestGetAuthRequestByStateNotFound(t *testing.T) {
	docs.Description("retrieving a nonexistent auth request should fail")
	tstSetup()
	defer tstShutdown()
	state := "inexistent-state"
	ar, err := cut.GetAuthRequestByState(context.TODO(), state)
	require.NotNil(t, err, "no error occurred, although it should have")
	require.Equal(t, fmt.Sprintf("cannot get auth request '%s' - not present", state), err.Error(), "unexpected error message")
	require.Nil(t, ar, "result entity should be nil")
}

func TestDeleteAuthRequestByState(t *testing.T) {
	docs.Description("deleting an existing auth request should succeed and it should be gone afterwards")
	tstSetup()
	defer tstShutdown()
	state := "test-state"
	ar := &entity.AuthRequest{State: state}
	err := cut.AddAuthRequest(context.TODO(), ar)
	require.Nil(t, err, "unexpected error during add")

	err2 := cut.DeleteAuthRequestByState(context.TODO(), state)
	require.Nil(t, err2, "unexpected error during delete")

	ar2, err3 := cut.GetAuthRequestByState(context.TODO(), state)
	require.NotNil(t, err3, "no error occurred, although it should have")
	require.Equal(t, fmt.Sprintf("cannot get auth request '%s' - not present", state), err3.Error(), "unexpected error message")
	require.Nil(t, ar2, "result entity should be nil")
}

func TestDeleteAuthRequestByStateNotFound(t *testing.T) {
	docs.Description("deleting a nonexistent auth request should fail")
	tstSetup()
	defer tstShutdown()
	state := "inexistent-state"
	err := cut.DeleteAuthRequestByState(context.TODO(), state)
	require.NotNil(t, err, "no error occurred, although it should have")
	require.Equal(t, fmt.Sprintf("cannot delete auth request '%s' - not present", state), err.Error(), "unexpected error message")
}

func TestPruneAuthRequestsEmpty(t *testing.T) {
	docs.Description("it should be possible to prune auth requests even if none are available")
	tstSetup()
	defer tstShutdown()

	pruneCount, err := cut.PruneAuthRequests(context.TODO())
	require.Nil(t, err, "unexpected error during prune")
	require.Equal(t, uint(0), pruneCount, "unexpected number of pruned entities")
}

func TestPruneAuthRequestsNoExpired(t *testing.T) {
	docs.Description("it should be possible to prune auth requests without removing any when none have expired")
	tstSetup()
	defer tstShutdown()
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: "test-state-1", ExpiresAt: time.Now().Add(time.Hour)})
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: "test-state-2", ExpiresAt: time.Now().Add(time.Hour)})
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: "test-state-3", ExpiresAt: time.Now().Add(time.Hour)})

	pruneCount, err := cut.PruneAuthRequests(context.TODO())
	require.Nil(t, err, "unexpected error during prune")
	require.Equal(t, uint(0), pruneCount, "unexpected number of pruned entities")
}

func TestPruneAuthRequestsSingleExpired(t *testing.T) {
	docs.Description("it should be possible to prune auth requests when only one has expired and verify it has been removed")
	tstSetup()
	defer tstShutdown()
	expiredState2 := "test-state-2-expired"
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: "test-state-1", ExpiresAt: time.Now().Add(time.Hour)})
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: expiredState2, ExpiresAt: time.Now().Add(-time.Hour)})
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: "test-state-3", ExpiresAt: time.Now().Add(time.Hour)})

	pruneCount, err := cut.PruneAuthRequests(context.TODO())
	require.Nil(t, err, "unexpected error during prune")
	require.Equal(t, uint(1), pruneCount, "unexpected number of pruned entities")

	ar, err := cut.GetAuthRequestByState(context.TODO(), expiredState2)
	require.NotNil(t, err, "no error occurred, although it should have")
	require.Equal(t, fmt.Sprintf("cannot get auth request '%s' - not present", expiredState2), err.Error(), "unexpected error message")
	require.Nil(t, ar, "result entity should be nil")
}

func TestPruneAuthRequestsMultipleExpired(t *testing.T) {
	docs.Description("it should be possible to prune auth requests when multiple have expired and verify they have been removed")
	tstSetup()
	defer tstShutdown()
	expiredState1 := "test-state-1-expired"
	expiredState3 := "test-state-3-expired"
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: expiredState1, ExpiresAt: time.Now().Add(-time.Hour)})
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: "test-state-2", ExpiresAt: time.Now().Add(time.Hour)})
	cut.AddAuthRequest(context.TODO(), &entity.AuthRequest{State: expiredState3, ExpiresAt: time.Now().Add(-time.Hour)})

	pruneCount, err := cut.PruneAuthRequests(context.TODO())
	require.Nil(t, err, "unexpected error during prune")
	require.Equal(t, uint(2), pruneCount, "unexpected number of pruned entities")

	ar1, err := cut.GetAuthRequestByState(context.TODO(), expiredState1)
	require.NotNil(t, err, "no error occurred, although it should have")
	require.Equal(t, fmt.Sprintf("cannot get auth request '%s' - not present", expiredState1), err.Error(), "unexpected error message")
	require.Nil(t, ar1, "result entity should be nil")

	ar2, err := cut.GetAuthRequestByState(context.TODO(), expiredState3)
	require.NotNil(t, err, "no error occurred, although it should have")
	require.Equal(t, fmt.Sprintf("cannot get auth request '%s' - not present", expiredState3), err.Error(), "unexpected error message")
	require.Nil(t, ar2, "result entity should be nil")
}

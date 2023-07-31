package redistool

import (
	"context"
	"fmt"

	"github.com/redis/rueidis"
)

// maxTransactionIterations is a safety mechanism to avoid infinite attempts.
const maxTransactionIterations = 100

// Optimistic locking pattern.
// See https://redis.io/docs/interact/transactions/
// See https://github.com/redis/rueidis#cas-pattern
func transaction(ctx context.Context, c rueidis.DedicatedClient, cb func(context.Context) ([]rueidis.Completed, error), keys ...string) (retErr error) {
	execCalled := false
	defer func() {
		if execCalled {
			return
		}
		// x. UNWATCH if there was an error or nothing to delete.
		err := c.Do(ctx, c.B().Unwatch().Build()).Error()
		if retErr == nil {
			retErr = err
		}
	}()
	for i := 0; i < maxTransactionIterations; i++ {
		// 1. WATCH
		execCalled = false // Enable deferred cleanup (for retries)
		err := c.Do(ctx, c.B().Watch().Key(keys...).Build()).Error()
		if err != nil {
			return err
		}
		// 2. READ
		cmds, err := cb(ctx)
		if err != nil {
			return err
		}
		if len(cmds) == 0 {
			return nil
		}
		// 3. Mutation via MULTI+EXEC
		multiExec := make([]rueidis.Completed, 0, len(cmds)+2)
		multiExec = append(multiExec, c.B().Multi().Build())
		multiExec = append(multiExec, cmds...)
		multiExec = append(multiExec, c.B().Exec().Build())
		resp := c.DoMulti(ctx, multiExec...)
		execCalled = true                         // Disable deferred UNWATCH as Redis UNWATCHes all keys on EXEC.
		err = MultiFirstError(resp[:len(resp)-1]) // all but the last one, which is EXEC
		if err != nil {                           // Something is wrong with commands or I/O, abort
			return err
		}
		// EXEC error
		switch resp[len(resp)-1].Error() { // nolint: errorlint
		case nil: // Success!
			return nil
		case rueidis.Nil: // EXEC detected a conflict, retry.
		default: // EXEC failed in a bad way, abort
			return err
		}
	}
	return fmt.Errorf("failed to execute Redis transaction %d times", maxTransactionIterations)
}

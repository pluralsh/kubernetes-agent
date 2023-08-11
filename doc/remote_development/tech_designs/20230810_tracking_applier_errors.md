# Technical Design for tracking applier errors
---

# Overview

- There may be errors at the time of applying k8s resources. At present these are captured within a goroutine but are not reported in Rails. So in case of errors, the user will be misled in the UI as the workspace will be stuck in the last possible stage.

# Special considerations

- Implementation of the tracking logic interleaves well with any combination of full/partial sync cycle
- Since applier errors are received asynchronously, the implementation should be able to accommodate edge cases pertaining to concurrency issues

# Solution Overview

The following components are affected as a part of this change:

- agentk
- rails

## Changes in agentk

An error tracker will be introduced in the reconciler to keep track of errors encountered as a part of the reconciliation loop. The tracker will be equipped to handle errors available in synchronous as well as asynchronous calls

### Error Tracker schema

```go
type errorTracker struct {
	mx    sync.Mutex
	store map[errorTrackerKey]operationState
}

type errorTrackerKey struct {
	name string
	namespace string
}

type operationState struct {
	version uint64
	err error // err contains the error available at the end of an async operation
}

type ErrorType string
const (
    ErrorTypeApplier ErrorType = "applier"
)
```

The error tracker will be created and bound to the lifecycle of the reconciler. Therefore, for every full sync, both the existing reconciler and the `errorTracker` shall be dropped and created from scratch for subsequent loops.

The `errorTracker` works similar to the `terminationTracker` in that it may require multiple reconciliation cycles to fully report the error details being tracked to rails as well as clean up the state.

The working of the errorTracker can be understood through its three distinct functions

- capturing an error
- reporting an error
- cleanup after successful reporting

**Capturing an error**

This takes place within the `applyWorkspaceChanges` method but it may occur asynchronously as well. This is especially relevant to calls the `k8sClient.applier.Run` function where an error channel is returned. In such cases, the tracker may be updated asynchronously.

Some questions to consider:

- can there be cases where multiple errors are available in the `chan error` in a single call to `k8sClient.applier.Run`?
  - For the sake of simplicity, the proposed changes will capture all errors per action published under the channel. Once the channel is closed, these errors will be merged and reported to rails in one call. The alternative implementation to capture/publish errors as they arrive will be slightly trickier but can be considered if absolutely necessary.

**Reporting an error**

At the start of each reconciliation loop, a snapshot of `errorTracker` will be created to collect the latest errors per workspace tracked _so far_. This is done in a manner similar to how a snapshot of workspaces is prepared using the `k8sInformer`. Another benefit of using a snapshot is that it renders handling of some edge cases unnecessary, thereby making the logic simpler. These cases primarily concern scenarios where the `errorTracker` may be updated in the background.

The main work of the reporting happens inside `generateWorkspaceAgentInfos`

The updated logic will behave as follows:

- For each existing workspace, check if either an unpersisted state exists. If yes, collect it for reporting to rails
- Iterate the `terminatingTracker` to collect information on workspaces with a termination progress to report. The data within `terminationTracker` should be independent of the data within the error snapshot
- Iterate through existing entries in the snapshot of errors and collect/prepare payload to report to Rails

Note: there may be entries in the error snapshot which may also have unpersisted data to report to Rails. In such a case, both the latest data as well applier errors must be sent to rails in the same request

**Cleanup after reporting**

The errorTracker _may_ cleanup an entry of a workspace within the `applyWorkspaceChanges` in a manner similar to the `terminationTracker`. If the response from Rails is successful, an entry in the `errorTracker` may be evicted IFF the current version of the entry in the `errorTracker` is the same as the version in the snapshot. This safety check ensures that an entry that tracks errors for a newer action is unaffected.

**IMPORTANT: Versioning within errorTracker**

One notable aspect of the errorTracker is that it only tracks the latest **version** of errors received per workspace/namespace. In order to understand why this is necessary, consider a hypothetical implementation errorTracker that doesn't have any concept of versioning i.e. it is essentially a `map[(workspace, namespace)]err`. Let's evaluate how this implementation behaves when subject to the following sequence of actions:

- partial sync 1: receive config 1 for workspace A
- partial sync 2: receive config 2 for workspace A
- (async): applying config 1 throws an error which is captured by the error tracker
- (async): config 2 can be applied successfully without an error
- partial sync 3: an error will be reported for a stale operation where instead it should be Running

With versioning, stale writes can be avoided. Rewriting the same scenario with versioning support:

- partial sync 1: receive config 1 for workspace A. Create entry in errorTracker for workspace A with version 1 and `err` nil
- partial sync 2: receive config 2 for workspace A. Update entry in errorTracker for workspace A with version 2 and `err` nil
- (async): applying config 1 throws an error. The errorTracker will ignore writes for workspace A with version 1 as the latest entry has a higher version
- (async): config 2 can be applied successfully without an error. Since there are no errors, entry for workspace A in errorTracker will be removed IFF the version is 2. This protects entries with a higher version from being removed by a different goroutine
- partial sync 3: no error is reported for a stale operation

The above example also works for cases where a stale error for an earlier action may override errors captured for a more recent action. Since the versions will be synchronously created per reconciliation cycle, they can be made to be monotonically increasing either by using an atomic counter.

## Changes in Rails

In case of errors in the request payload, the workspace `ActualState` should transition to `Error` if the DesiredState of the workspace has NOT changed since the last reconciliation cycle. This implies that the user has not initiated any other operation in between the reconciliation cycles and the error in the payload will correspond to the last action.

However, if the `DesiredState` of a workspace has been modified between reconciliation cycles, then it could be misleading to change the `ActualState` of the workspace to `Error` as the user may interpret this error to have been caused by the latest action. It should be ok for **the reconciliation API on Rails** to suppress the error details received in such cases and have the agentk carry out the latest instruction and focus on reporting its errors, if any.

**Example of the payload received in Rails:**

```json
{
  "update_type": "partial",
  "workspace_agent_infos": [{
    "name": "test-workspace",
    "namespace": "test-namespace",
    "latest_k8s_deployment_info": ..., // may or may not be populated alongside error_details
    "termination_progress": ..., // may or may not be populated alongside error_details
    "error_details": {
      "error_type": "applier",
      "error_message": "what went wrong while applying the configs"
    }
  }]
}
```

# Questions raised in the [issue](https://gitlab.com/gitlab-org/gitlab/-/issues/397001 "Agentk: Track agentk applier errors and send them to Rails")

> We would have to think about what would happen if there is some error in local tracker which has not yet been persisted in Rails and agentk restarted

Yes, this can happen. However, the restart would just result in the configuration being re-applied. If re-applying the resource doesn't create the error then reporting the lost error serves no purpose and the UI _should_ reflect the latest working state. If the applier fails again, then the error will be propagated to rails in the next reconciliation cycle.

> What kind of errors does the applier throw? Ideally what we are interested in are the errors when the kubernetes resources failed to apply for due to XYZ reasons. The applier errors pertain to errors when applying kubernetes sources derived from the provided devfile. In other words, we cannot expect the user to always understand what went wrong due to this translation layer as the user may only be aware of the devfile. So upfront, any sort of categorization would rely on being able to accurately classify errors into categories where the user response to an error may be different.

One way would be to figure out all/most types of errors the applier can throw and to then map them to a category of errors. However, after digging into the code, there are way too many possible errors that may be returned by the applier. In addition, some of the errors are parameterized as well, thereby complicating any categorization at the application layer. There is also a maintenance cost associated with having to update these categories of errors with lib upgrades if the error message content changes.

For the first iteration, perhaps the approach can be to just avoid such a classification of errors and limit the scope to just reporting the existence of an applier error. Even with such a limitation, a workspace user will have enough information (erroneous workspace id/name) to reach out to the cluster administrator to aid with troubleshooting.

> What is the structure of the Error field that we return? Is it just a string? Or do we need to pass some additional structured metadata, error type sub-field, etc?

Since there is a clear distinction between an Error and Failure, it would make sense to not use Failure for any of these cases. The payload can be of the form:

```protobuf
enum ErrorType {
  APPLIER = "applier";
}
```

```json
{
  "error_details": {
    "error_type": "applier",
    "error_message": "..."
  }
}
```

We can reserve 2 error types for starters: unknown type of error and an applier error. In the future, this can be extended to capture and deal with errors from other places.

> Do we need to consider other scenarios when designing the error field structure other than just applier errors? For example, any of the error scenarios described in Robust Error Handling and Logging ( &amp;10461)?

Perhaps all errors returned by `applyWorkspaceChanges` can be tracked and returned to Rails to indicate something going wrong at the time of making changes to kubernetes resources for a particular workspace

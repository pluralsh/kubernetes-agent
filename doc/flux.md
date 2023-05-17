# Flux Module

This document describes the current implementation of the [Flux module](internal/module/flux)
in the `gitlab-agent` project.

## Features

The Flux module implements the following features:

- Manage [Flux Receivers](https://fluxcd.io/flux/components/notification/receiver/) and their [associated secrets](https://fluxcd.io/flux/components/notification/receiver/#secret-reference) for [GitRepository resources](https://fluxcd.io/flux/components/source/gitrepositories/) that are referencing projects the same GitLab instance as the agent is connected to.

## Implementation Details

This section document the implementation details for the [aforementioned features](#features).

### Notification Receivers

The Flux module implements a [Kubernetes controller](https://kubernetes.io/docs/concepts/architecture/controller/)
that on a high-level watches for changes of the Flux GitRepository resources and if necessary
creates or updates associated notification receivers (and secrets for them). These notification receivers
install a webhook in the Flux notification controller (done by Flux itself) which are called whenever
GitLab detects a push on that particular Git repository.

The GitRepository controller does the following:

- watch for GitRepository resource changes (add and update) (delete is not necessary as Kubernetes garbage collection will clean up for us).
- watch for Receiver resource changes (add, update and delete).
- when a GitRepository resource event is triggered the controller will start a reconciliation for that particular object, which consists of:
  - check if the GitRepository references a GitLab project on the same instance as the agent is connected to (determined by the hostname).
  - create or update the Kubernetes Secret that is used to authenticate the receiver webhook caller
    - currently this secret is empty and no authenticate is done when calling the receiver webhook, nonetheless the secret is required by Flux.
  - create or update the Flux Receiver for that GitRepository.
  - both secret and receiver are created in the same namespace as the GitRepository
    and with an owner reference to that and have a 1:1 relationship which enables
    [Kubernetes garbage collection](https://kubernetes.io/docs/concepts/architecture/garbage-collection/) when the GitRepository object is deleted.
  - If any of the above fails and the error is not fatal that GitRepository is retried (rate limited and possibly delayed).
- when a Receiver resource event is triggered the controller will:
  - make sure that Receiver object is owned by the controller or else drop that resource
  - trigger an update for which projects should be reconciled - *reconciled* here refers what receiver webhooks to call -
    basically reconcile the GitRepository (the Flux reconciliation).
  - enqueue the owner GitRepository for reconciliation - may be that we need to reconcile the actual to the desired state again,
    because the object was modified or deleted.

#### Efficient lookup of Receivers by GitLab project

To efficiently lookup receivers by GitLab project (e.g. when a push event is received) a `project` index
is created on the receiver informer which maps GitLab project paths to receiver objects.
After a lookup those receiver objects are used to retrieve the webhook path that Flux created for the receiver object.

#### Access to GitLab projects

In the agent cluster a user may add GitRepository objects referencing any GitLab repository (or any repository for that matter, though we filter those out).
Among those there may be one that the agent doesn't have access to and for those we don't want to call the receiver webhook,
or we'd e.g. leak what projects exists on that particular instance and which are updated.

This *authorization* check is done on the server part which for a receiver Git push event checks back with GitLab if
a particular agent has access to that GitLab project before forwarding the event to the agent.
The checks back to GitLab are cached in Redis to limit load.

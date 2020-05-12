# gitlab-agent

GitLab Kubernetes Agent is an active in-cluster component for solving any GitLab<->Kubernetes integration tasks.

## Use cases and ideas

Below are some ideas that can be built using the agent.

* “Real-time” and resilient web hooks. Polling git repos scales poorly and so webhooks were invented. They remove polling, easing the load on infrastructure, and reduce the "event happened->it got noticed in an external system" latency. However, "webhooks" analog cannot work if cluster is behind a firewall. So an agent, runnning in the cluster, can connect to GitLab and receive a message when a change happens. Like web hooks, but the actual connection is initiated from the client, not from the server. Then the agent could:

  * Emulate a webhook inside of the cluster

  * Update a Kubernetes object with a new state. It can be a GitLab-specific object with some concrete schema about a git repository. Then we can have third-parties integrate with us via this object-based API. It can also be some integration-specific object.

* “Real-time” data access. Agent can stream requested data back to GitLab. See https://gitlab.com/gitlab-org/gitlab/-/issues/212810. 

* Feature/component discovery. GitLab may need a third-party component to be installed in a cluster for a particular feature to work. Agent can do that component discovery. E.g. we need Prometheus for metrics and we probably can find it in the cluster (is this a bad example? it illustrates the idea though).

* Better [GitOps](https://www.gitops.tech/) support. A repository can be used as a IaC repo. On successful CI run on the main repo, a commit is merged into that IaC repo. Commit describes the new desired state of infrastructure in a particular cluster (or clusters). An agent in a corresponding cluster(s) picks up the update and applies it to the objects in the cluster. We can work with Argo-cd/Flux here to try to reuse existing code and integrate with the community-built tools.

* “Infrastructure drift detection”. Monitor and alert on unexpected changes in Kubernetes objects that are managed in the IaC repo. Should support various ways to describe infrastructure (kustomize/helm/plain yaml/etc). 

* Preview changes to IaC specs against the current state of the corresponding cluster right in the MR. 

* “Live diff”. Building on top of the previous feature. In repo browser when a directory with IaC specs is opened, show a live comparison of what is in the repo and what is in the corresponding cluster. 

* Kubernetes has audit logs. We could build a page to view them and perhaps correlate with other GitLab events? 

* See how we can support https://github.com/kubernetes-sigs/application.

  * In repo browser detect resource specs with the defined annotations and show the relevant meta information bits
  * Have a panel showing live list of installed applications based on the annotations from the specification

* Kubernetes API proxying. Today GitLab cannot integrate with clusters behind a firewall. If we put an agent into such clusters and another agent next to GitLab, we can emulate Kubernetes API and proxy it into the actual cluster via the agents. The same approach can be used to build other features from this list. See the scheme below.
    
    ```mermaid
    graph TB
      agentk -- gRPC bidirectional streaming --> agentg
      
      subgraph "GitLab"
      agentg[agentg]
      GitLabRoR[GitLab RoR] -- gRPC --> agentg
      end
    
      subgraph "Kubernetes cluster"
      agentk[agentk]
      end  
    ```
    
    * `agentk` is our agent. It keeps a connection established to a GitLab instance. It waits for requests from it to process.
    
    * `agentg` is what accepts requests from `agentk`. It also listens for requests from `GitLab RoR`. The job of `agentg` is to match incoming requests from `GitLab RoR` with existing connections from `agentk`, forward the request to it and forward responses back.
    
    * `GitLab RoR` is the main GitLab application. It uses gRPC to talk to `agentg`. We could also support Kubernetes API to simplify migration of existing code onto this architecture. Could support both, depending on the need.
    
    [Bidirectional streaming](https://grpc.io/docs/guides/concepts/#bidirectional-streaming-rpc) is used between `agentk` and `agentg` to allow forwarding multiple concurrent requests though a single connection. This allows the connection acceptor i.e. gRPC server (`agentg`) to act as a client, sending requests as gRPC replies. Inverting client-server relationship is needed because the connection has to be initiated from the inside of the Kubernetes cluster i.e. from behind the firewall.

## Open questions and things to consider

### GitLab.com + `agentk`

We have CloudFlare CDN in front of GitLab.com. The connections that `agentk` establishes are long-running by design. It may or may not be in issue. See https://gitlab.com/groups/gitlab-com/gl-infra/-/epics/228

### HTTP/2 / gRPC to the backend

Using gRPC with CloudFlare CDN may or may not be an issue. [This comment](https://community.cloudflare.com/t/grpc-support/127798) suggests it is supported by [CloudFlare Spectrum](https://www.cloudflare.com/products/cloudflare-spectrum/) but [this tweet](https://twitter.com/prdonahue/status/1252886427475611650) says they are working on it.

Another potential issue is HAProxy that we use as our front door after CDN. We currently run 1.8 but HTTP/2-to-the-backend and hence gRPC-to-the-backend support [was added only in 1.9](https://www.haproxy.com/blog/haproxy-1-9-2-adds-grpc-support/). We'd need to upgrade to use this functionality.

### High availability and scalability

#### `agentk` - agent on the Kubernetes side

Multiple `Pod`s per deployment would be needed to have HA. This might be trivial but might require doing leader election, depending on the feature the agent will need to support. 

#### `agentg` - agent on the GitLab side

The difficulty of having multiple copies of the program is that only one of the copies has an active connection from a particular Kubernetes cluster. So to server a request from GitLab targeted at that cluster some sort of request routing would be needed. There are options (off the top of my head):

- [Consistent hashing](https://en.wikipedia.org/wiki/Consistent_hashing) could be used to:
  - Minimize disruptions if a copy of `agentg` goes missing
  - Route traffic to the correct copy
  
  We use [`nginx` as our ingress controller](https://docs.gitlab.com/charts/charts/nginx/index.html) and it does [support consistent hashing](https://www.nginx.com/resources/wiki/modules/consistent_hash/).

- We could have `nginx` ask one (any) of `agentg` where the right copy (the one that has the connection) running. `agentg` can do a lookup in Redis where each cluster connection is registered and return the address to `nginx` via [X-Accel-Redirect](https://www.nginx.com/resources/wiki/start/topics/examples/x-accel/#x-accel-redirect) header.

- We could teach `agentg` to gossip with its copies so that they tell each other where a connection for each cluster is. Each copy will know about all the connections and either proxy the request or use `X-Accel-Redirect` to redirect the traffic.

### What to build first?

What is the feature to prioritize first?

### `agentk` topologies

In a cluster `agentk` can be deployed:

- Cluster-wide deployment that works across all namespaces. Useful to manage cluster-wide state if we support GitOps.
- One or more per-namespace deployments, each concerned only with what is happening in a particular namespace. Note that namespace where `agentk` is deployed, and the namespace it's managing might be different namespaces.
- Both of the above at the same time.

Because of the above, each `agentk` copy must have its own identity for the system to be able to tell one from the other e.g. for troubleshooting reasons. Each copy should get a URL of `agentg` and fetch the configuration from it. In that configuration the agent should see if it's per-namespace or cluster-wide. Configuration is stored in a Git repo on GitLab to promote IaC approach.   

### Identity and authentication

`agentk` authenticates to GitLab using a token. That token also encodes agent's identity. It must be possible to rotate the token but keep the identity of the agent.

Each cluster has an identity too. The agent learns the identifier from the configuration it fetches. 

### Other items?

Please add more here.

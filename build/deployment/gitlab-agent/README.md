# GitLab Agent for Kubernetes package

## Description

Package for installing GitLab Agent on a Kubernetes cluster.

`agentk` is the GitLab Agent. It keeps a connection established to a GitLab instance, waiting for requests to process. It may also actively send information about things happening in the cluster.

## Prerequisites

- [Kustomize](https://kustomize.io/)
- [`kpt`](https://kpt.dev/book/01-getting-started/01-system-requirements) version v1.0+
  - for pre-v1.0 versions, see *link* (previous package management with kpt [v0.39.3](https://github.com/GoogleContainerTools/kpt/releases/tag/v0.39.3))
- [Git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
- [Docker](https://docs.docker.com/get-docker/); specifically, the `docker-cli`
- `cluster-admin` access to a Kubernetes cluster

## Configuration

GitLab Agent needs two pieces of configuration to connect to a GitLab instance:

1. URL - something like `wss://kas.gitlab.com`.

    > The agent can use WebSockets or gRPC protocols to connect to GitLab. Depending on how your GitLab instance is configured, you may need to use one or the other.
    >
    > - Specify `grpc` scheme (e.g. `grpc://127.0.0.1:8150`) to use gRPC. The connection is **not encrypted**.
    > - Specify `grpcs` scheme to use an encrypted gRPC connection.
    > - Specify `ws` scheme to use WebSocket connection. The connection is **not encrypted**.
    > - Specify `wss` scheme to use an encrypted WebSocket connection.

1. Token - obtained through UI when registering the agent. Go to "Infrastructure -> Kubernetes clusters" in the menu to register an agent.

## Deploy the package

1. Get the package.

    ```shell
    kpt pkg get https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent.git/build/deployment/gitlab-agent gitlab-agent
    ```

   - *If you are not using `kpt`, clone the repository*:

     ```shell
     git clone https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent.git
     cd gitlab-agent/build/deployment/gitlab-agent
     ```

1. (Optional) Edit the package if you want to change the agent name, the namespace where the agent is deployed, the image version, etc... This step is also necessary if you are installing multiple agents in your cluster.

    - Commit the initial package to version control.
    - Edit values in the `kpt-setter-configmap.yaml` file:

        | Name                    | Default Value                                                                   | Description                                                        |
        | ----------------------- | ------------------------------------------------------------------------------- | ------------------------------------------------------------------ |
        | **agent-image-ref**     | `registry.gitlab.com/gitlab-org/cluster-integration/gitlab-agent/agentk:stable` | Image ref name and tag                                             |
        | **kas-args**            | `- --token-file=/config/token\n- --kas-address\n- wss://kas.gitlab.com`         | Image args for agentk container                                    |
        | **name-prefix**         | `""`                                                                            | Prefix for resource names (multiple agents must have unique names) |
        | **namespace**           | `gitlab-agent`                                                                  | Namespace to install Agent into                                    |
        | **prometheus-scrape**   | `true`                                                                          | Enable or disable scraping of agentk metrics                       |
        | **serviceaccount-name** | `gitlab-agent`                                                                  |                                                                    |

    - Apply the change(s):

        ```shell
        kpt fn render
        ```

1. Add the agent token.

    - Write the token you obtained from the agent registration step in GitLab to `secrets/agent.token`

        ```shell
        echo -n "<agent token>" > secrets/agent.token
        ```

        *Note that this should not be committed to a repository and is excluded by the `.gitignore` in the `secrets` directory by default.*

1. Initialize the package for deployment. This step will create `resourcegroup.yaml`.

    ```shell
    kpt live init
    ```

1. (Optional but recommended) Commit changes to version control with a new version tag.

1. Apply the package to the cluster.

    ```shell
    kustomize build | kpt live apply - --reconcile-timeout=2m --install-resource-group --server-side --show-status-events
    ```

## Updating the Agent

1. Retrieve the latest version of the package.

    ```shell
    kpt pkg update gitlab-agent@{newVersion}
    ```

1. (Optional but recommended) Commit changes to version control with a new version tag.

1. Apply the package to the cluster.

    ```shell
    kustomize build | kpt live apply - --reconcile-timeout=2m --server-side --show-status-events
    ```

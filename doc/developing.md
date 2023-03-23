# Development guide

[[_TOC_]]

## Repository overview

Most up-to-date video describing how this repository is
structured:

[![GitLab Agent repository overview - It's time to Go! Episode 5](https://img.youtube.com/vi/Mh7PG4_cBxI/0.jpg)](https://www.youtube.com/watch?v=Mh7PG4_cBxI "GitLab Agent repository overview")

### Past recordings

#### 2021-11-20

[![GitLab Agent repository overview](http://img.youtube.com/vi/j8CyaCWroUY/1.jpg)](http://www.youtube.com/watch?v=j8CyaCWroUY "GitLab Agent repository overview")

## Running kas and agentk locally

[![GitLab Agent development environment setup](https://img.youtube.com/vi/UWptMO-Amtc/0.jpg)](https://www.youtube.com/watch?v=UWptMO-Amtc "GitLab Agent development environment setup")

You can run `kas` and `agentk` locally to test the Agent yourself.

1. Create a `token.txt`. This is the token for
   [the agent you created](https://docs.gitlab.com/ee/user/clusters/agent/index.html#create-an-agent-record-in-gitlab).
   This file must not contain a newline character. You can create the file with this command:

   ```shell
   echo -n "<TOKEN>" > token.txt
   ```

1. [Setup GDK](https://gitlab.com/gitlab-org/gitlab-development-kit#installation).

1. [Setup kas in GDK](https://gitlab.com/gitlab-org/gitlab-development-kit/-/blob/main/doc/howto/kubernetes_agent.md).

1. Start the binaries with the following commands:

   ```shell
   # Start GitLab but stop GDK's version of kas.
   gdk start && gdk stop gitlab-k8s-agent

   # Let kas know it's own private API url
   # This is needed for CI tunnel and any other reverse tunnel-based features.
   export OWN_PRIVATE_API_URL=grpc://127.0.0.1:8155
   # Start kas
   bazel run //cmd/kas -- --configuration-file="gdk/dir/gitlab-k8s-agent-config.yml"
   ```

1. In a new terminal window, run this command to start `agentk`:

   ```shell
   # These are used for leader election, etc. Make sure the namespace exists in the cluster.
   export POD_NAMESPACE=ns
   export POD_NAME=agent1
   kubectl create ns "$POD_NAMESPACE"

   # Set --kas-address correctly, depending on how kas is setup.
   # Set --context to a kubectl context to use. Can be omitted to use the current context, but that is risky
   # as the behavior is not deterministic in that case.
   # Get the list of contexts with kubectl config get-contexts
   bazel run //cmd/agentk -- --kas-address=grpc://127.0.0.1:8150 --token-file="$(pwd)/token.txt" --context=minikube
   ```

You can also inspect the [Makefile](../Makefile) for more targets. You can run `kas` and/or `agentk` from your
editor too. Just make sure to setup

## Run tests locally

You can run all tests, or a subset of tests, locally.

- **To run all tests**: Run the command `make test`.
- **To run all test targets in the directory**: Run the command
  `bazel test //internal/module/gitops/server:all`.

  You can use `*` in the command, instead of `all`, but it must be quoted to
  avoid shell expansion: `bazel test '//internal/module/gitops/server:*'`.
- **To run all tests in a directory and its subdirectories**: Run the command
  `bazel test //internal/module/gitops/server/...`.

### Run specific test scenarios

To run only a specific test scenario, you need the directory name and the target
name of the test. For example, to run the tests at
`internal/module/gitops/server/module_test.go`, the `BUILD.bazel` file that
defines the test's target name lives at `internal/module/gitops/server/BUILD.bazel`.
In the latter, the target name is defined like:

```bazel
go_test(
    name = "server_test",
    size = "small",
    srcs = [
        "module_test.go",
```

The target name is `server_test` and the directory is `internal/module/gitops/server/`.
Run the test scenario with this command:

```shell
bazel test //internal/module/gitops/server:server_test
```

### Additional resources

- Bazel documentation about [specifying targets to build](https://docs.bazel.build/versions/master/guide.html#specifying-targets-to-build).
- [The Bazel query](https://docs.bazel.build/versions/master/query.html)
- [Bazel query how to](https://docs.bazel.build/versions/master/query-how-to.html)


## Debugging locally

For local debugging we don't use Bazel and instead build agent "from source".

### dlv

Debug `agentk` with the following command:

```sh
export POD_NAMESPACE=default
export POD_NAME=agentk

dlv cmd/agentk/main.go -- \
    --kas-address "${kas_address}" \
    --token-file "${token_file}"
```

Debug `kas` with the following command:

```sh
export OWN_PRIVATE_API_URL="$(gdk config get gitlab_k8s_agent.__private_api_url)"
dlv cmd/kas/main.go -- --configuration-file "$(gdk config get gitlab_k8s_agent.__config_file)"
```

### VS Code

To debug in VS Code, use the following [Launch Configuration](https://code.visualstudio.com/docs/editor/debugging#_launch-configurations). Replace `<path-to-your-gdk>` with the full path to your GDK (you can find that with `gdk config get gitlab_k8s_agent.__config_file`).

```json
{
    // Use IntelliSense to learn about possible attributes.
    // Hover to view descriptions of existing attributes.
    // For more information, visit: https://go.microsoft.com/fwlink/?linkid=830387
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch agentk",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/agentk/",
            "args": ["--kas-address","grpc://172.16.123.1:8150", "--token-file", "${workspaceFolder}/token.txt"],
            "env": {
                "POD_NAMESPACE": "default",
                "POD_NAME": "agentk"
            }
        },
        {
            "name": "Launch kas",
            "type": "go",
            "request": "launch",
            "mode": "auto",
            "program": "${workspaceFolder}/cmd/kas/",
            "args": ["--configuration-file", "<path-to-your-gdk>/gitlab-k8s-agent-config.yml"],
            "env": {
                "OWN_PRIVATE_API_URL": "grpc://172.16.123.1:8155",
            }
        }
    ]
}
```

### JetBrains GoLand

Add the following run/debug configurations:

For `kas`:

![kas run configuration](https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/uploads/12c3ad3f7c92a6d5ce2b5380ef4be5a2/Screen_Shot_2023-01-10_at_9.19.36_AM.png)

For `agentk`:

![agentk run configuration](https://gitlab.com/gitlab-org/cluster-integration/gitlab-agent/uploads/5aadad03c98c4136f94bfa7c702c4725/Screen_Shot_2023-01-10_at_9.19.46_AM.png)

It's optional, but consider also specifying `--context=<desired context>` command line argument to not depend on the currently selected context.
## kas QA tests

This section describes how to run kas tests against different GitLab environments based on the
[GitLab QA orchestrator](https://gitlab.com/gitlab-org/gitlab-qa).

### Status

The `kas` QA tests currently have some limitations. You can run them manually on GDK, but they don't
run automatically with the nightly jobs against the live environment. See the section below
to learn how to run them against different environments.

### Prepare

Before performing any of these tests, if you have a `k3s` instance running, make sure to
stop it manually before running them. Otherwise, the tests might fail with the message
`failed to remove k3s cluster`.

You might need to specify the correct Agent image version that matches the `kas` image version. You can use the `GITLAB_AGENTK_VERSION` local environment for this.

### Against `staging`

1. Go to your local `qa/qa/service/cluster_provider/k3s.rb` and comment out
   [this line](https://gitlab.com/gitlab-org/gitlab/-/blob/5b15540ea78298a106150c3a1d6ed26416109b9d/qa/qa/service/cluster_provider/k3s.rb#L8) and
   [this line](https://gitlab.com/gitlab-org/gitlab/-/blob/5b15540ea78298a106150c3a1d6ed26416109b9d/qa/qa/service/cluster_provider/k3s.rb#L36).
   We don't allow local connections on `staging` as they require an admin user.
1. Ensure you don't have an `EE_LICENSE` environment variable set as this would force an admin login.
1. Go to your GDK root folder and `cd gitlab/qa`.
1. Login with your user in staging and create a group to be used as sandbox.
   Something like: `username-qa-sandbox`.
1. Create an access token for your user with the `api` permission.
1. Replace the values given below with your own and run:

   ```shell
   GITLAB_SANDBOX_NAME="<THE GROUP ID YOU CREATED ON STEP 2>" \
   GITLAB_QA_ACCESS_TOKEN="<THE ACCESS TOKEN YOU CREATED ON STEP 3>" \
   GITLAB_USERNAME="<YOUR STAGING USERNAME>" \
   GITLAB_PASSWORD="<YOUR STAGING PASSWORD>" \
   bundle exec bin/qa Test::Instance::All https://staging.gitlab.com -- --tag quarantine qa/specs/features/ee/api/7_configure/kubernetes/kubernetes_agent_spec.rb
   ```

### Against GDK

1. Go to your `qa/qa/fixtures/kubernetes_agent/agentk-manifest.yaml.erb` and comment out [this line](https://gitlab.com/gitlab-org/gitlab/-/blob/a55b78532cfd29426cf4e5b4edda81407da9d449/qa/qa/fixtures/kubernetes_agent/agentk-manifest.yaml.erb#L27) and uncomment [this line](https://gitlab.com/gitlab-org/gitlab/-/blob/a55b78532cfd29426cf4e5b4edda81407da9d449/qa/qa/fixtures/kubernetes_agent/agentk-manifest.yaml.erb#L28).
   GDK's `kas` listens on `grpc`, not on `wss`.
1. Go to the GDK's root folder and `cd gitlab/qa`.
1. On the contrary to staging, run the QA test in GDK as admin, which is the default choice. To do so, use the default sandbox group and run the command below. Make sure to adjust your credentials if necessary, otherwise, the test might fail:

   ```shell
   GITLAB_USERNAME=root \
   GITLAB_PASSWORD="5iveL\!fe" \
   GITLAB_ADMIN_USERNAME=root \
   GITLAB_ADMIN_PASSWORD="5iveL\!fe" \
   bundle exec bin/qa Test::Instance::All http://gdk.test:3000 -- --tag quarantine qa/specs/features/ee/api/7_configure/kubernetes/kubernetes_agent_spec.rb
   ```

## Optimizing build performance

Bazel creates a lot of files during the build:
- For [sandboxing](https://docs.bazel.build/versions/main/sandboxing.html)
  purposes. These are temporary files, not taking up disk space permanently.
- Cache of completed actions such as compiled packages, ready for linking.

To speed up Bazel builds on your machine, you can put all those files onto a RAM disk. On Linux, you can use
`/dev/shm` (where it's mounted depends on your distribution). On macOS you can create a 10 GiB RAM disk and
mount it under `/Volumes/ramdisk` using the following command:

```shell
diskutil partitionDisk $(hdiutil attach -nomount ram://20971520) 1 GPTFormat APFS 'ramdisk' '100%'
```

10 GiB should be enough, although more space may be required by this project in the future. This disk will be gone once you reboot your Mac. Just recreate it.

To use the RAM disk, specify the path to it. There are two options:

- **Move sandbox only**. This requires less space i.e. less RAM. Use
  [`--sandbox_base`](https://docs.bazel.build/versions/main/command-line-reference.html#flag--sandbox_base)
  option.

- **Move all temporary files**, including sandbox and cache. Use
  [`--output_user_root`](https://docs.bazel.build/versions/main/command-line-reference.html#flag--output_user_root)
  option. See
  [Choosing the output base](https://docs.bazel.build/versions/main/guide.html#choosing-the-output-base).

Put the option you chose into your `~/.bazelrc` like that:

```plaintext
#build --sandbox_base=/Volumes/ramdisk
startup --output_user_root=/Volumes/ramdisk
```

[Docs about the `.bazelrc` file](https://docs.bazel.build/versions/main/guide.html#bazelrc-the-bazel-configuration-file).

# Development guide

## Repository overview

You can watch the video below to understand how this repository is
structured:

[![GitLab Agent repository overview](http://img.youtube.com/vi/j8CyaCWroUY/0.jpg)](http://www.youtube.com/watch?v=j8CyaCWroUY "GitLab Agent repository overview")

## Running kas and agentk locally

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
   # Set --kas-address correctly, depending on how kas is setup.
   bazel run //cmd/agentk -- --kas-address=grpc://127.0.0.1:8150 --token-file="$(pwd)/token.txt"
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

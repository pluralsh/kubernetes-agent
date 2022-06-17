# Release Process

1. On the 15th, an automatic Slack notification reminds the Configure team to create a monthly release.
1. On the 15th, tag a new version that matches the upcoming GitLab minor version. E.g. If the upcoming milestone is 13.7,
   then tag `v13.7.0`.
1. Make a release of the [gitlab-agent chart](https://gitlab.com/gitlab-org/charts/gitlab-agent#publishing-a-new-release):
   - In `Chart.yaml`
      - Update `appVersion` to be the exactly the same as the version tag in `cluster-integrations/gitlab-agent`
      - Bump `version` according to semantic versioning. For a GitLab milestone release, this will generally be the minor version.
   - Push a tag `vX.Y.Z` where `X.Y.Z` is the `version` in `Chart.yaml`
1. The [`GITLAB_KAS_VERSION`](https://gitlab.com/gitlab-org/gitlab/-/blob/master/GITLAB_KAS_VERSION) file in
   the GitLab rails monolith is updated to that new tag in a new MR.
   This MR should be accepted by the maintainer with no questions asked, typically.
   [Example MR](https://gitlab.com/gitlab-org/gitlab/-/merge_requests/74462).
1. Wait for the MR to get deployed to a .com environment (can be pre, gstg, gprd etc; this is shown in the MR widget).
1. Find the latest image built of KAS in the [dev.gitlab.org registry]( https://dev.gitlab.org/gitlab/charts/components/images/container_registry/426?orderBy=NAME&sort=asc):
   - Go to [`#releases` in Slack](https://gitlab.slack.com/archives/C0XM5UU6B) (internal link)
   - Look at the latest comment from `auto-deploy-bot`. It should be something like:

     > New auto-deploy branch: `14-10-auto-deploy-2022041215`

     Pick out the timestamp from there, e.g. `2022041215` and enter it into the search box. This should give the tag, e.g. `dev.gitlab.org:5005/gitlab/charts/components/images/gitlab-kas:14-10-202204121520-1ad684ad2e5`.
   - Ensure that the version is correct. Running the image with `--help` should match `GITLAB_KAS_VERSION` from `gitlab-org/gitlab`. For example:

     ```shell
     docker login dev.gitlab.org:5005 # use your username and a personal access token with the read_registry scope
     docker run --rm dev.gitlab.org:5005/gitlab/charts/components/images/gitlab-kas:14-10-202204121520-1ad684ad2e5 --version
     ```

     outputs

     ```
     kas version v14.10.0-rc2, commit: 68ae893, built: 20220412.152619
     ```

1. Make an MR to update `kas` image tag for `pre` and `gstg` in
   https://gitlab.com/gitlab-com/gl-infra/k8s-workloads/gitlab-com, and get an SRE to deploy that MR.
   [Example MR](https://gitlab.com/gitlab-com/gl-infra/k8s-workloads/gitlab-com/-/merge_requests/1318).
1. We should test the new version of `kas` on `gstg` with a real agent. An end-to-end QA test with a real agent
   and GitLab that runs automatically as part of the nightly QA process.
   is [planned](https://gitlab.com/groups/gitlab-org/-/epics/4949). For now, you can use `gstg` + an agent in
   a Kubernetes cluster, including locally running `agentk` and cluster.
1. To get the new change deployed to `gprd`, follow the
   [Change Request Workflows](https://about.gitlab.com/handbook/engineering/infrastructure/change-management/#change-request-workflows):
   - Make an MR to update `kas` image tag in `gprd` in
     https://gitlab.com/gitlab-com/gl-infra/k8s-workloads/gitlab-com.
     [Example MR](https://gitlab.com/gitlab-com/gl-infra/k8s-workloads/gitlab-com/-/merge_requests/1319).
   - Open a [production change issue](https://gitlab.com/gitlab-com/gl-infra/production/-/issues) to get
     the above MR merged and rolled out. [Example issue](https://gitlab.com/gitlab-com/gl-infra/production/-/issues/7266).
1. If there are breaking changes to the `kas` config file, then MRs need to be raised for
   [Omnibus](https://gitlab.com/gitlab-org/omnibus-gitlab/), and
   [charts](https://gitlab.com/gitlab-org/charts/gitlab/).
1. On the 22nd, the `GITLAB_KAS_VERSION` is automatically synced for when release team tags new Omnibus and chart releases.

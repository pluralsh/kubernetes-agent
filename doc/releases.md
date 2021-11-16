# Release Process

1. On the 15th, an automatic Slack notification reminds the Configure team to create a monthly release.
1. On the 15th we should always tag a new version that matches the upcoming GitLab minor version. E.g. If GitLab 13.7
   will be released on the 22nd, then we should tag `v13.7.0`.
1. The [`GITLAB_KAS_VERSION`](https://gitlab.com/gitlab-org/gitlab/-/blob/master/GITLAB_KAS_VERSION) file in
   the GitLab rails monolith is updated to that new tag in a new MR.
   This MR should be accepted by the maintainer with no questions asked, typically.
   [Example MR](https://gitlab.com/gitlab-org/gitlab/-/merge_requests/74462).
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
     the above MR merged and rolled out. [Example issue](https://gitlab.com/gitlab-com/gl-infra/production/-/issues/5821).
1. If there are breaking changes to the `kas` config file, then MRs need to be raised for
   [Omnibus](https://gitlab.com/gitlab-org/omnibus-gitlab/), and
   [charts](https://gitlab.com/gitlab-org/charts/gitlab/).
1. On the 22nd, the `GITLAB_KAS_VERSION` is automatically synced for when release team tags new Omnibus and chart releases.

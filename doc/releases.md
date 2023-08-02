# Release Process

1. On the 15th, an automatic Slack notification reminds the Environments team to create a monthly release.
1. On the 15th, tag a new version that matches the upcoming GitLab minor version. E.g. If the upcoming milestone is 13.7,
   then tag `v13.7.0`.
1. Make a release of the [gitlab-agent chart](https://gitlab.com/gitlab-org/charts/gitlab-agent#publishing-a-new-release):
   - In `Chart.yaml`
      - Update `appVersion` to be the exactly the same as the version tag in `cluster-integrations/gitlab-agent`
      - Bump `version` according to semantic versioning. For a GitLab milestone release, this will generally be the minor version.
   - Push a tag `vX.Y.Z` where `X.Y.Z` is the `version` in `Chart.yaml`
1. (Automated) The [`GITLAB_KAS_VERSION`](https://gitlab.com/gitlab-org/gitlab/-/blob/master/GITLAB_KAS_VERSION) file in
   the GitLab rails monolith is updated to that new tag in a new MR.
   This MR should be accepted by the maintainer with no questions asked, typically.
   [Example MR](https://gitlab.com/gitlab-org/gitlab/-/merge_requests/111845). [List of bot-created MRs](https://gitlab.com/gitlab-org/gitlab/-/merge_requests?scope=all&state=all&label_name[]=group%3A%3Aconfigure&author_username=gitlab-dependency-update-bot). [Bot configuration](https://gitlab.com/gitlab-org/frontend/renovate-gitlab-bot/-/blob/main/renovate/gitlab/version-files.config.js).
1. Wait for the MR to get deployed to a .com environment (can be pre, gstg, gprd etc; this is shown in the MR widget).
1. Find the latest image built of KAS in the [dev.gitlab.org registry]( https://dev.gitlab.org/gitlab/charts/components/images/container_registry/426?orderBy=NAME&sort=asc) (*Note: the images are built by the `gitlab-kas` job in [this project pipeline](https://gitlab.com/gitlab-org/security/charts/components/images)*):
   - Construct a timestamp from today's or yesterday's date according to the `YYYYMMDD` format (e.g. `20230316`) and enter it into the search box. This should give you a list of tags published on that day.
     Pick the latest one, e.g. `dev.gitlab.org:5005/gitlab/charts/components/images/gitlab-kas:15-10-202303160020-bc2cbbf9e9d`.
   - Ensure that the version is correct. Running the image with `--version` should match `GITLAB_KAS_VERSION` from `gitlab-org/gitlab`. For example:

     ```shell
     docker login dev.gitlab.org:5005 # use your username and a personal access token with the read_registry scope
     docker run --rm dev.gitlab.org:5005/gitlab/charts/components/images/gitlab-kas:15-10-202303160020-bc2cbbf9e9d --version
     ```

     outputs

     ```
     kas version v15.10.0, commit: v15.10.0, built: 20230316.002657
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
1. If there are breaking changes to the `kas` config file, then MRs need to be raised for
   [Omnibus](https://gitlab.com/gitlab-org/omnibus-gitlab/), and
   [charts](https://gitlab.com/gitlab-org/charts/gitlab/).
1. On the 22nd, the `GITLAB_KAS_VERSION` is automatically synced for when release team tags new Omnibus and chart releases.

## Troubleshooting

### Image with desired tag is not on `dev.gitlab.org`

The image deployed to `dev.gitlab.org` is built within the pipeline of the 
[security images project](https://gitlab.com/gitlab-org/security/charts/components/images) in a 3 hour frequency. 
The version to built is taken from the [`GITLAB_KAS_VERSION`](https://gitlab.com/gitlab-org/gitlab/-/blob/master/GITLAB_KAS_VERSION).

### Image Pull Error in gstg / dev cluster

This may happen because the image from `dev.gitlab.org` is not yet synced with the high availability registry,
where the KAS image is eventually pulled from. This registry may lag 1.5-2 hours behind.
This sync happens in the `sync-images-artifact-registry` job of the pipeline in the 
[security images project](https://gitlab.com/gitlab-org/security/charts/components/images).
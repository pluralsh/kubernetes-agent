module gitlab.com/gitlab-org/cluster-integration/gitlab-agent/v15

go 1.17

require (
	cloud.google.com/go/profiler v0.3.0
	github.com/aquasecurity/starboard v0.15.8
	github.com/ash2k/stager v0.3.0
	github.com/bmatcuk/doublestar/v2 v2.0.4
	github.com/envoyproxy/protoc-gen-validate v0.6.13
	github.com/getsentry/sentry-go v0.14.0
	github.com/go-logr/zapr v1.2.3
	github.com/go-redis/redis/v8 v8.11.5
	github.com/go-redis/redismock/v8 v8.0.6
	github.com/golang-jwt/jwt/v4 v4.4.2
	github.com/golang/mock v1.6.0
	github.com/google/go-cmp v0.5.9
	github.com/grpc-ecosystem/go-grpc-middleware v1.3.0
	github.com/grpc-ecosystem/go-grpc-prometheus v1.2.1-0.20200507082539-9abf3eb82b4a
	github.com/hashicorp/go-retryablehttp v0.7.1
	github.com/piotrkowalczuk/promgrpc/v4 v4.0.4
	github.com/prometheus/client_golang v1.13.0
	github.com/robfig/cron/v3 v3.0.1
	github.com/spf13/cobra v1.6.0
	github.com/stretchr/testify v1.8.0
	gitlab.com/gitlab-org/gitaly/v15 v15.4.2
	gitlab.com/gitlab-org/security-products/analyzers/report/v3 v3.7.1
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.36.2
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.36.2
	go.opentelemetry.io/otel v1.10.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.10.0
	go.opentelemetry.io/otel/metric v0.32.2
	go.opentelemetry.io/otel/sdk v1.10.0
	go.opentelemetry.io/otel/trace v1.10.0
	go.uber.org/zap v1.23.0
	golang.org/x/net v0.0.0-20221004154528-8021a29435af
	golang.org/x/sync v0.0.0-20220929204114-8fcdb60fdcc0
	golang.org/x/time v0.0.0-20220922220347-f3bd1da661af
	golang.org/x/tools v0.1.12
	google.golang.org/api v0.98.0
	google.golang.org/genproto v0.0.0-20221010155953-15ba04fc1c0e
	google.golang.org/grpc v1.50.0
	google.golang.org/protobuf v1.28.1
	helm.sh/helm/v3 v3.9.4
	k8s.io/api v0.24.6
	k8s.io/apimachinery v0.24.6
	k8s.io/apiserver v0.24.6
	k8s.io/cli-runtime v0.24.6
	k8s.io/client-go v0.24.6
	k8s.io/klog/v2 v2.70.1
	k8s.io/kubectl v0.24.6
	k8s.io/utils v0.0.0-20220728103510-ee6ede2d64ed
	nhooyr.io/websocket v1.8.7
	sigs.k8s.io/cli-utils v0.33.0
	sigs.k8s.io/controller-runtime v0.12.3
	sigs.k8s.io/yaml v1.3.0
)

require (
	cloud.google.com/go v0.102.0 // indirect
	cloud.google.com/go/compute v1.7.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20210617225240-d185dfc1b5a1 // indirect
	github.com/BurntSushi/toml v1.2.0 // indirect
	github.com/MakeNowJust/heredoc v0.0.0-20170808103936-bb23615498cd // indirect
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver/v3 v3.1.1 // indirect
	github.com/Masterminds/sprig/v3 v3.2.2 // indirect
	github.com/Masterminds/squirrel v1.5.3 // indirect
	github.com/OneOfOne/xxhash v1.2.8 // indirect
	github.com/agnivade/levenshtein v1.1.1 // indirect
	github.com/asaskevich/govalidator v0.0.0-20200428143746-21a406dcc535 // indirect
	github.com/beevik/ntp v0.3.0 // indirect
	github.com/benbjohnson/clock v1.3.0 // indirect
	github.com/beorn7/perks v1.0.1 // indirect
	github.com/caarlos0/env/v6 v6.10.0 // indirect
	github.com/cenkalti/backoff/v4 v4.1.3 // indirect
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/chai2010/gettext-go v0.0.0-20160711120539-c6fed771bfd5 // indirect
	github.com/containerd/containerd v1.6.8 // indirect
	github.com/cyphar/filepath-securejoin v0.2.3 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/docker/cli v20.10.17+incompatible // indirect
	github.com/docker/distribution v2.8.1+incompatible // indirect
	github.com/docker/docker v20.10.17+incompatible // indirect
	github.com/docker/docker-credential-helpers v0.6.4 // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-metrics v0.0.1 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/emicklei/go-restful/v3 v3.8.0 // indirect
	github.com/evanphx/json-patch v5.6.0+incompatible // indirect
	github.com/exponent-io/jsonpath v0.0.0-20151013193312-d6023ce2651d // indirect
	github.com/fatih/camelcase v1.0.0 // indirect
	github.com/fatih/color v1.13.0 // indirect
	github.com/felixge/httpsnoop v1.0.3 // indirect
	github.com/fsnotify/fsnotify v1.5.4 // indirect
	github.com/fvbommel/sortorder v1.0.1 // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/go-errors/errors v1.4.2 // indirect
	github.com/go-gorp/gorp/v3 v3.0.2 // indirect
	github.com/go-logr/logr v1.2.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/jsonreference v0.20.0 // indirect
	github.com/go-openapi/swag v0.21.1 // indirect
	github.com/gobwas/glob v0.2.3 // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20210331224755-41bb18bfe9da // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/gnostic v0.6.9 // indirect
	github.com/google/go-containerregistry v0.11.0 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/google/pprof v0.0.0-20220412212628-83db2b799d1f // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.1.0 // indirect
	github.com/googleapis/gax-go/v2 v2.4.0 // indirect
	github.com/gorhill/cronexpr v0.0.0-20180427100037-88b0669f7d75 // indirect
	github.com/gorilla/mux v1.8.0 // indirect
	github.com/gosuri/uitable v0.0.4 // indirect
	github.com/gregjones/httpcache v0.0.0-20190611155906-901d90724c79 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.7.0 // indirect
	github.com/hashicorp/go-cleanhttp v0.5.1 // indirect
	github.com/hashicorp/yamux v0.1.1 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.13 // indirect
	github.com/inconshreveable/mousetrap v1.0.1 // indirect
	github.com/jmoiron/sqlx v1.3.5 // indirect
	github.com/jonboulle/clockwork v0.2.2 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/compress v1.15.11 // indirect
	github.com/lann/builder v0.0.0-20180802200727-47ae307949d0 // indirect
	github.com/lann/ps v0.0.0-20150810152359-62de8c46ede0 // indirect
	github.com/lib/pq v1.10.6 // indirect
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.16 // indirect
	github.com/mattn/go-runewidth v0.0.9 // indirect
	github.com/matttproud/golang_protobuf_extensions v1.0.2-0.20181231171920-c182affec369 // indirect
	github.com/mitchellh/copystructure v1.2.0 // indirect
	github.com/mitchellh/go-wordwrap v1.0.0 // indirect
	github.com/mitchellh/reflectwalk v1.0.2 // indirect
	github.com/moby/locker v1.0.1 // indirect
	github.com/moby/spdystream v0.2.0 // indirect
	github.com/moby/term v0.0.0-20210619224110-3f7ff695adc6 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/monochromegane/go-gitignore v0.0.0-20200626010858-205db1a8cc00 // indirect
	github.com/morikuni/aec v1.0.0 // indirect
	github.com/munnerz/goautoneg v0.0.0-20191010083416-a7dc8b61c822 // indirect
	github.com/oklog/ulid/v2 v2.0.2 // indirect
	github.com/open-policy-agent/opa v0.44.0 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.0.3-0.20220114050600-8b9d41f48198 // indirect
	github.com/opentracing/opentracing-go v1.2.0 // indirect
	github.com/pelletier/go-toml v1.9.5 // indirect
	github.com/peterbourgon/diskv v2.0.1+incompatible // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/prometheus/client_model v0.2.0 // indirect
	github.com/prometheus/common v0.37.0 // indirect
	github.com/prometheus/procfs v0.8.0 // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/rubenv/sql-migrate v1.1.2 // indirect
	github.com/russross/blackfriday v1.6.0 // indirect
	github.com/sebest/xff v0.0.0-20210106013422-671bd2870b3a // indirect
	github.com/shopspring/decimal v1.2.0 // indirect
	github.com/sirupsen/logrus v1.9.0 // indirect
	github.com/spf13/cast v1.4.1 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/spyzhov/ajson v0.4.2 // indirect
	github.com/tchap/go-patricia/v2 v2.3.1 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/xlab/treeprint v0.0.0-20181112141820-a009c3971eca // indirect
	github.com/yashtewari/glob-intersection v0.1.0 // indirect
	gitlab.com/gitlab-org/labkit v1.16.0 // indirect
	gitlab.com/gitlab-org/security-products/analyzers/common/v2 v2.24.0 // indirect
	gitlab.com/gitlab-org/security-products/analyzers/ruleset v1.0.0 // indirect
	go.opencensus.io v0.23.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/internal/retry v1.10.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.10.0 // indirect
	go.opentelemetry.io/proto/otlp v0.19.0 // indirect
	go.starlark.net v0.0.0-20200306205701-8dd3e2ee1dd5 // indirect
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	golang.org/x/crypto v0.0.0-20220926161630-eccd6366d1be // indirect
	golang.org/x/mod v0.6.0-dev.0.20220419223038-86c51ed26bb4 // indirect
	golang.org/x/oauth2 v0.0.0-20220822191816-0ebed06d0094 // indirect
	golang.org/x/sys v0.0.0-20220928140112-f11e5e49a4ec // indirect
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467 // indirect
	golang.org/x/text v0.3.7 // indirect
	gomodules.xyz/jsonpatch/v2 v2.2.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	gopkg.in/inf.v0 v0.9.1 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/apiextensions-apiserver v0.24.2 // indirect
	k8s.io/component-base v0.24.6 // indirect
	k8s.io/kube-openapi v0.0.0-20220627174259-011e075b9cb8 // indirect
	oras.land/oras-go v1.2.0 // indirect
	sigs.k8s.io/json v0.0.0-20220525155127-227cbc7cc124 // indirect
	sigs.k8s.io/kustomize/api v0.11.4 // indirect
	sigs.k8s.io/kustomize/kyaml v0.13.6 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.1 // indirect
)

// same version as used by rules_go to maintain compatibility with patches - see the WORKSPACE file
replace golang.org/x/tools => golang.org/x/tools v0.1.12

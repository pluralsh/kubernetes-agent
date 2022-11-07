load("@bazel_gazelle//:deps.bzl", "go_repository")

def go_repositories():
    go_repository(
        name = "cc_mvdan_gofumpt",
        build_file_proto_mode = "disable_global",
        importpath = "mvdan.cc/gofumpt",
        sum = "h1:bi/1aS/5W00E2ny5q65w9SnKpWEF/UIOqDYBILpo9rA=",
        version = "v0.1.1",
    )
    go_repository(
        name = "cc_mvdan_interfacer",
        build_file_proto_mode = "disable_global",
        importpath = "mvdan.cc/interfacer",
        sum = "h1:WX1yoOaKQfddO/mLzdV4wptyWgoH/6hwLs7QHTixo0I=",
        version = "v0.0.0-20180901003855-c20040233aed",
    )
    go_repository(
        name = "cc_mvdan_lint",
        build_file_proto_mode = "disable_global",
        importpath = "mvdan.cc/lint",
        sum = "h1:DxJ5nJdkhDlLok9K6qO+5290kphDJbHOQO1DFFFTeBo=",
        version = "v0.0.0-20170908181259-adc824a0674b",
    )
    go_repository(
        name = "cc_mvdan_unparam",
        build_file_proto_mode = "disable_global",
        importpath = "mvdan.cc/unparam",
        sum = "h1:HT3e4Krq+IE44tiN36RvVEb6tvqeIdtsVSsxmNPqlFU=",
        version = "v0.0.0-20210104141923-aac4ce9116a7",
    )

    go_repository(
        name = "co_honnef_go_tools",
        build_file_proto_mode = "disable_global",
        importpath = "honnef.co/go/tools",
        sum = "h1:/EPr//+UMMXwMTkXvCCoaJDq8cpjMO80Ou+L4PDo2mY=",
        version = "v0.2.1",
    )
    go_repository(
        name = "com_4d63_gochecknoglobals",
        build_file_proto_mode = "disable_global",
        importpath = "4d63.com/gochecknoglobals",
        sum = "h1:zeZSRqj5yCg28tCkIV/z/lWbwvNm5qnKVS15PI8nhD0=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_acomagu_bufpipe",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/acomagu/bufpipe",
        sum = "h1:fxAGrHZTgQ9w5QqVItgzwj235/uYZYgbXitB+dLupOk=",
        version = "v1.0.3",
    )

    go_repository(
        name = "com_github_adalogics_go_fuzz_headers",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/AdaLogics/go-fuzz-headers",
        sum = "h1:V8krnnfGj4pV65YLUm3C0/8bl7V5Nry2Pwvy3ru/wLc=",
        version = "v0.0.0-20210715213245-6c3934b029d8",
    )
    go_repository(
        name = "com_github_agnivade_levenshtein",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/agnivade/levenshtein",
        sum = "h1:QY8M92nrzkmr798gCo3kmMyqXFzdQVpxLlGPRBij0P8=",
        version = "v1.1.1",
    )

    go_repository(
        name = "com_github_alecthomas_template",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/alecthomas/template",
        sum = "h1:JYp7IbQjafoB+tBA3gMyHYHrpOtNuDiK/uB5uXxq5wM=",
        version = "v0.0.0-20190718012654-fb15b899a751",
    )
    go_repository(
        name = "com_github_alecthomas_units",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/alecthomas/units",
        sum = "h1:UQZhZ2O0vMHr2cI+DC1Mbh0TJxzA3RcLoMsFw+aXw7E=",
        version = "v0.0.0-20190924025748-f65c72e2690d",
    )
    go_repository(
        name = "com_github_alexbrainman_sspi",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/alexbrainman/sspi",
        sum = "h1:Kk6a4nehpJ3UuJRqlA3JxYxBZEqCeOmATOvrbT4p9RA=",
        version = "v0.0.0-20210105120005-909beea2cc74",
    )

    go_repository(
        name = "com_github_alexkohler_prealloc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/alexkohler/prealloc",
        sum = "h1:Hbq0/3fJPQhNkN0dR95AVrr6R7tou91y0uHG5pOcUuw=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_andreyvit_diff",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/andreyvit/diff",
        sum = "h1:bvNMNQO63//z+xNgfBlViaCIJKLlCJ6/fmUseuG0wVQ=",
        version = "v0.0.0-20170406064948-c7f18ee00883",
    )

    go_repository(
        name = "com_github_andybalholm_brotli",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/andybalholm/brotli",
        sum = "h1:V7DdXeJtZscaqfNuAdSRuRFzuiKlHSC/Zh3zl9qY3JY=",
        version = "v1.0.4",
    )

    go_repository(
        name = "com_github_antihax_optional",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/antihax/optional",
        sum = "h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_antlr_antlr4_runtime_go_antlr",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/antlr/antlr4/runtime/Go/antlr",
        sum = "h1:ue9pVfIcP+QMEjfgo/Ez4ZjNZfonGgR6NgjMaJMu1Cg=",
        version = "v0.0.0-20220418222510-f25a4f6275ed",
    )
    go_repository(
        name = "com_github_antonboom_errname",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Antonboom/errname",
        sum = "h1:IM+A/gz0pDhKmlt5KSNTVAvfLMb+65RxavBXpRtCUEg=",
        version = "v0.1.5",
    )
    go_repository(
        name = "com_github_antonboom_nilnil",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Antonboom/nilnil",
        sum = "h1:DLDavmg0a6G/F4Lt9t7Enrbgb3Oph6LnDE6YVsmTt74=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_aokoli_goutils",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aokoli/goutils",
        sum = "h1:7fpzNGoJ3VA8qcrm++XEE1QUe0mIwNeLa02Nwq7RDkg=",
        version = "v1.0.1",
    )

    go_repository(
        name = "com_github_aquasecurity_starboard",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aquasecurity/starboard",
        sum = "h1:VzwNZqE1Z6sI4JGJ2Nm2E11g4PCR3JM6ap573OWmjZY=",
        version = "v0.15.8",
    )
    go_repository(
        name = "com_github_arbovm_levenshtein",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/arbovm/levenshtein",
        sum = "h1:jfIu9sQUG6Ig+0+Ap1h4unLjW6YQJpKZVmUzxsD4E/Q=",
        version = "v0.0.0-20160628152529-48b4e1c0c4d0",
    )

    go_repository(
        name = "com_github_armon_circbuf",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/armon/circbuf",
        sum = "h1:QEF07wC0T1rKkctt1RINW/+RMTVmiwxETico2l3gxJA=",
        version = "v0.0.0-20150827004946-bbbad097214e",
    )
    go_repository(
        name = "com_github_armon_consul_api",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/armon/consul-api",
        sum = "h1:G1bPvciwNyF7IUmKXNt9Ak3m6u9DE1rF+RmtIkBpVdA=",
        version = "v0.0.0-20180202201655-eb2c6b5be1b6",
    )

    go_repository(
        name = "com_github_armon_go_metrics",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/armon/go-metrics",
        sum = "h1:8GUt8eRujhVEGZFFEjBj46YV4rDjvGrNxb0KMWYkL2I=",
        version = "v0.0.0-20180917152333-f0300d1749da",
    )
    go_repository(
        name = "com_github_armon_go_radix",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/armon/go-radix",
        sum = "h1:F4z6KzEeeQIMeLFa97iZU6vupzoecKdU5TX24SNppXI=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_armon_go_socks5",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/armon/go-socks5",
        sum = "h1:0CwZNZbxp69SHPdPJAN/hZIm0C4OItdklCFmMRWYpio=",
        version = "v0.0.0-20160902184237-e75332964ef5",
    )

    go_repository(
        name = "com_github_asaskevich_govalidator",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/asaskevich/govalidator",
        sum = "h1:4daAzAu0S6Vi7/lbWECcX0j45yZReDZ56BQsrVBOEEY=",
        version = "v0.0.0-20200428143746-21a406dcc535",
    )
    go_repository(
        name = "com_github_ash2k_stager",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ash2k/stager",
        sum = "h1:3gpvrLGQzRmUCk1+A1EmrcSBu/JyeoU09fnITBl4Z/k=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_ashanbrown_forbidigo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ashanbrown/forbidigo",
        sum = "h1:RMlEFupPCxQ1IogYOQUnIQwGEUGK8g5vAPMRyJoSxbc=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_ashanbrown_makezero",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ashanbrown/makezero",
        sum = "h1:YOsoVXsZQPA9aOTy1g0lAJv5VzZUvwQuZqug8XPeqfM=",
        version = "v0.0.0-20210520155254-b6261585ddde",
    )

    go_repository(
        name = "com_github_avast_retry_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/avast/retry-go",
        sum = "h1:4SOWQ7Qs+oroOTQOYnAHqelpCO0biHSxpiH9JdtuBj0=",
        version = "v3.0.0+incompatible",
    )

    go_repository(
        name = "com_github_aws_aws_sdk_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go",
        sum = "h1:yJZIr8nMV1hXjAvvOLUFqZRJcHV7udPQBfhJqawDzI0=",
        version = "v1.43.31",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2",
        sum = "h1:fqlCk6Iy3bnCumtrLz9r3mJ/2gUT0pJ0wLFVIdWh+JA=",
        version = "v1.16.2",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_aws_protocol_eventstream",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/aws/protocol/eventstream",
        sum = "h1:SdK4Ppk5IzLs64ZMvr6MrSficMtjY2oS0WOORXTlxwU=",
        version = "v1.4.1",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_config",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/config",
        sum = "h1:5AlQD0jhVXlGzwo+VORKiUuogkG7pQcLJNzIzK7eodw=",
        version = "v1.15.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_credentials",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/credentials",
        sum = "h1:RQQ5fzclAKJyY5TvF+fkjJEwzK4hnxQCLOu5JXzDmQo=",
        version = "v1.11.2",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_feature_ec2_imds",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/feature/ec2/imds",
        sum = "h1:LWPg5zjHV9oz/myQr4wMs0gi4CjnDN/ILmyZUFYXZsU=",
        version = "v1.12.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_feature_s3_manager",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/feature/s3/manager",
        sum = "h1:ir7iEq78s4txFGgwcLqD6q9IIPzTQNRJXulJd9h/zQo=",
        version = "v1.11.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_configsources",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/configsources",
        sum = "h1:onz/VaaxZ7Z4V+WIN9Txly9XLTmoOh1oJ8XcAC3pako=",
        version = "v1.1.9",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_endpoints_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/endpoints/v2",
        sum = "h1:9stUQR/u2KXU6HkFJYlqnZEjBnbgrVbG6I5HN09xZh0=",
        version = "v2.4.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_internal_ini",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/internal/ini",
        sum = "h1:by9P+oy3P/CwggN4ClnW2D4oL91QV7pBzBICi1chZvQ=",
        version = "v1.3.10",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_accept_encoding",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/accept-encoding",
        sum = "h1:T4pFel53bkHjL2mMo+4DKE6r6AuoZnM0fg7k1/ratr4=",
        version = "v1.9.1",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_checksum",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/checksum",
        sum = "h1:I0dcwWitE752hVSMrsLCxqNQ+UdEp3nACx2bYNMQq+k=",
        version = "v1.1.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_presigned_url",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/presigned-url",
        sum = "h1:Gh1Gpyh01Yvn7ilO/b/hr01WgNpaszfbKMUgqM186xQ=",
        version = "v1.9.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_internal_s3shared",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/service/internal/s3shared",
        sum = "h1:BKjwCJPnANbkwQ8vzSbaZDKawwagDubrH/z/c0X+kbQ=",
        version = "v1.13.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_s3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/service/s3",
        sum = "h1:rMPtwA7zzkSQZhhz9U3/SoIDz/NZ7Q+iRn4EIO8rSyU=",
        version = "v1.26.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_sso",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/service/sso",
        sum = "h1:frW4ikGcxfAEDfmQqWgMLp+F1n4nRo9sF39OcIb5BkQ=",
        version = "v1.11.3",
    )
    go_repository(
        name = "com_github_aws_aws_sdk_go_v2_service_sts",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/aws-sdk-go-v2/service/sts",
        sum = "h1:cJGRyzCSVwZC7zZZ1xbx9m32UnrKydRYhOvcD1NYP9Q=",
        version = "v1.16.3",
    )
    go_repository(
        name = "com_github_aws_smithy_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aws/smithy-go",
        sum = "h1:eG/N+CcUMAvsdffgMvjMKwfyDzIkjM6pfxMJ8Mzc6mE=",
        version = "v1.11.2",
    )
    go_repository(
        name = "com_github_aymerick_douceur",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/aymerick/douceur",
        sum = "h1:Mv+mAeH1Q+n9Fr+oyamOlAkUNPWPlA8PPGR0QAaYuPk=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_github_azure_azure_pipeline_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/azure-pipeline-go",
        sum = "h1:7U9HBg1JFK3jHl5qmo4CTZKFTVgMwdFHMVtCdfBE21U=",
        version = "v0.2.3",
    )

    go_repository(
        name = "com_github_azure_azure_storage_blob_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/azure-storage-blob-go",
        sum = "h1:1BCg74AmVdYwO3dlKwtFU1V0wU2PZdREkXvAmZJRUlM=",
        version = "v0.14.0",
    )

    go_repository(
        name = "com_github_azure_go_ansiterm",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/go-ansiterm",
        sum = "h1:UQHMgLO+TxOElx5B5HZ4hJQsoJ/PvUvKRhJHDQXO8P8=",
        version = "v0.0.0-20210617225240-d185dfc1b5a1",
    )
    go_repository(
        name = "com_github_azure_go_autorest",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/go-autorest",
        sum = "h1:V5VMDjClD3GiElqLWO7mz2MxNAK/vTfRHdAubSIPRgs=",
        version = "v14.2.0+incompatible",
    )
    go_repository(
        name = "com_github_azure_go_autorest_autorest",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/go-autorest/autorest",
        sum = "h1:F3R3q42aWytozkV8ihzcgMO4OA4cuqr3bNlsEuF6//A=",
        version = "v0.11.27",
    )
    go_repository(
        name = "com_github_azure_go_autorest_autorest_adal",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/go-autorest/autorest/adal",
        sum = "h1:gJ3E98kMpFB1MFqQCvA1yFab8vthOeD4VlFRQULxahg=",
        version = "v0.9.20",
    )

    go_repository(
        name = "com_github_azure_go_autorest_autorest_date",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/go-autorest/autorest/date",
        sum = "h1:7gUk1U5M/CQbp9WoqinNzJar+8KY+LPI6wiWrP/myHw=",
        version = "v0.3.0",
    )

    go_repository(
        name = "com_github_azure_go_autorest_logger",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/go-autorest/logger",
        sum = "h1:IG7i4p/mDa2Ce4TRyAO8IHnVhAVF3RFU+ZtXWSmf4Tg=",
        version = "v0.2.1",
    )
    go_repository(
        name = "com_github_azure_go_autorest_tracing",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Azure/go-autorest/tracing",
        sum = "h1:TYi4+3m5t6K48TGI9AUdb+IzbnSxvnvUMfuitfgcfuo=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_github_beevik_ntp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/beevik/ntp",
        sum = "h1:xzVrPrE4ziasFXgBVBZJDP0Wg/KpMwk2KHJ4Ba8GrDw=",
        version = "v0.3.0",
    )

    go_repository(
        name = "com_github_benbjohnson_clock",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/benbjohnson/clock",
        sum = "h1:ip6w0uFQkncKQ979AypyG0ER7mqUSBdKLOgAle/AT8A=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_beorn7_perks",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/beorn7/perks",
        sum = "h1:VlbKKnNfV8bJzeqoa4cOKqO6bYr3WgKZxO8Z16+hsOM=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_bgentry_speakeasy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bgentry/speakeasy",
        sum = "h1:ByYyxL9InA1OWqxJqqp2A5pYHUrCiAL6K3J+LKSsQkY=",
        version = "v0.1.0",
    )

    go_repository(
        name = "com_github_bketelsen_crypt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bketelsen/crypt",
        sum = "h1:w/jqZtC9YD4DS/Vp9GhWfWcCpuAL58oTnLoI8vE9YHU=",
        version = "v0.0.4",
    )
    go_repository(
        name = "com_github_bkielbasa_cyclop",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bkielbasa/cyclop",
        sum = "h1:7Jmnh0yL2DjKfw28p86YTd/B4lRGcNuu12sKE35sM7A=",
        version = "v1.2.0",
    )

    go_repository(
        name = "com_github_blang_semver",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/blang/semver",
        sum = "h1:cQNTCjp13qL8KC3Nbxr/y2Bqb63oX6wdnnjpJbkM4JQ=",
        version = "v3.5.1+incompatible",
    )
    go_repository(
        name = "com_github_blang_semver_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/blang/semver/v4",
        sum = "h1:1PFHFE6yCCTv8C1TeyNNarDzntLi7wMI5i/pzqYIsAM=",
        version = "v4.0.0",
    )
    go_repository(
        name = "com_github_blizzy78_varnamelen",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/blizzy78/varnamelen",
        sum = "h1:80mYO7Y5ppeEefg1Jzu+NBg16iwToOQVnDnNIoWSShs=",
        version = "v0.3.0",
    )

    go_repository(
        name = "com_github_bmatcuk_doublestar_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bmatcuk/doublestar/v2",
        sum = "h1:6I6oUiT/sU27eE2OFcWqBhL1SwjyvQuOssxT4a1yidI=",
        version = "v2.0.4",
    )

    go_repository(
        name = "com_github_bombsimon_wsl_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bombsimon/wsl/v3",
        sum = "h1:Mka/+kRLoQJq7g2rggtgQsjuI/K5Efd87WX96EWFxjM=",
        version = "v3.3.0",
    )
    go_repository(
        name = "com_github_breml_bidichk",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/breml/bidichk",
        sum = "h1:Qpy8Rmgos9qdJxhka0K7ADEE5bQZX9PQUthkgggHpFM=",
        version = "v0.1.1",
    )

    go_repository(
        name = "com_github_bshuster_repo_logrus_logstash_hook",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bshuster-repo/logrus-logstash-hook",
        sum = "h1:e+C0SB5R1pu//O4MQ3f9cFuPGoOVeF2fE4Og9otCc70=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_buger_jsonparser",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/buger/jsonparser",
        sum = "h1:2PnMjfWD7wBILjqQbt530v576A/cAbQvEW9gGIpYMUs=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_bugsnag_bugsnag_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bugsnag/bugsnag-go",
        sum = "h1:rFt+Y/IK1aEZkEHchZRSq9OQbsSzIT/OrI8YFFmRIng=",
        version = "v0.0.0-20141110184014-b1d153021fcd",
    )
    go_repository(
        name = "com_github_bugsnag_osext",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bugsnag/osext",
        sum = "h1:otBG+dV+YK+Soembjv71DPz3uX/V/6MMlSyD9JBQ6kQ=",
        version = "v0.0.0-20130617224835-0dd3f918b21b",
    )
    go_repository(
        name = "com_github_bugsnag_panicwrap",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bugsnag/panicwrap",
        sum = "h1:nvj0OLI3YqYXer/kZD8Ri1aaunCxIEsOst1BVJswV0o=",
        version = "v0.0.0-20151223152923-e2c28503fcd0",
    )
    go_repository(
        name = "com_github_burntsushi_toml",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/BurntSushi/toml",
        sum = "h1:Rt8g24XnyGTyglgET/PRUNlrUeu9F5L+7FilkXfZgs0=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_burntsushi_xgb",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/BurntSushi/xgb",
        sum = "h1:1BDTz0u9nC3//pOCMdNH+CiXJVYJh5UQNCOBG7jbELc=",
        version = "v0.0.0-20160522181843-27f122750802",
    )
    go_repository(
        name = "com_github_butuzov_ireturn",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/butuzov/ireturn",
        sum = "h1:QvrO2QF2+/Cx1WA/vETCIYBKtRjc30vesdoPUNo1EbY=",
        version = "v0.1.1",
    )

    go_repository(
        name = "com_github_bytecodealliance_wasmtime_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/bytecodealliance/wasmtime-go",
        sum = "h1:B6thr7RMM9xQmouBtUqm1RpkJjuLS37m6nxX+iwsQSc=",
        version = "v0.36.0",
    )

    go_repository(
        name = "com_github_caarlos0_env_v6",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/caarlos0/env/v6",
        sum = "h1:lA7sxiGArZ2KkiqpOQNf8ERBRWI+v8MWIH+eGjSN22I=",
        version = "v6.10.0",
    )

    go_repository(
        name = "com_github_cenkalti_backoff_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cenkalti/backoff/v4",
        sum = "h1:cFAlzYUlVYDysBEH2T5hyJZMh3+5+WCBvSnK6Q8UtC4=",
        version = "v4.1.3",
    )

    go_repository(
        name = "com_github_census_instrumentation_opencensus_proto",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/census-instrumentation/opencensus-proto",
        sum = "h1:t/LhUZLVitR1Ow2YOnduCsavhwFUklBMoGVYUCqmCqk=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_certifi_gocertifi",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/certifi/gocertifi",
        sum = "h1:S2NE3iHSwP0XV47EEXL8mWmRdEfGscSJ+7EgePNgt0s=",
        version = "v0.0.0-20210507211836-431795d63e8d",
    )
    go_repository(
        name = "com_github_cespare_xxhash",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cespare/xxhash",
        sum = "h1:a6HrQnmkObjyL+Gs60czilIUGqrzKutQD6XZog3p+ko=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_cespare_xxhash_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cespare/xxhash/v2",
        sum = "h1:YRXhKfTDauu4ajMg1TPgFO5jnlC2HCbmLXMcTG5cbYE=",
        version = "v2.1.2",
    )
    go_repository(
        name = "com_github_chai2010_gettext_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/chai2010/gettext-go",
        sum = "h1:1Lwwip6Q2QGsAdl/ZKPCwTe9fe0CjlUbqj5bFNSjIRk=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_charithe_durationcheck",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/charithe/durationcheck",
        sum = "h1:mPP4ucLrf/rKZiIG/a9IPXHGlh8p4CzgpyTy6EEutYk=",
        version = "v0.0.9",
    )
    go_repository(
        name = "com_github_chavacava_garif",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/chavacava/garif",
        sum = "h1:spmv8nSH9h5oCQf40jt/ufBCt9j0/58u4G+rkeMqXGI=",
        version = "v0.0.0-20210405164556-e8a0a408d6af",
    )

    go_repository(
        name = "com_github_chzyer_logex",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/chzyer/logex",
        sum = "h1:Swpa1K6QvQznwJRcfTfQJmTE72DqScAa40E+fbHEXEE=",
        version = "v1.1.10",
    )
    go_repository(
        name = "com_github_chzyer_readline",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/chzyer/readline",
        sum = "h1:fY5BOSpyZCqRo5OhCuC+XN+r/bBCmeuuJtjz+bCNIf8=",
        version = "v0.0.0-20180603132655-2972be24d48e",
    )
    go_repository(
        name = "com_github_chzyer_test",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/chzyer/test",
        sum = "h1:q763qf9huN11kDQavWsoZXJNW3xEE4JJyHa5Q25/sd8=",
        version = "v0.0.0-20180213035817-a1ea475d72b1",
    )

    go_repository(
        name = "com_github_cilium_ebpf",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cilium/ebpf",
        sum = "h1:1k/q3ATgxSXRdrmPfH8d7YK0GfqVsEKZAX9dQZvs56k=",
        version = "v0.7.0",
    )

    go_repository(
        name = "com_github_client9_misspell",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/client9/misspell",
        sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
        version = "v0.3.4",
    )
    go_repository(
        name = "com_github_client9_reopen",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/client9/reopen",
        sum = "h1:8tpLVR74DLpLObrn2KvsyxJY++2iORGR17WLUdSzUws=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_cloudflare_circl",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cloudflare/circl",
        sum = "h1:bZgT/A+cikZnKIwn7xL2OBj012Bmvho/o6RpRvv3GKY=",
        version = "v1.1.0",
    )

    go_repository(
        name = "com_github_cloudflare_tableflip",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cloudflare/tableflip",
        sum = "h1:8I+B99QnnEWPHOY3fWipwVKxS70LGgUsslG7CSfmHMw=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_cloudykit_fastprinter",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/CloudyKit/fastprinter",
        sum = "h1:sR+/8Yb4slttB4vD+b9btVEnWgL3Q00OBTzVT8B9C0c=",
        version = "v0.0.0-20200109182630-33d98a066a53",
    )
    go_repository(
        name = "com_github_cloudykit_jet_v6",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/CloudyKit/jet/v6",
        sum = "h1:hvO96X345XagdH1fAoBjpBYG4a1ghhL/QzalkduPuXk=",
        version = "v6.1.0",
    )

    go_repository(
        name = "com_github_cncf_udpa_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cncf/udpa/go",
        sum = "h1:hzAQntlaYRkVSFEfj9OTWlVV1H155FMD8BTKktLv0QI=",
        version = "v0.0.0-20210930031921-04548b0d99d4",
    )
    go_repository(
        name = "com_github_cncf_xds_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cncf/xds/go",
        sum = "h1:zH8ljVhhq7yC0MIeUL/IviMtY8hx2mK8cN9wEYb8ggw=",
        version = "v0.0.0-20211011173535-cb28da3451f1",
    )

    go_repository(
        name = "com_github_cockroachdb_datadriven",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cockroachdb/datadriven",
        sum = "h1:OaNxuTZr7kxeODyLWsRMC+OD03aFUH+mW6r2d+MWa5Y=",
        version = "v0.0.0-20190809214429-80d97fb3cbaa",
    )

    go_repository(
        name = "com_github_codegangsta_inject",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/codegangsta/inject",
        sum = "h1:sDMmm+q/3+BukdIpxwO365v/Rbspp2Nt5XntgQRXq8Q=",
        version = "v0.0.0-20150114235600-33e0aa1cb7c0",
    )
    go_repository(
        name = "com_github_containerd_aufs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/aufs",
        sum = "h1:2oeJiwX5HstO7shSrPZjrohJZLzK36wvpdmzDRkL/LY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_containerd_btrfs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/btrfs",
        sum = "h1:osn1exbzdub9L5SouXO5swW4ea/xVdJZ3wokxN5GrnA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_containerd_cgroups",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/cgroups",
        sum = "h1:jN/mbWBEaz+T1pi5OFtnkQ+8qnmEbAr1Oo1FRm5B0dA=",
        version = "v1.0.4",
    )
    go_repository(
        name = "com_github_containerd_console",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/console",
        sum = "h1:lIr7SlA5PxZyMV30bDW0MGbiOPXwc63yRuCP0ARubLw=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_containerd_containerd",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/containerd",
        sum = "h1:h4dOFDwzHmqFEP754PgfgTeVXFnLiRc6kiqC7tplDJs=",
        version = "v1.6.8",
    )
    go_repository(
        name = "com_github_containerd_continuity",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/continuity",
        sum = "h1:QSqfxcn8c+12slxwu00AtzXrsami0MJb/MQs9lOLHLA=",
        version = "v0.2.2",
    )
    go_repository(
        name = "com_github_containerd_fifo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/fifo",
        sum = "h1:6PirWBr9/L7GDamKr+XM0IeUFXu5mf3M/BPpH9gaLBU=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_containerd_go_cni",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/go-cni",
        sum = "h1:el5WPymG5nRRLQF1EfB97FWob4Tdc8INg8RZMaXWZlo=",
        version = "v1.1.6",
    )
    go_repository(
        name = "com_github_containerd_go_runc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/go-runc",
        sum = "h1:oU+lLv1ULm5taqgV/CJivypVODI4SUz1znWjv3nNYS0=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_containerd_imgcrypt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/imgcrypt",
        sum = "h1:iKTstFebwy3Ak5UF0RHSeuCTahC5OIrPJa6vjMAM81s=",
        version = "v1.1.4",
    )
    go_repository(
        name = "com_github_containerd_nri",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/nri",
        sum = "h1:6QioHRlThlKh2RkRTR4kIT3PKAcrLo3gIWnjkM4dQmQ=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_containerd_stargz_snapshotter_estargz",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/stargz-snapshotter/estargz",
        sum = "h1:idtwRTLjk2erqiYhPWy2L844By8NRFYEwYHcXhoIWPM=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_github_containerd_ttrpc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/ttrpc",
        sum = "h1:GbtyLRxb0gOLR0TYQWt3O6B0NvT8tMdorEHqIQo/lWI=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_containerd_typeurl",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/typeurl",
        sum = "h1:Chlt8zIieDbzQFzXzAeBEF92KhExuE4p9p92/QmY7aY=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_containerd_zfs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containerd/zfs",
        sum = "h1:cXLJbx+4Jj7rNsTiqVfm6i+RNLx6FFA2fMmDlEf+Wm8=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_containernetworking_cni",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containernetworking/cni",
        sum = "h1:ky20T7c0MvKvbMOwS/FrlbNwjEoqJEUUYfsL4b0mc4k=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_containernetworking_plugins",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containernetworking/plugins",
        sum = "h1:+AGfFigZ5TiQH00vhR8qPeSatj53eNGz0C1d3wVYlHE=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_containers_ocicrypt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/containers/ocicrypt",
        sum = "h1:uMxn2wTb4nDR7GqG3rnZSfpJXqWURfzZ7nKydzIeKpA=",
        version = "v1.1.3",
    )

    go_repository(
        name = "com_github_coreos_etcd",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/coreos/etcd",
        sum = "h1:jFneRYjIvLMLhDLCzuTuU4rSJUjRplcJQ7pD7MnhC04=",
        version = "v3.3.10+incompatible",
    )
    go_repository(
        name = "com_github_coreos_go_etcd",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/coreos/go-etcd",
        sum = "h1:bXhRBIXoTm9BYHS3gE0TtQuyNZyeEMux2sDi4oo5YOo=",
        version = "v2.0.0+incompatible",
    )

    go_repository(
        name = "com_github_coreos_go_oidc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/coreos/go-oidc",
        sum = "h1:sdJrfw8akMnCuUlaZU3tE/uYXFgfqom8DBE9so9EBsM=",
        version = "v2.1.0+incompatible",
    )
    go_repository(
        name = "com_github_coreos_go_semver",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/coreos/go-semver",
        sum = "h1:wkHLiw0WNATZnSG7epLsujiMCgPAc9xhjJ4tgnAxmfM=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_coreos_go_systemd",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/coreos/go-systemd",
        sum = "h1:W8b4lQ4tFF21aspRGoBuCNV6V2fFJBF+pm1J6OY8Lys=",
        version = "v0.0.0-20190620071333-e64a0ec8b42a",
    )
    go_repository(
        name = "com_github_coreos_go_systemd_v22",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/coreos/go-systemd/v22",
        sum = "h1:D9/bQk5vlXQFZ6Kwuu6zaiXJ9oTPe68++AzAJc1DzSI=",
        version = "v22.3.2",
    )
    go_repository(
        name = "com_github_coreos_pkg",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/coreos/pkg",
        sum = "h1:lBNOc5arjvs8E5mO2tbpBpLoyyu8B6e44T7hJy6potg=",
        version = "v0.0.0-20180928190104-399ea9e2e55f",
    )
    go_repository(
        name = "com_github_cpuguy83_go_md2man",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cpuguy83/go-md2man",
        sum = "h1:BSKMNlYxDvnunlTymqtgONjNnaRV1sTpcovwwjF22jk=",
        version = "v1.0.10",
    )

    go_repository(
        name = "com_github_cpuguy83_go_md2man_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cpuguy83/go-md2man/v2",
        sum = "h1:p1EgwI/C7NhT0JmVkwCD2ZBK8j4aeHQX2pMHHBfMQ6w=",
        version = "v2.0.2",
    )
    go_repository(
        name = "com_github_creack_pty",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/creack/pty",
        sum = "h1:07n33Z8lZxZ2qwegKbObQohDhXDQxiMMz1NOUGYlesw=",
        version = "v1.1.11",
    )
    go_repository(
        name = "com_github_cyphar_filepath_securejoin",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/cyphar/filepath-securejoin",
        sum = "h1:YX6ebbZCZP7VkM3scTTokDgBL2TY741X51MTk3ycuNI=",
        version = "v0.2.3",
    )

    go_repository(
        name = "com_github_daixiang0_gci",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/daixiang0/gci",
        sum = "h1:iwJvwQpBZmMg31w+QQ6jsyZ54KEATn6/nfARbBNW294=",
        version = "v0.2.9",
    )

    go_repository(
        name = "com_github_danieljoos_wincred",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/danieljoos/wincred",
        sum = "h1:3RNcEpBg4IhIChZdFRSdlQt1QjCp1sMAPIrOnm7Yf8g=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_data_dog_go_sqlmock",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/DATA-DOG/go-sqlmock",
        sum = "h1:Shsta01QNfFxHCfpW6YH2STWB0MudeXXEWMr20OEh60=",
        version = "v1.5.0",
    )

    go_repository(
        name = "com_github_datadog_datadog_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/DataDog/datadog-go",
        sum = "h1:R7WqXWP4fIOAqWJtUKmSfuc7eDsBT58k9AY5WSHVosk=",
        version = "v4.4.0+incompatible",
    )

    go_repository(
        name = "com_github_datadog_sketches_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/DataDog/sketches-go",
        sum = "h1:chm5KSXO7kO+ywGWJ0Zs6tdmWU8PBXSbywFVciL6BG4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_davecgh_go_spew",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/davecgh/go-spew",
        sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_daviddengcn_go_colortext",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/daviddengcn/go-colortext",
        sum = "h1:ANqDyC0ys6qCSvuEK7l3g5RaehL/Xck9EX8ATG8oKsE=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_denis_tingajkin_go_header",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/denis-tingajkin/go-header",
        sum = "h1:jEeSF4sdv8/3cT/WY8AgDHUoItNSoEZ7qg9dX7pc218=",
        version = "v0.4.2",
    )
    go_repository(
        name = "com_github_denisenkom_go_mssqldb",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/denisenkom/go-mssqldb",
        sum = "h1:RSohk2RsiZqLZ0zCjtfn3S4Gp4exhpBWHyQ7D0yGjAk=",
        version = "v0.9.0",
    )

    go_repository(
        name = "com_github_dgraph_io_badger_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dgraph-io/badger/v3",
        sum = "h1:dpyM5eCJAtQCBcMCZcT4UBZchuTJgCywerHHgmxfxM8=",
        version = "v3.2103.2",
    )
    go_repository(
        name = "com_github_dgraph_io_ristretto",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dgraph-io/ristretto",
        sum = "h1:Jv3CGQHp9OjuMBSne1485aDpUkTKEcUqF+jm/LuerPI=",
        version = "v0.1.0",
    )

    go_repository(
        name = "com_github_dgrijalva_jwt_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dgrijalva/jwt-go",
        sum = "h1:7qlOGliEKZXTDg6OTjfoBKDXWrumCAMpl/TFQ4/5kLM=",
        version = "v3.2.0+incompatible",
    )
    go_repository(
        name = "com_github_dgryski_go_farm",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dgryski/go-farm",
        sum = "h1:fAjc9m62+UWV/WAFKLNi6ZS0675eEUC9y3AlwSbQu1Y=",
        version = "v0.0.0-20200201041132-a6ae2369ad13",
    )

    go_repository(
        name = "com_github_dgryski_go_minhash",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dgryski/go-minhash",
        sum = "h1:54Y/2GF52MSJ4n63HWvNDFRtztgm6tq2UrOX61sjGKc=",
        version = "v0.0.0-20170608043002-7fe510aff544",
    )
    go_repository(
        name = "com_github_dgryski_go_rendezvous",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dgryski/go-rendezvous",
        sum = "h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=",
        version = "v0.0.0-20200823014737-9f7001d12a5f",
    )

    go_repository(
        name = "com_github_dgryski_trifles",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dgryski/trifles",
        sum = "h1:fRzb/w+pyskVMQ+UbP35JkH8yB7MYb4q/qhBarqZE6g=",
        version = "v0.0.0-20200323201526-dd97f9abfb48",
    )

    go_repository(
        name = "com_github_distribution_distribution_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/distribution/distribution/v3",
        sum = "h1:hbCT8ZPPMqefiAWD2ZKjn7ypokIGViTvBBg/ExLSdCk=",
        version = "v3.0.0-20220526142353-ffbd94cbe269",
    )
    go_repository(
        name = "com_github_djarvur_go_err113",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Djarvur/go-err113",
        sum = "h1:sHglBQTwgx+rWPdisA5ynNEsoARbiCBOyGcJM4/OzsM=",
        version = "v0.0.0-20210108212216-aea10b59be24",
    )

    go_repository(
        name = "com_github_docker_cli",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/cli",
        sum = "h1:eO2KS7ZFeov5UJeaDmIs1NFEDRf32PaqRpvoEkKBy5M=",
        version = "v20.10.17+incompatible",
    )
    go_repository(
        name = "com_github_docker_distribution",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/distribution",
        sum = "h1:Q50tZOPR6T/hjNsyc9g8/syEs6bk8XXApsHjKukMl68=",
        version = "v2.8.1+incompatible",
    )
    go_repository(
        name = "com_github_docker_docker",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/docker",
        sum = "h1:JYCuMrWaVNophQTOrMMoSwudOVEfcegoZZrleKc1xwE=",
        version = "v20.10.17+incompatible",
    )
    go_repository(
        name = "com_github_docker_docker_credential_helpers",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/docker-credential-helpers",
        sum = "h1:axCks+yV+2MR3/kZhAmy07yC56WZ2Pwu/fKWtKuZB0o=",
        version = "v0.6.4",
    )
    go_repository(
        name = "com_github_docker_go_connections",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/go-connections",
        sum = "h1:El9xVISelRB7BuFusrZozjnkIM5YnzCViNKohAFqRJQ=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_docker_go_events",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/go-events",
        sum = "h1:+pKlWGMw7gf6bQ+oDZB4KHQFypsfjYlq/C4rfL7D3g8=",
        version = "v0.0.0-20190806004212-e31b211e4f1c",
    )
    go_repository(
        name = "com_github_docker_go_metrics",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/go-metrics",
        sum = "h1:AgB/0SvBxihN0X8OR4SjsblXkbMvalQ8cjmtKQ2rQV8=",
        version = "v0.0.1",
    )
    go_repository(
        name = "com_github_docker_go_units",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/go-units",
        sum = "h1:3uh0PgVws3nIA0Q+MwDC8yjEPf9zjRfZZWXZYDct3Tw=",
        version = "v0.4.0",
    )

    go_repository(
        name = "com_github_docker_libtrust",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docker/libtrust",
        sum = "h1:ZClxb8laGDf5arXfYcAtECDFgAgHklGI8CxgjHnXKJ4=",
        version = "v0.0.0-20150114040149-fa567046d9b1",
    )

    go_repository(
        name = "com_github_docopt_docopt_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/docopt/docopt-go",
        sum = "h1:bWDMxwH3px2JBh6AyO7hdCn/PkvCZXii8TGj7sbtEbQ=",
        version = "v0.0.0-20180111231733-ee0de3bc6815",
    )
    go_repository(
        name = "com_github_dpotapov_go_spnego",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dpotapov/go-spnego",
        sum = "h1:3EIvol92cLWG5m13Me3/LfEtQqr23xzwZCxiWoT5gGE=",
        version = "v0.0.0-20220426193508-b7f82e4507db",
    )
    go_repository(
        name = "com_github_dustin_go_humanize",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/dustin/go-humanize",
        sum = "h1:VSnTsYCnlFHaM2/igO1h6X3HA71jcobQuxemgkq4zYo=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_eknkc_amber",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/eknkc/amber",
        sum = "h1:clC1lXBpe2kTj2VHdaIu9ajZQe4kcEY9j0NsnDDBZ3o=",
        version = "v0.0.0-20171010120322-cdade1c07385",
    )
    go_repository(
        name = "com_github_ekzhu_minhash_lsh",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ekzhu/minhash-lsh",
        sum = "h1:/7G7q8SDJdrah5jDYqZI8pGFjSqiCzfSEO+NgqKCYX0=",
        version = "v0.0.0-20171225071031-5c06ee8586a1",
    )
    go_repository(
        name = "com_github_elazarl_goproxy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/elazarl/goproxy",
        sum = "h1:yUdfgN0XgIJw7foRItutHYUIhlcKzcSf5vDpdhQAKTc=",
        version = "v0.0.0-20180725130230-947c36da3153",
    )
    go_repository(
        name = "com_github_emicklei_go_restful",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/emicklei/go-restful",
        sum = "h1:spTtZBk5DYEvbxMVutUuTyh1Ao2r4iyvLdACqsl/Ljk=",
        version = "v2.9.5+incompatible",
    )
    go_repository(
        name = "com_github_emicklei_go_restful_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/emicklei/go-restful/v3",
        sum = "h1:eCZ8ulSerjdAiaNpF7GxXIE7ZCMo1moN1qX+S609eVw=",
        version = "v3.8.0",
    )

    go_repository(
        name = "com_github_emirpasic_gods",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/emirpasic/gods",
        sum = "h1:FXtiHYKDGKCW2KzwZKx0iC0PQmdlorYgdFG9jPXJ1Bc=",
        version = "v1.18.1",
    )
    go_repository(
        name = "com_github_envoyproxy_go_control_plane",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/envoyproxy/go-control-plane",
        sum = "h1:xvqufLtNVwAhN8NMyWklVgxnWohi+wtMGQMhtxexlm0=",
        version = "v0.10.2-0.20220325020618-49ff273808a1",
    )
    go_repository(
        name = "com_github_esimonov_ifshort",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/esimonov/ifshort",
        sum = "h1:JD6x035opqGec5fZ0TLjXeROD2p5H7oLGn8MKfy9HTM=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_ettle_strcase",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ettle/strcase",
        sum = "h1:htFueZyVeE1XNnMEfbqp5r67qAN/4r6ya1ysq8Q+Zcw=",
        version = "v0.1.1",
    )

    go_repository(
        name = "com_github_evanphx_json_patch",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/evanphx/json-patch",
        sum = "h1:jBYDEEiFBPxA0v50tFdvOzQQTCvpL6mnFh5mB2/l16U=",
        version = "v5.6.0+incompatible",
    )
    go_repository(
        name = "com_github_evanphx_json_patch_v5",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/evanphx/json-patch/v5",
        sum = "h1:b91NhWfaz02IuVxO9faSllyAtNXHMPkC5J8sJCLunww=",
        version = "v5.6.0",
    )

    go_repository(
        name = "com_github_exponent_io_jsonpath",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/exponent-io/jsonpath",
        sum = "h1:105gxyaGwCFad8crR9dcMQWvV9Hvulu6hwUh4tWPJnM=",
        version = "v0.0.0-20151013193312-d6023ce2651d",
    )

    go_repository(
        name = "com_github_fatih_camelcase",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fatih/camelcase",
        sum = "h1:hxNvNX/xYBp0ovncs8WyWZrOrpBNub/JfaMvbURyft8=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_fatih_color",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fatih/color",
        sum = "h1:8LOYc1KYPPmyKMuN8QV2DNRWNbLo6LZ0iLs8+mlH53w=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_github_fatih_structs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fatih/structs",
        sum = "h1:Q7juDM0QtcnhCpeyLGQKyg4TOIghuNXrkL32pHAUMxo=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_fatih_structtag",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fatih/structtag",
        sum = "h1:/OdNE99OxoI/PqaW/SuSK9uxxT3f/tcSZgon/ssNSx4=",
        version = "v1.2.0",
    )

    go_repository(
        name = "com_github_felixge_httpsnoop",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/felixge/httpsnoop",
        sum = "h1:s/nj+GCswXYzN5v2DpNMuMQYe+0DDwt5WVCU6CWBdXk=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_flosch_pongo2_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/flosch/pongo2/v4",
        sum = "h1:gv+5Pe3vaSVmiJvh/BZa82b7/00YUGm0PIyVVLop0Hw=",
        version = "v4.0.2",
    )
    go_repository(
        name = "com_github_flowstack_go_jsonschema",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/flowstack/go-jsonschema",
        sum = "h1:dCrjGJRXIlbDsLAgTJZTjhwUJnnxVWl1OgNyYh5nyDc=",
        version = "v0.1.1",
    )

    go_repository(
        name = "com_github_form3tech_oss_jwt_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/form3tech-oss/jwt-go",
        sum = "h1:7ZaBxOI7TMoYBfyA3cQHErNNyAWIKUMIwqxEtgHOs5c=",
        version = "v3.2.3+incompatible",
    )
    go_repository(
        name = "com_github_fortytw2_leaktest",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fortytw2/leaktest",
        sum = "h1:u8491cBMTQ8ft8aeV+adlcytMZylmA5nnwwkRZjI8vw=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_foxcpp_go_mockdns",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/foxcpp/go-mockdns",
        sum = "h1:E52jfcE64UG42SwLmrW0QByONfGynWuzBvm86BoB9z8=",
        version = "v0.0.0-20210729171921-fb145fc6f897",
    )

    go_repository(
        name = "com_github_fsnotify_fsnotify",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fsnotify/fsnotify",
        sum = "h1:jRbGcIw6P2Meqdwuo0H1p6JVLbL5DHKAKlYndzMwVZI=",
        version = "v1.5.4",
    )

    go_repository(
        name = "com_github_fullstorydev_grpcurl",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fullstorydev/grpcurl",
        sum = "h1:p8BB6VZF8O7w6MxGr3KJ9E6EVKaswCevSALK6FBtMzA=",
        version = "v1.6.0",
    )

    go_repository(
        name = "com_github_fvbommel_sortorder",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fvbommel/sortorder",
        sum = "h1:dSnXLt4mJYH25uDDGa3biZNQsozaUWDSWeKJ0qqFfzE=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_fzipp_gocyclo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/fzipp/gocyclo",
        sum = "h1:A9UeX3HJSXTBzvHzhqoYVuE0eAhe+aM8XBCCwsPMZOc=",
        version = "v0.3.1",
    )

    go_repository(
        name = "com_github_getkin_kin_openapi",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/getkin/kin-openapi",
        sum = "h1:j77zg3Ec+k+r+GA3d8hBoXpAc6KX9TbBPrwQGBIy2sY=",
        version = "v0.76.0",
    )
    go_repository(
        name = "com_github_getsentry_raven_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/getsentry/raven-go",
        sum = "h1:no+xWJRb5ZI7eE8TWgIq1jLulQiIoLG0IfYxv5JYMGs=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_getsentry_sentry_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/getsentry/sentry-go",
        sum = "h1:rlOBkuFZRKKdUnKO+0U3JclRDQKlRu5vVQtkWSQvC70=",
        version = "v0.14.0",
    )
    go_repository(
        name = "com_github_ghodss_yaml",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ghodss/yaml",
        sum = "h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_gin_contrib_sse",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gin-contrib/sse",
        sum = "h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_gin_gonic_gin",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gin-gonic/gin",
        sum = "h1:4+fr/el88TOO3ewCmQr8cx/CtZ/umlIRIs5M4NTNjf8=",
        version = "v1.8.1",
    )

    go_repository(
        name = "com_github_git_lfs_git_lfs_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/git-lfs/git-lfs/v3",
        sum = "h1:pQQ1GYjxGJE92TqlWASyTAvNig1vS8mDWdWcXgSP3GA=",
        version = "v3.2.0",
    )

    go_repository(
        name = "com_github_git_lfs_gitobj_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/git-lfs/gitobj/v2",
        sum = "h1:5BUDAMga0Sv9msMXolrn6xplkiG5RaVEkOir2HSznog=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_github_git_lfs_go_netrc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/git-lfs/go-netrc",
        sum = "h1:6pskVZacdMUL93pCpMAYnMDLjH1yDFhssPYGe32sjdk=",
        version = "v0.0.0-20210914205454-f0c862dd687a",
    )
    go_repository(
        name = "com_github_git_lfs_pktline",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/git-lfs/pktline",
        sum = "h1:riQhgheTL7tMF4d5raz9t3+IzoR1i1wqxE1kZC6dY+U=",
        version = "v0.0.0-20210330133718-06e9096e2825",
    )

    go_repository(
        name = "com_github_git_lfs_wildmatch_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/git-lfs/wildmatch/v2",
        sum = "h1:Ds+aobrV5bK0wStILUOn9irllPyf9qrFETbKzwzoER8=",
        version = "v2.0.1",
    )
    go_repository(
        name = "com_github_go_critic_go_critic",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-critic/go-critic",
        sum = "h1:lS8B9LH/VVsvQQP7Ao5TJyQqteVKVs3E4dXiHMyubtI=",
        version = "v0.6.1",
    )

    go_repository(
        name = "com_github_go_enry_go_enry_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-enry/go-enry/v2",
        sum = "h1:BwvNrN58JqBJhyyVdZSl5QD3xoxEEGYUrRyPh31FGhw=",
        version = "v2.8.3",
    )

    go_repository(
        name = "com_github_go_enry_go_license_detector_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-enry/go-license-detector/v4",
        sum = "h1:OFlQAVNw5FlKUjX4OuW8JOabu8MQHjTKDb9pdeNYMUw=",
        version = "v4.3.0",
    )
    go_repository(
        name = "com_github_go_enry_go_oniguruma",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-enry/go-oniguruma",
        sum = "h1:k8aAMuJfMrqm/56SG2lV9Cfti6tC4x8673aHCcBk+eo=",
        version = "v1.2.1",
    )

    go_repository(
        name = "com_github_go_errors_errors",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-errors/errors",
        sum = "h1:J6MZopCL4uSllY1OfXM374weqZFFItUbrImctkmUxIA=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_go_git_gcfg",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-git/gcfg",
        sum = "h1:Q5ViNfGF8zFgyJWPqYwA7qGFoMTEiBmdlkcfRmpIMa4=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_go_git_go_billy_v5",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-git/go-billy/v5",
        sum = "h1:CPiOUAzKtMRvolEKw+bG1PLRpT7D3LIs3/3ey4Aiu34=",
        version = "v5.3.1",
    )

    go_repository(
        name = "com_github_go_git_go_git_v5",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-git/go-git/v5",
        sum = "h1:BXyZu9t0VkbiHtqrsvdq39UDhGJTl1h55VW6CSC4aY4=",
        version = "v5.4.2",
    )
    go_repository(
        name = "com_github_go_gl_glfw",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-gl/glfw",
        sum = "h1:QbL/5oDUmRBzO9/Z7Seo6zf912W/a6Sr4Eu0G/3Jho0=",
        version = "v0.0.0-20190409004039-e6da0acd62b1",
    )
    go_repository(
        name = "com_github_go_gl_glfw_v3_3_glfw",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-gl/glfw/v3.3/glfw",
        sum = "h1:WtGNWLvXpe6ZudgnXrq0barxBImvnnJoMEhXAzcbM0I=",
        version = "v0.0.0-20200222043503-6f7a984d4dc4",
    )
    go_repository(
        name = "com_github_go_gorp_gorp_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-gorp/gorp/v3",
        sum = "h1:ULqJXIekoqMx29FI5ekXXFoH1dT2Vc8UhnRzBg+Emz4=",
        version = "v3.0.2",
    )

    go_repository(
        name = "com_github_go_ini_ini",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-ini/ini",
        sum = "h1:z6ZrTEZqSWOTyH2FlglNbNgARyHG8oLW9gMELqKr06A=",
        version = "v1.67.0",
    )
    go_repository(
        name = "com_github_go_kit_kit",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-kit/kit",
        sum = "h1:wDJmvq38kDhkVxi50ni9ykkdUr1PKgqKOoi01fa0Mdk=",
        version = "v0.9.0",
    )
    go_repository(
        name = "com_github_go_kit_log",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-kit/log",
        sum = "h1:7i2K3eKTos3Vc0enKCfnVcgHh2olr/MyfboYq7cAcFw=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_go_logfmt_logfmt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-logfmt/logfmt",
        sum = "h1:otpy5pqBCBZ1ng9RQ0dPu4PN7ba75Y/aA+UpowDyNVA=",
        version = "v0.5.1",
    )
    go_repository(
        name = "com_github_go_logr_logr",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-logr/logr",
        sum = "h1:2DntVwHkVopvECVRSlL5PSo9eG+cAkDCuckLubN+rq0=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_go_logr_stdr",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-logr/stdr",
        sum = "h1:hSWxHoqTgW2S2qGc0LTAI563KZ5YKYRhT3MFKZMbjag=",
        version = "v1.2.2",
    )

    go_repository(
        name = "com_github_go_logr_zapr",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-logr/zapr",
        sum = "h1:a9vnzlIBPQBBkeaR9IuMUfmVOrQlkoC4YfPoFkX3T7A=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_go_martini_martini",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-martini/martini",
        sum = "h1:xveKWz2iaueeTaUgdetzel+U7exyigDYBryyVfV/rZk=",
        version = "v0.0.0-20170121215854-22fa46961aab",
    )
    go_repository(
        name = "com_github_go_ole_go_ole",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-ole/go-ole",
        sum = "h1:/Fpf6oFPoeFik9ty7siob0G6Ke8QvQEuVcuChpwXzpY=",
        version = "v1.2.6",
    )

    go_repository(
        name = "com_github_go_openapi_jsonpointer",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-openapi/jsonpointer",
        sum = "h1:gZr+CIYByUqjcgeLXnQu2gHYQC9o73G2XUeOFYEICuY=",
        version = "v0.19.5",
    )
    go_repository(
        name = "com_github_go_openapi_jsonreference",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-openapi/jsonreference",
        sum = "h1:MYlu0sBgChmCfJxxUKZ8g1cPWFOB37YSZqewK7OKeyA=",
        version = "v0.20.0",
    )

    go_repository(
        name = "com_github_go_openapi_swag",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-openapi/swag",
        sum = "h1:wm0rhTb5z7qpJRHBdPOMuY4QjVUMbF6/kwoYeRAOrKU=",
        version = "v0.21.1",
    )

    go_repository(
        name = "com_github_go_playground_assert_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-playground/assert/v2",
        sum = "h1:MsBgLAaY856+nPRTKrp3/OZK38U/wa0CcBYNjji3q3A=",
        version = "v2.0.1",
    )
    go_repository(
        name = "com_github_go_playground_locales",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-playground/locales",
        sum = "h1:u50s323jtVGugKlcYeyzC0etD1HifMjqmJqb8WugfUU=",
        version = "v0.14.0",
    )
    go_repository(
        name = "com_github_go_playground_universal_translator",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-playground/universal-translator",
        sum = "h1:82dyy6p4OuJq4/CByFNOn/jYrnRPArHwAcmLoJZxyho=",
        version = "v0.18.0",
    )
    go_repository(
        name = "com_github_go_playground_validator_v10",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-playground/validator/v10",
        sum = "h1:prmOlTVv+YjZjmRmNSF3VmspqJIxJWXmqUsHwfTRRkQ=",
        version = "v10.11.1",
    )
    go_repository(
        name = "com_github_go_redis_redis",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-redis/redis",
        sum = "h1:BKZuG6mCnRj5AOaWJXoCgf6rqTYnYJLe4en2hxT7r9o=",
        version = "v6.15.8+incompatible",
    )

    go_repository(
        name = "com_github_go_redis_redis_v8",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-redis/redis/v8",
        sum = "h1:AcZZR7igkdvfVmQTPnu9WE37LRrO/YrBH5zWyjDC0oI=",
        version = "v8.11.5",
    )
    go_repository(
        name = "com_github_go_redis_redismock_v8",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-redis/redismock/v8",
        sum = "h1:rtuijPgGynsRB2Y7KDACm09WvjHWS4RaG44Nm7rcj4Y=",
        version = "v8.0.6",
    )
    go_repository(
        name = "com_github_go_sql_driver_mysql",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-sql-driver/mysql",
        sum = "h1:BCTh4TKNUYmOmMUcQ3IipzF5prigylS7XXjEkfCHuOE=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_go_stack_stack",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-stack/stack",
        sum = "h1:5SgMzNM5HxrEjV0ww2lTmX6E2Izsfxas4+YHWRs3Lsk=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_github_go_task_slim_sprig",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-task/slim-sprig",
        sum = "h1:p104kn46Q8WdvHunIJ9dAyjPVtrBPhSr3KT2yUst43I=",
        version = "v0.0.0-20210107165309-348f09dbbbc0",
    )
    go_repository(
        name = "com_github_go_toolsmith_astcast",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/astcast",
        sum = "h1:JojxlmI6STnFVG9yOImLeGREv8W2ocNUM+iOhR6jE7g=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_toolsmith_astcopy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/astcopy",
        sum = "h1:OMgl1b1MEpjFQ1m5ztEO06rz5CUd3oBv9RF7+DyvdG8=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_toolsmith_astequal",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/astequal",
        sum = "h1:JbSszi42Jiqu36Gnf363HWS9MTEAz67vTQLponh3Moc=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_go_toolsmith_astfmt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/astfmt",
        sum = "h1:A0vDDXt+vsvLEdbMFJAUBI/uTbRw1ffOPnxsILnFL6k=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_toolsmith_astinfo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/astinfo",
        sum = "h1:wP6mXeB2V/d1P1K7bZ5vDUO3YqEzcvOREOxZPEu3gVI=",
        version = "v0.0.0-20180906194353-9809ff7efb21",
    )
    go_repository(
        name = "com_github_go_toolsmith_astp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/astp",
        sum = "h1:alXE75TXgcmupDsMK1fRAy0YUzLzqPVvBKoyWV+KPXg=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_toolsmith_pkgload",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/pkgload",
        sum = "h1:4DFWWMXVfbcN5So1sBNW9+yeiMqLFGl1wFLTL5R0Tgg=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_toolsmith_strparse",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/strparse",
        sum = "h1:Vcw78DnpCAKlM20kSbAyO4mPfJn/lyYA4BJUDxe2Jb4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_go_toolsmith_typep",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-toolsmith/typep",
        sum = "h1:8xdsa1+FSIH/RhEkgnD1j2CJOy5mNllW1Q9tRiYwvlk=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_go_xmlfmt_xmlfmt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/go-xmlfmt/xmlfmt",
        sum = "h1:khEcpUM4yFcxg4/FHQWkvVRmgijNXRfzkIDHh23ggEo=",
        version = "v0.0.0-20191208150333-d5b6f63a941b",
    )
    go_repository(
        name = "com_github_gobuffalo_logger",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gobuffalo/logger",
        sum = "h1:nnZNpxYo0zx+Aj9RfMPBm+x9zAU2OayFh/xrAWi34HU=",
        version = "v1.0.6",
    )
    go_repository(
        name = "com_github_gobuffalo_packd",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gobuffalo/packd",
        sum = "h1:U2wXfRr4E9DH8IdsDLlRFwTZTK7hLfq9qT/QHXGVe/0=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_gobuffalo_packr_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gobuffalo/packr/v2",
        sum = "h1:xE1yzvnO56cUC0sTpKR3DIbxZgB54AftTFMhB2XEWlY=",
        version = "v2.8.3",
    )

    go_repository(
        name = "com_github_gobwas_glob",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gobwas/glob",
        sum = "h1:A4xDbljILXROh+kObIiy5kIaPYD8e96x1tgBhUI5J+Y=",
        version = "v0.2.3",
    )

    go_repository(
        name = "com_github_gobwas_httphead",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gobwas/httphead",
        sum = "h1:s+21KNqlpePfkah2I+gwHF8xmJWRjooY+5248k6m4A0=",
        version = "v0.0.0-20180130184737-2c6c146eadee",
    )
    go_repository(
        name = "com_github_gobwas_pool",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gobwas/pool",
        sum = "h1:QEmUOlnSjWtnpRGHF3SauEiOsy82Cup83Vf2LcMlnc8=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_gobwas_ws",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gobwas/ws",
        sum = "h1:CoAavW/wd/kulfZmSIBt6p24n4j7tHgNVCjsfHVNUbo=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_goccy_go_json",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/goccy/go-json",
        sum = "h1:/pAaQDLHEoCq/5FFmSKBswWmK6H0e8g4159Kc/X/nqk=",
        version = "v0.9.11",
    )

    go_repository(
        name = "com_github_godbus_dbus_v5",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/godbus/dbus/v5",
        sum = "h1:mkgN1ofwASrYnJ5W6U/BxG15eXXXjirgZc7CLqkcaro=",
        version = "v5.0.6",
    )
    go_repository(
        name = "com_github_godror_godror",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/godror/godror",
        sum = "h1:uxGAD7UdnNGjX5gf4NnEIGw0JAPTIFiqAyRBZTPKwXs=",
        version = "v0.24.2",
    )
    go_repository(
        name = "com_github_gofrs_flock",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gofrs/flock",
        sum = "h1:+gYjHKf32LDeiEEFhQaotPbLuUXjY5ZqxKgXy7n59aw=",
        version = "v0.8.1",
    )

    go_repository(
        name = "com_github_gogo_googleapis",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gogo/googleapis",
        sum = "h1:zgVt4UpGxcqVOw97aRGxT4svlcmdK35fynLNctY32zI=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_gogo_protobuf",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gogo/protobuf",
        sum = "h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=",
        version = "v1.3.2",
    )

    go_repository(
        name = "com_github_golang_glog",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golang/glog",
        sum = "h1:nfP3RFugxnNRyKgeWd4oI1nYvXpxrx8ck8ZrcizshdQ=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_golang_groupcache",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golang/groupcache",
        sum = "h1:oI5xCqsCo564l8iNU+DwB5epxmsaqB+rhGL0m5jtYqE=",
        version = "v0.0.0-20210331224755-41bb18bfe9da",
    )
    go_repository(
        name = "com_github_golang_jwt_jwt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golang-jwt/jwt",
        sum = "h1:IfV12K8xAKAnZqdXVzCZ+TOjboZ2keLg81eXfW3O+oY=",
        version = "v3.2.2+incompatible",
    )
    go_repository(
        name = "com_github_golang_jwt_jwt_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golang-jwt/jwt/v4",
        sum = "h1:rcc4lwaZgFMCZ5jxF9ABolDcIHdBytAFgqFPbSJQAYs=",
        version = "v4.4.2",
    )

    go_repository(
        name = "com_github_golang_mock",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golang/mock",
        sum = "h1:ErTB+efbowRARo13NNdxyJji2egdxLGQhRaY+DUumQc=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_golang_protobuf",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golang/protobuf",
        sum = "h1:ROPKBNFfQgOUMifHyP+KYbvpjbdoFNs+aK7DXlji0Tw=",
        version = "v1.5.2",
    )
    go_repository(
        name = "com_github_golang_snappy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golang/snappy",
        sum = "h1:yAGX7huGHXlcLOEtBnF4w7FQwA26wojNCwOYAEhLjQM=",
        version = "v0.0.4",
    )
    go_repository(
        name = "com_github_golang_sql_civil",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golang-sql/civil",
        sum = "h1:lXe2qZdvpiX5WZkZR4hgp4KJVfY3nMkvmwbVkpv1rVY=",
        version = "v0.0.0-20190719163853-cb61b32ac6fe",
    )
    go_repository(
        name = "com_github_golangci_check",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/check",
        sum = "h1:23T5iq8rbUYlhpt5DB4XJkc6BU31uODLD1o1gKvZmD0=",
        version = "v0.0.0-20180506172741-cfe4005ccda2",
    )
    go_repository(
        name = "com_github_golangci_dupl",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/dupl",
        sum = "h1:w8hkcTqaFpzKqonE9uMCefW1WDie15eSP/4MssdenaM=",
        version = "v0.0.0-20180902072040-3e9179ac440a",
    )
    go_repository(
        name = "com_github_golangci_go_misc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/go-misc",
        sum = "h1:9kfjN3AdxcbsZBf8NjltjWihK2QfBBBZuv91cMFfDHw=",
        version = "v0.0.0-20180628070357-927a3d87b613",
    )
    go_repository(
        name = "com_github_golangci_gofmt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/gofmt",
        sum = "h1:iR3fYXUjHCR97qWS8ch1y9zPNsgXThGwjKPrYfqMPks=",
        version = "v0.0.0-20190930125516-244bba706f1a",
    )
    go_repository(
        name = "com_github_golangci_golangci_lint",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/golangci-lint",
        sum = "h1:SLwZFEmDgopqZpfP495zCtV9REUf551JJlJ51Ql7NZA=",
        version = "v1.43.0",
    )
    go_repository(
        name = "com_github_golangci_lint_1",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/lint-1",
        sum = "h1:MfyDlzVjl1hoaPzPD4Gpb/QgoRfSBR0jdhwGyAWwMSA=",
        version = "v0.0.0-20191013205115-297bf364a8e0",
    )
    go_repository(
        name = "com_github_golangci_maligned",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/maligned",
        sum = "h1:kNY3/svz5T29MYHubXix4aDDuE3RWHkPvopM/EDv/MA=",
        version = "v0.0.0-20180506175553-b1d89398deca",
    )
    go_repository(
        name = "com_github_golangci_misspell",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/misspell",
        sum = "h1:pLzmVdl3VxTOncgzHcvLOKirdvcx/TydsClUQXTehjo=",
        version = "v0.3.5",
    )
    go_repository(
        name = "com_github_golangci_revgrep",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/revgrep",
        sum = "h1:SgM7GDZTxtTTQPU84heOxy34iG5Du7F2jcoZnvp+fXI=",
        version = "v0.0.0-20210930125155-c22e5001d4f2",
    )
    go_repository(
        name = "com_github_golangci_unconvert",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/golangci/unconvert",
        sum = "h1:zwtduBRr5SSWhqsYNgcuWO2kFlpdOZbP0+yRjmvPGys=",
        version = "v0.0.0-20180507085042-28b1c447d1f4",
    )

    go_repository(
        name = "com_github_gomodule_redigo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gomodule/redigo",
        sum = "h1:H5XSIre1MB5NbPYFp+i1NBbb5qN1W8Y8YAQoAYbkm8k=",
        version = "v1.8.2",
    )
    go_repository(
        name = "com_github_google_btree",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/btree",
        sum = "h1:gK4Kx5IaGY9CD5sPJ36FHiBJ6ZXl0kilRiiCj+jdYp4=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_google_cel_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/cel-go",
        sum = "h1:DmzaiSgoaqGCjtpPQWl26/gND+yRpim56H1jCVev6d8=",
        version = "v0.12.5",
    )

    go_repository(
        name = "com_github_google_certificate_transparency_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/certificate-transparency-go",
        sum = "h1:6JHXZhXEvilMcTjR4MGZn5KV0IRkcFl4CJx5iHVhjFE=",
        version = "v1.1.1",
    )

    go_repository(
        name = "com_github_google_flatbuffers",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/flatbuffers",
        sum = "h1:MVlul7pQNoDzWRLTw5imwYsl+usrS1TXG2H4jg6ImGw=",
        version = "v1.12.1",
    )
    go_repository(
        name = "com_github_google_gnostic",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/gnostic",
        sum = "h1:ZK/5VhkoX835RikCHpSUJV9a+S3e1zLh59YnyWeBW+0=",
        version = "v0.6.9",
    )

    go_repository(
        name = "com_github_google_go_cmp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/go-cmp",
        sum = "h1:O2Tfq5qg4qc4AmwVlvv0oLiVAGB7enBSJ2x2DqQFi38=",
        version = "v0.5.9",
    )
    go_repository(
        name = "com_github_google_go_containerregistry",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/go-containerregistry",
        sum = "h1:Xt8x1adcREjFcmDoDK8OdOsjxu90PHkGuwNP8GiHMLM=",
        version = "v0.11.0",
    )

    go_repository(
        name = "com_github_google_gofuzz",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/gofuzz",
        sum = "h1:xRy4A+RhZaiKjJ1bPfwQ8sedCA+YS2YcCHW6ec7JMi0=",
        version = "v1.2.0",
    )

    go_repository(
        name = "com_github_google_martian",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/martian",
        sum = "h1:/CP5g8u/VJHijgedC/Legn3BAbAaWPgecwXBIDzw5no=",
        version = "v2.1.0+incompatible",
    )
    go_repository(
        name = "com_github_google_martian_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/martian/v3",
        sum = "h1:d8MncMlErDFTwQGBK1xhv026j9kqhvw1Qv9IbWT1VLQ=",
        version = "v3.2.1",
    )
    go_repository(
        name = "com_github_google_pprof",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/pprof",
        sum = "h1:VrKTY4lquiy1oJzVZgXrauku9Jx9P+POv/gTLakG4Wk=",
        version = "v0.0.0-20220412212628-83db2b799d1f",
    )
    go_repository(
        name = "com_github_google_renameio",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/renameio",
        sum = "h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_google_shlex",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/shlex",
        sum = "h1:El6M4kTTCOh6aBiKaUGG7oYTSPP8MxqL4YI3kZKwcP4=",
        version = "v0.0.0-20191202100458-e7afc7fbc510",
    )
    go_repository(
        name = "com_github_google_trillian",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/trillian",
        sum = "h1:pPzJPkK06mvXId1LHEAJxIegGgHzzp/FUnycPYfoCMI=",
        version = "v1.3.11",
    )

    go_repository(
        name = "com_github_google_uuid",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/uuid",
        sum = "h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_google_wire",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/google/wire",
        sum = "h1:I7ELFeVBr3yfPIcc8+MWvrjk+3VjbcSzoXm3JVa+jD8=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_googleapis_enterprise_certificate_proxy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/googleapis/enterprise-certificate-proxy",
        sum = "h1:y8Yozv7SZtlU//QXbezB6QkpuE6jMD2/gfzk4AftXjs=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_github_googleapis_gax_go_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/googleapis/gax-go/v2",
        sum = "h1:SXk3ABtQYDT/OH8jAyvEOQ58mgawq5C4o/4/89qN2ZU=",
        version = "v2.6.0",
    )

    go_repository(
        name = "com_github_googleapis_go_type_adapters",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/googleapis/go-type-adapters",
        sum = "h1:9XdMn+d/G57qq1s8dNc5IesGCXHf6V2HZ2JwRxfA2tA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_googleapis_google_cloud_go_testing",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/googleapis/google-cloud-go-testing",
        sum = "h1:tlyzajkF3030q6M8SvmJSemC9DTHL/xaMa18b65+JM4=",
        version = "v0.0.0-20200911160855-bcd43fbb19e8",
    )

    go_repository(
        name = "com_github_gookit_color",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gookit/color",
        sum = "h1:tXy44JFSFkKnELV6WaMo/lLfu/meqITX3iAV52do7lk=",
        version = "v1.4.2",
    )

    go_repository(
        name = "com_github_gopherjs_gopherjs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gopherjs/gopherjs",
        sum = "h1:EGx4pi6eqNxGaHF6qqu48+N2wcFQ5qg5FXgOdqsJ5d8=",
        version = "v0.0.0-20181017120253-0766667cb4d1",
    )
    go_repository(
        name = "com_github_gordonklaus_ineffassign",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gordonklaus/ineffassign",
        sum = "h1:Nb2aRlC404yz7gQIfRZxX9/MLvQiqXyiBTJtgAy6yrI=",
        version = "v0.0.0-20210225214923-2e10b2664254",
    )

    go_repository(
        name = "com_github_gorhill_cronexpr",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gorhill/cronexpr",
        sum = "h1:f0n1xnMSmBLzVfsMMvriDyA75NB/oBgILX2GcHXIQzY=",
        version = "v0.0.0-20180427100037-88b0669f7d75",
    )
    go_repository(
        name = "com_github_gorilla_css",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gorilla/css",
        sum = "h1:BQqNyPTi50JCFMTw/b67hByjMVXZRwGha6wxVGkeihY=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_gorilla_handlers",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gorilla/handlers",
        sum = "h1:9lRY6j8DEeeBT10CvO9hGW0gmky0BprnvDI5vfhUHH4=",
        version = "v1.5.1",
    )
    go_repository(
        name = "com_github_gorilla_mux",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gorilla/mux",
        sum = "h1:i40aqfkR1h2SlN9hojwV5ZA91wcXFOvkdNIeFDP5koI=",
        version = "v1.8.0",
    )

    go_repository(
        name = "com_github_gorilla_websocket",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gorilla/websocket",
        sum = "h1:+/TMaTYc4QFitKJxsQ7Yye35DkWvkdLcvGKqM+x0Ufc=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_gostaticanalysis_analysisutil",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gostaticanalysis/analysisutil",
        sum = "h1:ZMCjoue3DtDWQ5WyU16YbjbQEQ3VuzwxALrpYd+HeKk=",
        version = "v0.7.1",
    )
    go_repository(
        name = "com_github_gostaticanalysis_comment",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gostaticanalysis/comment",
        sum = "h1:hlnx5+S2fY9Zo9ePo4AhgYsYHbM2+eAv8m/s1JiCd6Q=",
        version = "v1.4.2",
    )
    go_repository(
        name = "com_github_gostaticanalysis_forcetypeassert",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gostaticanalysis/forcetypeassert",
        sum = "h1:rx8127mFPqXXsfPSo8BwnIU97MKFZc89WHAHt8PwDVY=",
        version = "v0.0.0-20200621232751-01d4955beaa5",
    )
    go_repository(
        name = "com_github_gostaticanalysis_nilerr",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gostaticanalysis/nilerr",
        sum = "h1:ThE+hJP0fEp4zWLkWHWcRyI2Od0p7DlgYG3Uqrmrcpk=",
        version = "v0.1.1",
    )
    go_repository(
        name = "com_github_gostaticanalysis_testutil",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gostaticanalysis/testutil",
        sum = "h1:nhdCmubdmDF6VEatUNjgUZBJKWRqugoISdUv3PPQgHY=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_gosuri_uitable",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gosuri/uitable",
        sum = "h1:IG2xLKRvErL3uhY6e1BylFzG+aJiwQviDDTfOKeKTpY=",
        version = "v0.0.4",
    )

    go_repository(
        name = "com_github_gregjones_httpcache",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/gregjones/httpcache",
        sum = "h1:+ngKgrYPPJrOjhax5N+uePQ0Fh1Z7PheYoUI/0nzkPA=",
        version = "v0.0.0-20190611155906-901d90724c79",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_middleware",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/grpc-ecosystem/go-grpc-middleware",
        sum = "h1:+9834+KizmvFV7pXQGSXQTsaWhq2GjuNUt0aUU0YBYw=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_go_grpc_prometheus",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
        sum = "h1:tDEI6JHij2b2WFuvCLn6ZFk4WcpH+lWJPrDdIAgeBuQ=",
        version = "v1.2.1-0.20200507082539-9abf3eb82b4a",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_gateway",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/grpc-ecosystem/grpc-gateway",
        sum = "h1:gmcG1KaJ57LophUzW0Hy8NmPhnMZb4M0+kPpLofRdBo=",
        version = "v1.16.0",
    )
    go_repository(
        name = "com_github_grpc_ecosystem_grpc_gateway_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/grpc-ecosystem/grpc-gateway/v2",
        sum = "h1:BZHcxBETFHIdVyhyEfOvn/RdU/QGdLI4y34qQGjGWO0=",
        version = "v2.7.0",
    )

    go_repository(
        name = "com_github_hashicorp_consul_api",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/consul/api",
        sum = "h1:MwZJp86nlnL+6+W1Zly4JUuVn9YHhMggBirMpHGD7kw=",
        version = "v1.10.1",
    )
    go_repository(
        name = "com_github_hashicorp_consul_sdk",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/consul/sdk",
        sum = "h1:OJtKBtEjboEZvG6AOUdh4Z1Zbyu0WcxQ0qatRrZHTVU=",
        version = "v0.8.0",
    )
    go_repository(
        name = "com_github_hashicorp_errwrap",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/errwrap",
        sum = "h1:OxrOeh75EUXMY8TBjag2fzXGZ40LB6IKw45YeGUDY2I=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_cleanhttp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-cleanhttp",
        sum = "h1:dH3aiDG9Jvb5r5+bYHsikaOUIpcM0xvgMXVoDkXMzJM=",
        version = "v0.5.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_hclog",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-hclog",
        sum = "h1:d4QkX8FRTYaKaCZBoXYY8zJX2BXjWxurN/GA2tkrmZM=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_immutable_radix",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-immutable-radix",
        sum = "h1:AKDB1HM5PWEA7i4nhcpwOrO2byshxBjXVn/J/3+z5/0=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_msgpack",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-msgpack",
        sum = "h1:zKjpN5BK/P5lMYrLmBHdBULWbJ0XpYR+7NGzqkZzoD4=",
        version = "v0.5.3",
    )
    go_repository(
        name = "com_github_hashicorp_go_multierror",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-multierror",
        sum = "h1:H5DkEtf6CXdFp0N0Em5UCwQpXMWke8IA0+lD48awMYo=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_net",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go.net",
        sum = "h1:sNCoNyDEvN1xa+X0baata4RdcpKwcMS6DH+xwfqPgjw=",
        version = "v0.0.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_retryablehttp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-retryablehttp",
        sum = "h1:sUiuQAnLlbvmExtFQs72iFW/HXeUn8Z1aJLQ4LJJbTQ=",
        version = "v0.7.1",
    )
    go_repository(
        name = "com_github_hashicorp_go_rootcerts",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-rootcerts",
        sum = "h1:jzhAVGtqPKbwpyCPELlgNWhE1znq+qwJtW5Oi2viEzc=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_hashicorp_go_sockaddr",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-sockaddr",
        sum = "h1:GeH6tui99pF4NJgfnhp+L6+FfobzVW3Ah46sLo0ICXs=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_syslog",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-syslog",
        sum = "h1:KaodqZuhUoZereWVIYmpUgZysurB1kBLX2j0MwMrUAE=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hashicorp_go_uuid",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-uuid",
        sum = "h1:2gKiV6YVmrJ1i2CKKa9obLvRieoRGviZFL26PcT/Co8=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_hashicorp_go_version",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/go-version",
        sum = "h1:O293SZ2Eg+AAYijkVK3jR786Am1bhDEh2GHT0tIVE5E=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_hashicorp_golang_lru",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/golang-lru",
        sum = "h1:YDjusn29QI/Das2iO9M0BHnIbxPeyuCHsjMW+lJfyTc=",
        version = "v0.5.4",
    )
    go_repository(
        name = "com_github_hashicorp_hcl",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/hcl",
        sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_hashicorp_logutils",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/logutils",
        sum = "h1:dLEQVugN8vlakKOUE3ihGLTZJRB4j+M2cdTm/ORI65Y=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_hashicorp_mdns",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/mdns",
        sum = "h1:XFSOubp8KWB+Jd2PDyaX5xUd5bhSP/+pTDZVDMzZJM8=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_hashicorp_memberlist",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/memberlist",
        sum = "h1:5+RffWKwqJ71YPu9mWsF7ZOscZmwfasdA8kbdC7AO2g=",
        version = "v0.2.2",
    )
    go_repository(
        name = "com_github_hashicorp_serf",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/serf",
        sum = "h1:EBWvyu9tcRszt3Bxp3KNssBMP1KuHWyO51lz9+786iM=",
        version = "v0.9.5",
    )
    go_repository(
        name = "com_github_hashicorp_yamux",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hashicorp/yamux",
        sum = "h1:yrQxtgseBDrq9Y652vSRDvsKCJKOUD+GzTS4Y0Y8pvE=",
        version = "v0.1.1",
    )
    go_repository(
        name = "com_github_hdrhistogram_hdrhistogram_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/HdrHistogram/hdrhistogram-go",
        sum = "h1:cJXY5VLMHgejurPjZH6Fo9rIwRGLefBGdiaENZALqrg=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_hhatto_gorst",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hhatto/gorst",
        sum = "h1:Jdu2tbAxkRouSILp2EbposIb8h4gO+2QuZEn3d9sKAc=",
        version = "v0.0.0-20181029133204-ca9f730cac5b",
    )
    go_repository(
        name = "com_github_hpcloud_tail",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/hpcloud/tail",
        sum = "h1:nfCOvKYfkgYP8hkirhJocXT2+zOD8yUNjXaWfTlyFKI=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_huandu_xstrings",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/huandu/xstrings",
        sum = "h1:L18LIDzqlW6xN2rEkpdV8+oL/IXWJ1APd+vsdYy4Wdw=",
        version = "v1.3.2",
    )

    go_repository(
        name = "com_github_iancoleman_strcase",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/iancoleman/strcase",
        sum = "h1:05I4QRnGpI0m37iZQRuskXh+w77mr6Z41lwQzuHLwW0=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_ianlancetaylor_demangle",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ianlancetaylor/demangle",
        sum = "h1:uGg2frlt3IcT7kbV6LEp5ONv4vmoO2FW4qSO+my/aoM=",
        version = "v0.0.0-20210905161508-09a460cdf81d",
    )
    go_repository(
        name = "com_github_imdario_mergo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/imdario/mergo",
        sum = "h1:lFzP57bqS/wsqKssCGmtLAb8A0wKjLGrve2q3PPVcBk=",
        version = "v0.3.13",
    )

    go_repository(
        name = "com_github_inconshreveable_mousetrap",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/inconshreveable/mousetrap",
        sum = "h1:U3uMjPSQEBMNp1lFxmllqCPM6P5u/Xq7Pgzkat/bFNc=",
        version = "v1.0.1",
    )

    go_repository(
        name = "com_github_intel_goresctrl",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/intel/goresctrl",
        sum = "h1:JyZjdMQu9Kl/wLXe9xA6s1X+tF6BWsQPFGJMEeCfWzE=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_github_iris_contrib_jade",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/iris-contrib/jade",
        sum = "h1:WoYdfyJFfZIUgqNAeOyRfTNQZOksSlZ6+FnXR3AEpX0=",
        version = "v1.1.4",
    )

    go_repository(
        name = "com_github_iris_contrib_schema",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/iris-contrib/schema",
        sum = "h1:CPSBLyx2e91H2yJzPuhGuifVRnZBBJ3pCOMbOvPZaTw=",
        version = "v0.0.6",
    )

    go_repository(
        name = "com_github_jackc_chunkreader_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jackc/chunkreader/v2",
        sum = "h1:i+RDz65UE+mmpjTfyz0MoVTnzeYxroil2G82ki7MGG8=",
        version = "v2.0.1",
    )
    go_repository(
        name = "com_github_jackc_pgconn",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jackc/pgconn",
        sum = "h1:3L1XMNV2Zvca/8BYhzcRFS70Lr0WlDg16Di6SFGAbys=",
        version = "v1.13.0",
    )
    go_repository(
        name = "com_github_jackc_pgio",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jackc/pgio",
        sum = "h1:g12B9UwVnzGhueNavwioyEEpAmqMe1E/BN9ES+8ovkE=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_jackc_pgpassfile",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jackc/pgpassfile",
        sum = "h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_jackc_pgproto3_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jackc/pgproto3/v2",
        sum = "h1:nwj7qwf0S+Q7ISFfBndqeLwSwxs+4DPsbRFjECT1Y4Y=",
        version = "v2.3.1",
    )
    go_repository(
        name = "com_github_jackc_pgservicefile",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jackc/pgservicefile",
        sum = "h1:C8S2+VttkHFdOOCXJe+YGfa4vHYwlt4Zx+IVXQ97jYg=",
        version = "v0.0.0-20200714003250-2b9c44734f2b",
    )
    go_repository(
        name = "com_github_jackc_pgtype",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jackc/pgtype",
        sum = "h1:Dlq8Qvcch7kiehm8wPGIW0W3KsCCHJnRacKW0UM8n5w=",
        version = "v1.12.0",
    )
    go_repository(
        name = "com_github_jackc_pgx_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jackc/pgx/v4",
        sum = "h1:0Ut0rpeKwvIVbMQ1KbMBU4h6wxehBI535LK6Flheh8E=",
        version = "v4.17.2",
    )

    go_repository(
        name = "com_github_jbenet_go_context",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jbenet/go-context",
        sum = "h1:BQSFePA1RWJOlocH6Fxy8MmwDt+yVQYULKfN0RoTN8A=",
        version = "v0.0.0-20150711004518-d14ea06fba99",
    )
    go_repository(
        name = "com_github_jcmturner_aescts_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jcmturner/aescts/v2",
        sum = "h1:9YKLH6ey7H4eDBXW8khjYslgyqG2xZikXP0EQFKrle8=",
        version = "v2.0.0",
    )
    go_repository(
        name = "com_github_jcmturner_dnsutils_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jcmturner/dnsutils/v2",
        sum = "h1:lltnkeZGL0wILNvrNiVCR6Ro5PGU/SeBvVO/8c/iPbo=",
        version = "v2.0.0",
    )

    go_repository(
        name = "com_github_jcmturner_gofork",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jcmturner/gofork",
        sum = "h1:J7uCkflzTEhUZ64xqKnkDxq3kzc96ajM1Gli5ktUem8=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_jcmturner_goidentity_v6",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jcmturner/goidentity/v6",
        sum = "h1:VKnZd2oEIMorCTsFBnJWbExfNN7yZr3EhJAxwOkZg6o=",
        version = "v6.0.1",
    )
    go_repository(
        name = "com_github_jcmturner_gokrb5_v8",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jcmturner/gokrb5/v8",
        sum = "h1:6ZIM6b/JJN0X8UM43ZOM6Z4SJzla+a/u7scXFJzodkA=",
        version = "v8.4.2",
    )
    go_repository(
        name = "com_github_jcmturner_rpc_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jcmturner/rpc/v2",
        sum = "h1:7FXXj8Ti1IaVFpSAziCZWNzbNuZmnvw/i6CqLNdWfZY=",
        version = "v2.0.3",
    )

    go_repository(
        name = "com_github_jdkato_prose",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jdkato/prose",
        sum = "h1:LpvmDGwbKGTgdCH3a8VJL56sr7p/wOFPw/R4lM4PfFg=",
        version = "v1.1.0",
    )

    go_repository(
        name = "com_github_jessevdk_go_flags",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jessevdk/go-flags",
        sum = "h1:4IU2WS7AumrZ/40jfhf4QVDMsQwqA7VEHozFRrGARJA=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_jgautheron_goconst",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jgautheron/goconst",
        sum = "h1:HxVbL1MhydKs8R8n/HE5NPvzfaYmQJA3o879lE4+WcM=",
        version = "v1.5.1",
    )
    go_repository(
        name = "com_github_jhump_protoreflect",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jhump/protoreflect",
        sum = "h1:4/2yi5LyDPP7nN+Hiird1SAJ6YoxUm13/oxHGRnbPd8=",
        version = "v1.6.1",
    )
    go_repository(
        name = "com_github_jingyugao_rowserrcheck",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jingyugao/rowserrcheck",
        sum = "h1:zibz55j/MJtLsjP1OF4bSdgXxwL1b+Vn7Tjzq7gFzUs=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_jirfag_go_printf_func_name",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jirfag/go-printf-func-name",
        sum = "h1:KA9BjwUk7KlCh6S9EAGWBt1oExIUv9WyNCiRz5amv48=",
        version = "v0.0.0-20200119135958-7558a9eaa5af",
    )

    go_repository(
        name = "com_github_jmespath_go_jmespath",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jmespath/go-jmespath",
        sum = "h1:BEgLn5cpjn8UN1mAw4NjwDrS35OdebyEtFe+9YPoQUg=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_jmespath_go_jmespath_internal_testify",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jmespath/go-jmespath/internal/testify",
        sum = "h1:shLQSRRSCCPj3f2gpwzGwWFoC7ycTf1rcQZHOlsJ6N8=",
        version = "v1.5.1",
    )

    go_repository(
        name = "com_github_jmoiron_sqlx",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jmoiron/sqlx",
        sum = "h1:vFFPA71p1o5gAeqtEAwLU4dnX2napprKtHr7PYIcN3g=",
        version = "v1.3.5",
    )

    go_repository(
        name = "com_github_jonboulle_clockwork",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jonboulle/clockwork",
        sum = "h1:UOGuzwb1PwsrDAObMuhUnj0p5ULPj8V/xJ7Kx9qUBdQ=",
        version = "v0.2.2",
    )
    go_repository(
        name = "com_github_josharian_intern",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/josharian/intern",
        sum = "h1:vlS4z54oSdjm0bgjRigI+G1HpF+tI+9rE5LLzOg8HmY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_josharian_txtarfs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/josharian/txtarfs",
        sum = "h1:8NZHLa6Gp0hW6xJ0c3F1Kse7dJw30fOcDzHuF9sLbnE=",
        version = "v0.0.0-20210218200122-0702f000015a",
    )

    go_repository(
        name = "com_github_jpillora_backoff",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jpillora/backoff",
        sum = "h1:uvFg412JmmHBHw7iwprIxkPMI+sGQ4kzOWsMeHnm2EA=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_json_iterator_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/json-iterator/go",
        sum = "h1:PV8peI4a0ysnczrg+LtxykD8LfKY9ML6u2jnxaEnrnM=",
        version = "v1.1.12",
    )
    go_repository(
        name = "com_github_jstemmer_go_junit_report",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jstemmer/go-junit-report",
        sum = "h1:6QPYqodiu3GuPL+7mfx+NwDdp2eTkp9IfEUpgAwUN0o=",
        version = "v0.9.1",
    )

    go_repository(
        name = "com_github_jtolds_gls",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/jtolds/gls",
        sum = "h1:xdiiI2gbIgH/gLH7ADydsJ1uDOEzR8yvV7C0MuV77Wo=",
        version = "v4.20.0+incompatible",
    )
    go_repository(
        name = "com_github_juju_ratelimit",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/juju/ratelimit",
        sum = "h1:+7AIFJVQ0EQgq/K9+0Krm7m530Du7tIz0METWzN0RgY=",
        version = "v1.0.1",
    )

    go_repository(
        name = "com_github_julienschmidt_httprouter",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/julienschmidt/httprouter",
        sum = "h1:U0609e9tgbseu3rBINet9P48AI/D3oJs4dN7jwJOQ1U=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_github_julz_importas",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/julz/importas",
        sum = "h1:XeSMXURZPtUffuWAaq90o6kLgZdgu+QA8wk4MPC8ikI=",
        version = "v0.0.0-20210419104244-841f0c0fe66d",
    )
    go_repository(
        name = "com_github_k0kubun_colorstring",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/k0kubun/colorstring",
        sum = "h1:uC1QfSlInpQF+M0ao65imhwqKnz3Q2z/d8PWZRMQvDM=",
        version = "v0.0.0-20150214042306-9440f1994b88",
    )
    go_repository(
        name = "com_github_karrick_godirwalk",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/karrick/godirwalk",
        sum = "h1:DynhcF+bztK8gooS0+NDJFrdNZjJ3gzVzC545UNA9iw=",
        version = "v1.16.1",
    )
    go_repository(
        name = "com_github_kataras_blocks",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kataras/blocks",
        sum = "h1:cF3RDY/vxnSRezc7vLFlQFTYXG/yAr1o7WImJuZbzC4=",
        version = "v0.0.7",
    )

    go_repository(
        name = "com_github_kataras_golog",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kataras/golog",
        sum = "h1:0TY5tHn5L5DlRIikepcaRR/6oInIr9AiWsxzt0vvlBE=",
        version = "v0.1.7",
    )
    go_repository(
        name = "com_github_kataras_iris_v12",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kataras/iris/v12",
        sum = "h1:grB/oCf5baZhmYIeDMfgN3LYrtEcmK8pbxlRvEZ2pgw=",
        version = "v12.2.0-beta5",
    )

    go_repository(
        name = "com_github_kataras_pio",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kataras/pio",
        sum = "h1:kqreJ5KOEXGMwHAWHDwIl+mjfNCPhAwZPa8gK7MKlyw=",
        version = "v0.0.11",
    )
    go_repository(
        name = "com_github_kataras_sitemap",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kataras/sitemap",
        sum = "h1:w71CRMMKYMJh6LR2wTgnk5hSgjVNB9KL60n5e2KHvLY=",
        version = "v0.0.6",
    )
    go_repository(
        name = "com_github_kataras_tunnel",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kataras/tunnel",
        sum = "h1:sCAqWuJV7nPzGrlb0os3j49lk2JhILT0rID38NHNLpA=",
        version = "v0.0.4",
    )

    go_repository(
        name = "com_github_kelseyhightower_envconfig",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kelseyhightower/envconfig",
        sum = "h1:Im6hONhd3pLkfDFsbRgu68RDNkGF1r3dvMUtDTo2cv8=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_github_kevinburke_ssh_config",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kevinburke/ssh_config",
        sum = "h1:DowS9hvgyYSX4TO5NpyC606/Z4SxnNYbT+WX27or6Ck=",
        version = "v0.0.0-20201106050909-4977a11b4351",
    )

    go_repository(
        name = "com_github_kisielk_errcheck",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kisielk/errcheck",
        sum = "h1:YTDO4pNy7AUN/021p+JGHycQyYNIyMoenM1YDVK6RlY=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_kisielk_gotool",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kisielk/gotool",
        sum = "h1:AV2c/EiW3KqPNT9ZKl07ehoAGi4C5/01Cfbblndcapg=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_klauspost_compress",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/klauspost/compress",
        sum = "h1:Lcadnb3RKGin4FYM/orgq0qde+nc15E5Cbqg4B9Sx9c=",
        version = "v1.15.11",
    )

    go_repository(
        name = "com_github_konsorten_go_windows_terminal_sequences",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/konsorten/go-windows-terminal-sequences",
        sum = "h1:CE8S1cTafDpPvMhIxNJKvHsGVBgn1xWYf1NbHQhywc8=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_kortschak_utter",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kortschak/utter",
        sum = "h1:AJVccwLrdrikvkH0aI5JKlzZIORLpfMeGBQ5tHfIXis=",
        version = "v1.0.1",
    )

    go_repository(
        name = "com_github_kr_fs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kr/fs",
        sum = "h1:Jskdu9ieNAYnjxsi0LbQp1ulIKZV1LAFgK1tWhpZgl8=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_kr_logfmt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kr/logfmt",
        sum = "h1:T+h1c/A9Gawja4Y9mFVWj2vyii2bbUNDw3kt9VxK2EY=",
        version = "v0.0.0-20140226030751-b84e30acd515",
    )
    go_repository(
        name = "com_github_kr_pretty",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kr/pretty",
        sum = "h1:Fmg33tUaq4/8ym9TJN1x7sLJnHVwhP33CNkpYV/7rwI=",
        version = "v0.2.1",
    )
    go_repository(
        name = "com_github_kr_pty",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kr/pty",
        sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_kr_text",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kr/text",
        sum = "h1:5Nx0Ya0ZqY2ygV366QzturHI13Jq95ApcVaJBhpS+AY=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_kulti_thelper",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kulti/thelper",
        sum = "h1:2Nx7XbdbE/BYZeoip2mURKUdtHQRuy6Ug+wR7K9ywNM=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_kunwardeep_paralleltest",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kunwardeep/paralleltest",
        sum = "h1:UdKIkImEAXjR1chUWLn+PNXqWUGs//7tzMeWuP7NhmI=",
        version = "v1.0.3",
    )
    go_repository(
        name = "com_github_kylelemons_godebug",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kylelemons/godebug",
        sum = "h1:RPNrshWIDI6G2gRW9EHilWtl7Z6Sb1BR0xunSBf0SNc=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_kyoh86_exportloopref",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/kyoh86/exportloopref",
        sum = "h1:5Ry/at+eFdkX9Vsdw3qU4YkvGtzuVfzT4X7S77LoN/M=",
        version = "v0.1.8",
    )

    go_repository(
        name = "com_github_labstack_echo_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/labstack/echo/v4",
        sum = "h1:wPOF1CE6gvt/kmbMR4dGzWvHMPT+sAEUJOwOTtvITVY=",
        version = "v4.9.0",
    )
    go_repository(
        name = "com_github_labstack_gommon",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/labstack/gommon",
        sum = "h1:OomWaJXm7xR6L1HmEtGyQf26TEn7V6X88mktX9kee9o=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_lann_builder",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/lann/builder",
        sum = "h1:SOEGU9fKiNWd/HOJuq6+3iTQz8KNCLtVX6idSoTLdUw=",
        version = "v0.0.0-20180802200727-47ae307949d0",
    )
    go_repository(
        name = "com_github_lann_ps",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/lann/ps",
        sum = "h1:P6pPBnrTSX3DEVR4fDembhRWSsG5rVo6hYhAB/ADZrk=",
        version = "v0.0.0-20150810152359-62de8c46ede0",
    )
    go_repository(
        name = "com_github_ldez_gomoddirectives",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ldez/gomoddirectives",
        sum = "h1:p9/sXuNFArS2RLc+UpYZSI4KQwGMEDWC/LbtF5OPFVg=",
        version = "v0.2.2",
    )
    go_repository(
        name = "com_github_ldez_tagliatelle",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ldez/tagliatelle",
        sum = "h1:693V8Bf1NdShJ8eu/s84QySA0J2VWBanVBa2WwXD/Wk=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_github_leodido_go_urn",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/leodido/go-urn",
        sum = "h1:BqpAaACuzVSgi/VLzGZIobT2z4v53pjosyNd9Yv6n/w=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_leonelquinteros_gotext",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/leonelquinteros/gotext",
        sum = "h1:ODY7LzLpZWWSJdAHnzhreOr6cwLXTAmc914FOauSkBM=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_letsencrypt_pkcs11key_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/letsencrypt/pkcs11key/v4",
        sum = "h1:qLc/OznH7xMr5ARJgkZCCWk+EomQkiNTOoOF5LAgagc=",
        version = "v4.0.0",
    )
    go_repository(
        name = "com_github_lib_pq",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/lib/pq",
        sum = "h1:jbk+ZieJ0D7EVGJYpL9QTz7/YW6UHbmdnZWYyK5cdBs=",
        version = "v1.10.6",
    )
    go_repository(
        name = "com_github_libgit2_git2go_v34",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/libgit2/git2go/v34",
        sum = "h1:UKoUaKLmiCRbOCD3PtUi2hD6hESSXzME/9OUZrGcgu8=",
        version = "v34.0.0",
    )

    go_repository(
        name = "com_github_liggitt_tabwriter",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/liggitt/tabwriter",
        sum = "h1:9TO3cAIGXtEhnIaL+V+BEER86oLrvS+kWobKpbJuye0=",
        version = "v0.0.0-20181228230101-89fcab3d43de",
    )
    go_repository(
        name = "com_github_lightstep_lightstep_tracer_common_golang_gogo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/lightstep/lightstep-tracer-common/golang/gogo",
        sum = "h1:YjW+hUb8Fh2S58z4av4t/0cBMK/Q0aP48RocCFsC8yI=",
        version = "v0.0.0-20210210170715-a8dfcb80d3a7",
    )
    go_repository(
        name = "com_github_lightstep_lightstep_tracer_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/lightstep/lightstep-tracer-go",
        sum = "h1:sGVnz8h3jTQuHKMbUe2949nXm3Sg09N1UcR3VoQNN5E=",
        version = "v0.25.0",
    )

    go_repository(
        name = "com_github_lithammer_dedent",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/lithammer/dedent",
        sum = "h1:VNzHMVCBNG1j0fh3OrsFRkVUwStdDArbgBWoPAffktY=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_logrusorgru_aurora",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/logrusorgru/aurora",
        sum = "h1:9MlwzLdW7QSDrhDjFlsEYmxpFyIoXmYRon3dt0io31k=",
        version = "v0.0.0-20181002194514-a7b3b318ed4e",
    )
    go_repository(
        name = "com_github_lufia_plan9stats",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/lufia/plan9stats",
        sum = "h1:6E+4a0GO5zZEnZ81pIr0yLvtUWk2if982qA3F3QD6H4=",
        version = "v0.0.0-20211012122336-39d0f177ccd0",
    )

    go_repository(
        name = "com_github_magiconair_properties",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/magiconair/properties",
        sum = "h1:b6kJs+EmPFMYGkow9GiUyCyOvIwYetYJ3fSaWak/Gls=",
        version = "v1.8.5",
    )
    go_repository(
        name = "com_github_mailgun_raymond_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mailgun/raymond/v2",
        sum = "h1:aOYHhvTpF5USySJ0o7cpPno/Uh2I5qg2115K25A+Ft4=",
        version = "v2.0.46",
    )

    go_repository(
        name = "com_github_mailru_easyjson",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mailru/easyjson",
        sum = "h1:UGYAvKxe3sBsEDzO8ZeWOSlIQfWFlxbzLZe7hwFURr0=",
        version = "v0.7.7",
    )
    go_repository(
        name = "com_github_makenowjust_heredoc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/MakeNowJust/heredoc",
        sum = "h1:cXCdzVdstXyiTqTvfqk9SDHpKNjxuom+DOlyEeQ4pzQ=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_maratori_testpackage",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/maratori/testpackage",
        sum = "h1:QtJ5ZjqapShm0w5DosRjg0PRlSdAdlx+W6cCKoALdbQ=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_markbates_errx",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/markbates/errx",
        sum = "h1:QDFeR+UP95dO12JgW+tgi2UVfo0V8YBHiUIOaeBPiEI=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_markbates_oncer",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/markbates/oncer",
        sum = "h1:E83IaVAHygyndzPimgUYJjbshhDTALZyXxvk9FOlQRY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_markbates_safe",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/markbates/safe",
        sum = "h1:yjZkbvRM6IzKj9tlu/zMJLS0n/V351OZWRnF3QfaUxI=",
        version = "v1.0.1",
    )

    go_repository(
        name = "com_github_masterminds_goutils",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Masterminds/goutils",
        sum = "h1:5nUrii3FMTL5diU80unEVvNevw1nH4+ZV4DSLVJLSYI=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_masterminds_semver",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Masterminds/semver",
        sum = "h1:H65muMkzWKEuNDnfl9d70GUjFniHKHRbFPGBuZ3QEww=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_github_masterminds_semver_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Masterminds/semver/v3",
        sum = "h1:hLg3sBzpNErnxhQtUy/mmLR2I9foDujNK030IGemrRc=",
        version = "v3.1.1",
    )
    go_repository(
        name = "com_github_masterminds_sprig",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Masterminds/sprig",
        sum = "h1:z4yfnGrZ7netVz+0EDJ0Wi+5VZCSYp4Z0m2dk6cEM60=",
        version = "v2.22.0+incompatible",
    )
    go_repository(
        name = "com_github_masterminds_sprig_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Masterminds/sprig/v3",
        sum = "h1:17jRggJu518dr3QaafizSXOjKYp94wKfABxUmyxvxX8=",
        version = "v3.2.2",
    )
    go_repository(
        name = "com_github_masterminds_squirrel",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Masterminds/squirrel",
        sum = "h1:YPpoceAcxuzIljlr5iWpNKaql7hLeG1KLSrhvdHpkZc=",
        version = "v1.5.3",
    )
    go_repository(
        name = "com_github_masterminds_vcs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Masterminds/vcs",
        sum = "h1:IIA2aBdXvfbIM+yl/eTnL4hb1XwdpvuQLglAix1gweE=",
        version = "v1.13.3",
    )
    go_repository(
        name = "com_github_matoous_godox",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/matoous/godox",
        sum = "h1:pWxk9e//NbPwfxat7RXkts09K+dEBJWakUWwICVqYbA=",
        version = "v0.0.0-20210227103229-6504466cf951",
    )
    go_repository(
        name = "com_github_matryer_is",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/matryer/is",
        sum = "h1:sosSmIWwkYITGrxZ25ULNDeKiMNzFSr4V/eqBQP0PeE=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_github_mattn_go_colorable",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mattn/go-colorable",
        sum = "h1:fFA4WZxdEF4tXPZVKMLwD8oUnCTTo08duU7wxecdEvA=",
        version = "v0.1.13",
    )
    go_repository(
        name = "com_github_mattn_go_ieproxy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mattn/go-ieproxy",
        sum = "h1:YkaHmK1CzE5C4O7A3hv3TCbfNDPSCf0RKZFX+VhBeYk=",
        version = "v0.0.3",
    )
    go_repository(
        name = "com_github_mattn_go_isatty",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mattn/go-isatty",
        sum = "h1:bq3VjFmv/sOjHtdEhmkEV4x1AJtvUvOJ2PFAZ5+peKQ=",
        version = "v0.0.16",
    )
    go_repository(
        name = "com_github_mattn_go_oci8",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mattn/go-oci8",
        sum = "h1:aEUDxNAyDG0tv8CA3TArnDQNyc4EhnWlsfxRgDHABHM=",
        version = "v0.1.1",
    )

    go_repository(
        name = "com_github_mattn_go_runewidth",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mattn/go-runewidth",
        sum = "h1:Lm995f3rfxdpd6TSmuVCHVb/QhupuXlYr8sCI/QdE+0=",
        version = "v0.0.9",
    )
    go_repository(
        name = "com_github_mattn_go_shellwords",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mattn/go-shellwords",
        sum = "h1:M2zGm7EW6UQJvDeQxo4T51eKPurbeFbe8WtebGE2xrk=",
        version = "v1.0.12",
    )
    go_repository(
        name = "com_github_mattn_go_sqlite3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mattn/go-sqlite3",
        sum = "h1:qZgc/Rwetq+MtyE18WhzjokPD93dNqLGNT3QJuLvBGw=",
        version = "v1.14.14",
    )
    go_repository(
        name = "com_github_mattn_goveralls",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mattn/goveralls",
        sum = "h1:7eJB6EqsPhRVxvwEXGnqdO2sJI0PTsrWoTMXEk9/OQc=",
        version = "v0.0.2",
    )

    go_repository(
        name = "com_github_matttproud_golang_protobuf_extensions",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/matttproud/golang_protobuf_extensions",
        sum = "h1:I0XW9+e1XWDxdcEniV4rQAIOPUGDq67JSCiRCgGCZLI=",
        version = "v1.0.2-0.20181231171920-c182affec369",
    )

    go_repository(
        name = "com_github_mbilski_exhaustivestruct",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mbilski/exhaustivestruct",
        sum = "h1:wCBmUnSYufAHO6J4AVWY6ff+oxWxsVFrwgOdMUQePUo=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_mgechev_dots",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mgechev/dots",
        sum = "h1:zpIH83+oKzcpryru8ceC6BxnoG8TBrhgAvRg8obzup0=",
        version = "v0.0.0-20210922191527-e955255bf517",
    )
    go_repository(
        name = "com_github_mgechev_revive",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mgechev/revive",
        sum = "h1:MiYA/o9M7REjvOF20QN43U8OtXDDHQFKLCtJnxLGLog=",
        version = "v1.1.2",
    )

    go_repository(
        name = "com_github_microcosm_cc_bluemonday",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/microcosm-cc/bluemonday",
        sum = "h1:dNH3e4PSyE4vNX+KlRGHT5KrSvjeUkoNPwEORjffHJg=",
        version = "v1.0.21",
    )
    go_repository(
        name = "com_github_microsoft_go_winio",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Microsoft/go-winio",
        sum = "h1:a9IhgEQBCUEk6QCdml9CiJGhAws+YwffDHEMp1VMrpA=",
        version = "v0.5.2",
    )
    go_repository(
        name = "com_github_microsoft_hcsshim",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Microsoft/hcsshim",
        sum = "h1:mnUj0ivWy6UzbB1uLFqKR6F+ZyiDc7j4iGgHTpO+5+I=",
        version = "v0.9.4",
    )

    go_repository(
        name = "com_github_miekg_dns",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/miekg/dns",
        sum = "h1:JKfpVSCB84vrAmHzyrsxB5NAr5kLoMXZArPSw7Qlgyg=",
        version = "v1.1.43",
    )
    go_repository(
        name = "com_github_miekg_pkcs11",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/miekg/pkcs11",
        sum = "h1:Ugu9pdy6vAYku5DEpVWVFPYnzV+bxB+iRdbuFSu7TvU=",
        version = "v1.1.1",
    )

    go_repository(
        name = "com_github_mistifyio_go_zfs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mistifyio/go-zfs",
        sum = "h1:aKW/4cBs+yK6gpqU3K/oIwk9Q/XICqd3zOX/UFuvqmk=",
        version = "v2.1.2-0.20190413222219-f784269be439+incompatible",
    )
    go_repository(
        name = "com_github_mitchellh_cli",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/cli",
        sum = "h1:qj8czE26AU4PbiaPXK5uVmMSM+V5BYsFBiM9HhGRLUA=",
        version = "v1.1.4",
    )
    go_repository(
        name = "com_github_mitchellh_copystructure",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/copystructure",
        sum = "h1:vpKXTN4ewci03Vljg/q9QvCGUDttBOGBIa15WveJJGw=",
        version = "v1.2.0",
    )

    go_repository(
        name = "com_github_mitchellh_go_homedir",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/go-homedir",
        sum = "h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_mitchellh_go_ps",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/go-ps",
        sum = "h1:i6ampVEEF4wQFF+bkYfwYgY+F/uYJDktmvLPf7qIgjc=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_mitchellh_go_testing_interface",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/go-testing-interface",
        sum = "h1:fzU/JVNcaqHQEcVFAKeR41fkiLdIPrefOvVG1VZ96U0=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_mitchellh_go_wordwrap",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/go-wordwrap",
        sum = "h1:6GlHJ/LTGMrIJbwgdqdl2eEH8o+Exx/0m8ir9Gns0u4=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_mitchellh_gox",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/gox",
        sum = "h1:lfGJxY7ToLJQjHHwi0EX6uYBdK78egf954SQl13PQJc=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_github_mitchellh_iochan",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/iochan",
        sum = "h1:C+X3KsSTLFVBr/tK1eYN/vs4rJcvsiLU338UhYPJWeY=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_mitchellh_mapstructure",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/mapstructure",
        sum = "h1:6h7AQ0yhTcIsmFmnAwQls75jp2Gzs4iB8W7pjMO+rqo=",
        version = "v1.4.2",
    )

    go_repository(
        name = "com_github_mitchellh_reflectwalk",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mitchellh/reflectwalk",
        sum = "h1:G2LzWKi524PWgd3mLHV8Y5k7s6XUvT0Gef6zxSIeXaQ=",
        version = "v1.0.2",
    )

    go_repository(
        name = "com_github_moby_locker",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/moby/locker",
        sum = "h1:fOXqR41zeveg4fFODix+1Ch4mj/gT0NE1XJbp/epuBg=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_moby_spdystream",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/moby/spdystream",
        sum = "h1:cjW1zVyyoiM0T7b6UoySUFqzXMoqRckQtXwGPiBhOM8=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_moby_sys_mountinfo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/moby/sys/mountinfo",
        sum = "h1:2Ks8/r6lopsxWi9m58nlwjaeSzUX9iiL1vj5qB/9ObI=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_moby_sys_signal",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/moby/sys/signal",
        sum = "h1:aDpY94H8VlhTGa9sNYUFCFsMZIUh5wm0B6XkIoJj/iY=",
        version = "v0.6.0",
    )

    go_repository(
        name = "com_github_moby_sys_symlink",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/moby/sys/symlink",
        sum = "h1:tk1rOM+Ljp0nFmfOIBtlV3rTDlWOwFRhjEeAhZB0nZc=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_moby_term",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/moby/term",
        sum = "h1:dcztxKSvZ4Id8iPpHERQBbIJfabdt4wUm5qy3wOL2Zc=",
        version = "v0.0.0-20210619224110-3f7ff695adc6",
    )
    go_repository(
        name = "com_github_modern_go_concurrent",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/modern-go/concurrent",
        sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
        version = "v0.0.0-20180306012644-bacd9c7ef1dd",
    )
    go_repository(
        name = "com_github_modern_go_reflect2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/modern-go/reflect2",
        sum = "h1:xBagoLtFs94CBntxluKeaWgTMpvLxC4ur3nMaC9Gz0M=",
        version = "v1.0.2",
    )
    go_repository(
        name = "com_github_mohae_deepcopy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mohae/deepcopy",
        sum = "h1:RWengNIwukTxcDr9M+97sNutRR1RKhG96O6jWumTTnw=",
        version = "v0.0.0-20170929034955-c48cc78d4826",
    )

    go_repository(
        name = "com_github_monochromegane_go_gitignore",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/monochromegane/go-gitignore",
        sum = "h1:n6/2gBQ3RWajuToeY6ZtZTIKv2v7ThUy5KKusIT0yc0=",
        version = "v0.0.0-20200626010858-205db1a8cc00",
    )
    go_repository(
        name = "com_github_montanaflynn_stats",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/montanaflynn/stats",
        sum = "h1:bsjNADsjHq0gjU7KO7zwoX5k3HtFdf6TDzB3ncl5iUs=",
        version = "v0.0.0-20151014174947-eeaced052adb",
    )
    go_repository(
        name = "com_github_moricho_tparallel",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/moricho/tparallel",
        sum = "h1:95FytivzT6rYzdJLdtfn6m1bfFJylOJK41+lgv/EHf4=",
        version = "v0.2.1",
    )

    go_repository(
        name = "com_github_morikuni_aec",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/morikuni/aec",
        sum = "h1:nP9CBfwrvYnBRgY6qfDQkygYDmYwOilePFkwzv4dU8A=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_mozilla_scribe",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mozilla/scribe",
        sum = "h1:29NKShH4TWd3lxCDUhS4Xe16EWMA753dtIxYtwddklU=",
        version = "v0.0.0-20180711195314-fb71baf557c1",
    )
    go_repository(
        name = "com_github_mozilla_tls_observatory",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mozilla/tls-observatory",
        sum = "h1:0KqC6/sLy7fDpBdybhVkkv4Yz+PmB7c9Dz9z3dLW804=",
        version = "v0.0.0-20210609171429-7bc42856d2e5",
    )

    go_repository(
        name = "com_github_munnerz_goautoneg",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/munnerz/goautoneg",
        sum = "h1:C3w9PqII01/Oq1c1nUAm88MOHcQC9l5mIlSMApZMrHA=",
        version = "v0.0.0-20191010083416-a7dc8b61c822",
    )
    go_repository(
        name = "com_github_mwitkow_go_conntrack",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mwitkow/go-conntrack",
        sum = "h1:KUppIJq7/+SVif2QVs3tOP0zanoHgBEVAwHxUSIzRqU=",
        version = "v0.0.0-20190716064945-2f068394615f",
    )
    go_repository(
        name = "com_github_mwitkow_go_proto_validators",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mwitkow/go-proto-validators",
        sum = "h1:F6LFfmgVnfULfaRsQWBbe7F7ocuHCr9+7m+GAeDzNbQ=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_github_mxk_go_flowrate",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/mxk/go-flowrate",
        sum = "h1:y5//uYreIhSUg3J1GEMiLbxo1LJaP8RfCpH6pymGZus=",
        version = "v0.0.0-20140419014527-cca7078d478f",
    )
    go_repository(
        name = "com_github_nakabonne_nestif",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/nakabonne/nestif",
        sum = "h1:wm28nZjhQY5HyYPx+weN3Q65k6ilSBxDb8v5S81B81U=",
        version = "v0.3.1",
    )
    go_repository(
        name = "com_github_nbutton23_zxcvbn_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/nbutton23/zxcvbn-go",
        sum = "h1:4kuARK6Y6FxaNu/BnU2OAaLF86eTVhP2hjTB6iMvItA=",
        version = "v0.0.0-20210217022336-fa2cb2858354",
    )

    go_repository(
        name = "com_github_niemeyer_pretty",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/niemeyer/pretty",
        sum = "h1:fD57ERR4JtEqsWbfPhv4DMiApHyliiK5xCTNVSPiaAs=",
        version = "v0.0.0-20200227124842-a10e7caefd8e",
    )
    go_repository(
        name = "com_github_nishanths_exhaustive",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/nishanths/exhaustive",
        sum = "h1:+ANTMqRNrqwInnP9aszg/0jDo+zbXa4x66U19Bx/oTk=",
        version = "v0.2.3",
    )
    go_repository(
        name = "com_github_nishanths_predeclared",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/nishanths/predeclared",
        sum = "h1:1TXtjmy4f3YCFjTxRd8zcFHOmoUir+gp0ESzjFzG2sw=",
        version = "v0.2.1",
    )

    go_repository(
        name = "com_github_nxadm_tail",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/nxadm/tail",
        sum = "h1:nPr65rt6Y5JFSKQO7qToXr7pePgD6Gwiw05lkbyAQTE=",
        version = "v1.4.8",
    )
    go_repository(
        name = "com_github_nytimes_gziphandler",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/NYTimes/gziphandler",
        sum = "h1:ZUDjpQae29j0ryrS0u/B8HZfJBtBQHjqw2rQ2cqUQ3I=",
        version = "v1.1.1",
    )

    go_repository(
        name = "com_github_oklog_ulid_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/oklog/ulid/v2",
        sum = "h1:r4fFzBm+bv0wNKNh5eXTwU7i85y5x+uwkxCUTNVQqLc=",
        version = "v2.0.2",
    )
    go_repository(
        name = "com_github_olekukonko_tablewriter",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/olekukonko/tablewriter",
        sum = "h1:P2Ga83D34wi1o9J6Wh1mRuqd4mF/x/lgBS7N7AbDhec=",
        version = "v0.0.5",
    )
    go_repository(
        name = "com_github_olekukonko_ts",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/olekukonko/ts",
        sum = "h1:LiZB1h0GIcudcDci2bxbqI6DXV8bF8POAnArqvRrIyw=",
        version = "v0.0.0-20171002115256-78ecb04241c0",
    )
    go_repository(
        name = "com_github_oneofone_xxhash",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/OneOfOne/xxhash",
        sum = "h1:31czK/TI9sNkxIKfaUfGlU47BAxQ0ztGgd9vPyqimf8=",
        version = "v1.2.8",
    )
    go_repository(
        name = "com_github_onsi_ginkgo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/onsi/ginkgo",
        sum = "h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=",
        version = "v1.16.5",
    )
    go_repository(
        name = "com_github_onsi_ginkgo_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/onsi/ginkgo/v2",
        sum = "h1:3ZNA3L1c5FYDFTTxbFeVGGD8jYvjYauHD30YgLxVsNI=",
        version = "v2.2.0",
    )

    go_repository(
        name = "com_github_onsi_gomega",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/onsi/gomega",
        sum = "h1:8uQq0zMgLEfa0vRrrBgaJF2gyW9Da9BmfGV+OyUzfkY=",
        version = "v1.20.2",
    )

    go_repository(
        name = "com_github_open_policy_agent_opa",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/open-policy-agent/opa",
        sum = "h1:sEZthsrWBqIN+ShTMJ0Hcz6a3GkYsY4FaB2S/ou2hZk=",
        version = "v0.44.0",
    )

    go_repository(
        name = "com_github_opencontainers_go_digest",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/opencontainers/go-digest",
        sum = "h1:apOUWs51W5PlhuyGyz9FCeeBIOUDA/6nW8Oi/yOhh5U=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_opencontainers_image_spec",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/opencontainers/image-spec",
        sum = "h1:+czc/J8SlhPKLOtVLMQc+xDCFBT73ZStMsRhSsUhsSg=",
        version = "v1.0.3-0.20220114050600-8b9d41f48198",
    )
    go_repository(
        name = "com_github_opencontainers_runc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/opencontainers/runc",
        sum = "h1:2VSZwLx5k/BfsBxMMipG/LYUnmqOD/BPkIVgQUcTlLw=",
        version = "v1.1.2",
    )
    go_repository(
        name = "com_github_opencontainers_runtime_spec",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/opencontainers/runtime-spec",
        sum = "h1:3snG66yBm59tKhhSPQrQ/0bCrv1LQbKt40LnUPiUxdc=",
        version = "v1.0.3-0.20210326190908-1c3f411f0417",
    )

    go_repository(
        name = "com_github_opencontainers_selinux",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/opencontainers/selinux",
        sum = "h1:09LIPVRP3uuZGQvgR+SgMSNBd1Eb3vlRbGqQpoHsF8w=",
        version = "v1.10.1",
    )
    go_repository(
        name = "com_github_openpeedeep_depguard",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/OpenPeeDeeP/depguard",
        sum = "h1:VlW4R6jmBIv3/u1JNlawEvJMM4J+dPORPaZasQee8Us=",
        version = "v1.0.1",
    )

    go_repository(
        name = "com_github_opentracing_opentracing_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/opentracing/opentracing-go",
        sum = "h1:uEJPy/1a5RIPAJ0Ov+OIO8OxWu77jEv+1B0VhjKrZUs=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_otiai10_copy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/otiai10/copy",
        sum = "h1:HvG945u96iNadPoG2/Ja2+AUJeW5YuFQMixq9yirC+k=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_otiai10_curr",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/otiai10/curr",
        sum = "h1:TJIWdbX0B+kpNagQrjgq8bCMrbhiuX73M2XwgtDMoOI=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_otiai10_mint",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/otiai10/mint",
        sum = "h1:BCmzIS3n71sGfHB5NMNDB3lHYPz8fWSkCAErHed//qc=",
        version = "v1.3.1",
    )

    go_repository(
        name = "com_github_pascaldekloe_goe",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pascaldekloe/goe",
        sum = "h1:Lgl0gzECD8GnQ5QCWA8o6BtfL6mDH5rQgM4/fX3avOs=",
        version = "v0.0.0-20180627143212-57f6aae5913c",
    )
    go_repository(
        name = "com_github_pborman_getopt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pborman/getopt",
        sum = "h1:BHT1/DKsYDGkUgQ2jmMaozVcdk+sVfz0+1ZJq4zkWgw=",
        version = "v0.0.0-20170112200414-7148bc3a4c30",
    )
    go_repository(
        name = "com_github_pborman_uuid",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pborman/uuid",
        sum = "h1:J7Q5mO4ysT1dv8hyrUGHb9+ooztCXu1D8MY8DZYsu3g=",
        version = "v1.2.0",
    )

    go_repository(
        name = "com_github_pelletier_go_toml",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pelletier/go-toml",
        sum = "h1:4yBQzkHv+7BHq2PQUZF3Mx0IYxG7LsP222s7Agd3ve8=",
        version = "v1.9.5",
    )
    go_repository(
        name = "com_github_pelletier_go_toml_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pelletier/go-toml/v2",
        sum = "h1:ipoSadvV8oGUjnUbMub59IDPPwfxF694nG/jwbMiyQg=",
        version = "v2.0.5",
    )

    go_repository(
        name = "com_github_peterbourgon_diskv",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/peterbourgon/diskv",
        sum = "h1:UBdAOUP5p4RWqPBg048CAvpKN+vxiaj6gdUUzhl4XmI=",
        version = "v2.0.1+incompatible",
    )
    go_repository(
        name = "com_github_peterh_liner",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/peterh/liner",
        sum = "h1:zapSxdmZYY6vJWXFKLQ+MkI+agc+HQyfrCGowDSHiKs=",
        version = "v0.0.0-20170211195444-bf27d3ba8e1d",
    )
    go_repository(
        name = "com_github_phayes_checkstyle",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/phayes/checkstyle",
        sum = "h1:CdDQnGF8Nq9ocOS/xlSptM1N3BbrA6/kmaep5ggwaIA=",
        version = "v0.0.0-20170904204023-bfd46e6a821d",
    )

    go_repository(
        name = "com_github_phayes_freeport",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/phayes/freeport",
        sum = "h1:Ii+DKncOVM8Cu1Hc+ETb5K+23HdAMvESYE3ZJ5b5cMI=",
        version = "v0.0.0-20220201140144-74d24b5ae9f5",
    )

    go_repository(
        name = "com_github_philhofer_fwd",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/philhofer/fwd",
        sum = "h1:GdGcTjf5RNAxwS4QLsiMzJYj5KEvPJD3Abr261yRQXQ=",
        version = "v1.1.1",
    )

    go_repository(
        name = "com_github_pingcap_errors",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pingcap/errors",
        sum = "h1:lFuQV/oaUMGcD2tqt+01ROSmJs75VG1ToEOkZIZ4nE4=",
        version = "v0.11.4",
    )
    go_repository(
        name = "com_github_piotrkowalczuk_promgrpc_v4",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/piotrkowalczuk/promgrpc/v4",
        sum = "h1:X8NIOHEeAb578lV/VVb4PeAMewtmLi0hfpy1pWlKsXQ=",
        version = "v4.0.4",
    )
    go_repository(
        name = "com_github_pkg_diff",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pkg/diff",
        sum = "h1:aoZm08cpOy4WuID//EZDgcC4zIxODThtZNPirFr42+A=",
        version = "v0.0.0-20210226163009-20ebb0f2a09e",
    )

    go_repository(
        name = "com_github_pkg_errors",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pkg/errors",
        sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
        version = "v0.9.1",
    )

    go_repository(
        name = "com_github_pkg_sftp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pkg/sftp",
        sum = "h1:I2qBYMChEhIjOgazfJmV3/mZM256btk6wkCDRmW7JYs=",
        version = "v1.13.1",
    )
    go_repository(
        name = "com_github_pmezard_go_difflib",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pmezard/go-difflib",
        sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_polyfloyd_go_errorlint",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/polyfloyd/go-errorlint",
        sum = "h1:Kq/3kL0k033ds3tyez5lFPrfQ74fNJ+OqCclRipubwA=",
        version = "v0.0.0-20210722154253-910bb7978349",
    )

    go_repository(
        name = "com_github_posener_complete",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/posener/complete",
        sum = "h1:NP0eAhjcjImqslEwo/1hq7gpajME0fTLTezBKDqfXqo=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_poy_onpar",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/poy/onpar",
        sum = "h1:oL4IBbcqwhhNWh31bjOX8C/OCy0zs9906d/VUru+bqg=",
        version = "v0.0.0-20190519213022-ee068f8ea4d1",
    )

    go_repository(
        name = "com_github_pquerna_cachecontrol",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pquerna/cachecontrol",
        sum = "h1:yJMy84ti9h/+OEWa752kBTKv4XC30OtVVHYv/8cTqKc=",
        version = "v0.1.0",
    )

    go_repository(
        name = "com_github_prometheus_client_golang",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/prometheus/client_golang",
        sum = "h1:3gMjIY2+/hzmqhtUC/aQNYldJA6DtH3CgQvwS+02K1c=",
        version = "v1.13.1",
    )
    go_repository(
        name = "com_github_prometheus_client_model",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/prometheus/client_model",
        sum = "h1:uq5h0d+GuxiXLJLNABMgp2qUWDPiLvgCzz2dUR+/W/M=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_github_prometheus_common",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/prometheus/common",
        sum = "h1:ccBbHCgIiT9uSoFY0vX8H3zsNR5eLt17/RQLUvn8pXE=",
        version = "v0.37.0",
    )
    go_repository(
        name = "com_github_prometheus_procfs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/prometheus/procfs",
        sum = "h1:ODq8ZFEaYeCaZOJlZZdJA2AbQR98dSHSM1KW/You5mo=",
        version = "v0.8.0",
    )

    go_repository(
        name = "com_github_protonmail_go_crypto",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ProtonMail/go-crypto",
        sum = "h1:QeeqI2zxxgZVe11UrYFXXx6gVxPVF40ygekjBzEg4XY=",
        version = "v0.0.0-20220930113650-c6815a8c17ad",
    )

    go_repository(
        name = "com_github_pseudomuto_protoc_gen_doc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pseudomuto/protoc-gen-doc",
        sum = "h1:61vWZuxYa8D7Rn4h+2dgoTNqnluBmJya2MgbqO32z6g=",
        version = "v1.3.2",
    )
    go_repository(
        name = "com_github_pseudomuto_protokit",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/pseudomuto/protokit",
        sum = "h1:hlnBDcy3YEDXH7kc9gV+NLaN0cDzhDvD1s7Y6FZ8RpM=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_github_puerkitobio_purell",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/PuerkitoBio/purell",
        sum = "h1:WEQqlqaGbrPkxLJWfBwQmfEAE1Z7ONdDLqrN38tNFfI=",
        version = "v1.1.1",
    )
    go_repository(
        name = "com_github_puerkitobio_urlesc",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/PuerkitoBio/urlesc",
        sum = "h1:d+Bc7a5rLufV/sSk/8dngufqelfh6jnri85riMAaF/M=",
        version = "v0.0.0-20170810143723-de5bf2ad4578",
    )
    go_repository(
        name = "com_github_quasilyte_go_consistent",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/quasilyte/go-consistent",
        sum = "h1:JoUA0uz9U0FVFq5p4LjEq4C0VgQ0El320s3Ms0V4eww=",
        version = "v0.0.0-20190521200055-c6f3937de18c",
    )
    go_repository(
        name = "com_github_quasilyte_go_ruleguard",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/quasilyte/go-ruleguard",
        sum = "h1:O1G41cq1jUr3cJmqp7vOUT0SokqjzmS9aESWJuIDRaY=",
        version = "v0.3.13",
    )
    go_repository(
        name = "com_github_quasilyte_go_ruleguard_dsl",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/quasilyte/go-ruleguard/dsl",
        sum = "h1:4tVlVVcBT+nNWoF+t/zrAMO13sHAqYotX1K12Gc8f8A=",
        version = "v0.3.10",
    )
    go_repository(
        name = "com_github_quasilyte_go_ruleguard_rules",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/quasilyte/go-ruleguard/rules",
        sum = "h1:cRLFDAB53r5wIkxYvtQUMnn3+B09uZTAOPmefNfVk5I=",
        version = "v0.0.0-20210428214800-545e0d2e0bf7",
    )
    go_repository(
        name = "com_github_quasilyte_regex_syntax",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/quasilyte/regex/syntax",
        sum = "h1:L8QM9bvf68pVdQ3bCFZMDmnt9yqcMBro1pC7F+IPYMY=",
        version = "v0.0.0-20200407221936-30656e2c4a95",
    )

    go_repository(
        name = "com_github_rcrowley_go_metrics",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/rcrowley/go-metrics",
        sum = "h1:N/ElC8H3+5XpJzTSTfLsJV/mx9Q9g7kxmchpfZyxgzM=",
        version = "v0.0.0-20201227073835-cf1acfcdf475",
    )
    go_repository(
        name = "com_github_robfig_cron_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/robfig/cron/v3",
        sum = "h1:WdRxkvbJztn8LMz/QEvLN5sBU+xKpSqwwUO1Pjr4qDs=",
        version = "v3.0.1",
    )

    go_repository(
        name = "com_github_rogpeppe_fastuuid",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/rogpeppe/fastuuid",
        sum = "h1:Ppwyp6VYCF1nvBTXL3trRso7mXMlRrw9ooo375wvi2s=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_rogpeppe_go_internal",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/rogpeppe/go-internal",
        sum = "h1:FCbCCtXNOY3UtUuHUYaghJg4y7Fd14rXifAYUAtL9R8=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_github_rs_cors",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/rs/cors",
        sum = "h1:+88SsELBHx5r+hZ8TCkggzSstaWNbDvThkVK8H6f9ik=",
        version = "v1.7.0",
    )

    go_repository(
        name = "com_github_rubenv_sql_migrate",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/rubenv/sql-migrate",
        sum = "h1:fOXMPLMd41sK7Tg75SXDec15k3zg5WNV6SjuDRiNfcU=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_rubyist_tracerx",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/rubyist/tracerx",
        sum = "h1:mncRSDOqYCng7jOD+Y6+IivdRI6Kzv2BLWYkWkdQfu0=",
        version = "v0.0.0-20170927163412-787959303086",
    )
    go_repository(
        name = "com_github_russross_blackfriday",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/russross/blackfriday",
        sum = "h1:KqfZb0pUVN2lYqZUYRddxF4OR8ZMURnJIG5Y3VRLtww=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_github_russross_blackfriday_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/russross/blackfriday/v2",
        sum = "h1:JIOH55/0cWyOuilr9/qlrm0BSXldqnqwMsf35Ld67mk=",
        version = "v2.1.0",
    )
    go_repository(
        name = "com_github_ryancurrah_gomodguard",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ryancurrah/gomodguard",
        sum = "h1:ww2fsjqocGCAFamzvv/b8IsRduuHHeK2MHTcTxZTQX8=",
        version = "v1.2.3",
    )
    go_repository(
        name = "com_github_ryanrolds_sqlclosecheck",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ryanrolds/sqlclosecheck",
        sum = "h1:AZx+Bixh8zdUBxUA1NxbxVAS78vTPq4rCb8OUZI9xFw=",
        version = "v0.3.0",
    )

    go_repository(
        name = "com_github_ryanuber_columnize",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ryanuber/columnize",
        sum = "h1:UFr9zpz4xgTnIE5yIMtWAMngCdZ9p/+q6lTbgelo80M=",
        version = "v0.0.0-20160712163229-9b3edd62028f",
    )

    go_repository(
        name = "com_github_sagikazarmark_crypt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sagikazarmark/crypt",
        sum = "h1:AyO7PGna28P9TMH93Bsxd7m9QC4xE6zyGQTXCo7ZrA8=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_sanposhiho_wastedassign_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sanposhiho/wastedassign/v2",
        sum = "h1:+6/hQIHKNJAUixEj6EmOngGIisyeI+T3335lYTyxRoA=",
        version = "v2.0.6",
    )

    go_repository(
        name = "com_github_schollz_closestmatch",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/schollz/closestmatch",
        sum = "h1:Uel2GXEpJqOWBrlyI+oY9LTiyyjYS17cCYRqP13/SHk=",
        version = "v2.1.0+incompatible",
    )

    go_repository(
        name = "com_github_sean_seed",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sean-/seed",
        sum = "h1:nn5Wsu0esKSJiIVhscUtVbo7ada43DJhG55ua/hjS5I=",
        version = "v0.0.0-20170313163322-e2103e2c3529",
    )
    go_repository(
        name = "com_github_sebest_xff",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sebest/xff",
        sum = "h1:iLcLb5Fwwz7g/DLK89F+uQBDeAhHhwdzB5fSlVdhGcM=",
        version = "v0.0.0-20210106013422-671bd2870b3a",
    )

    go_repository(
        name = "com_github_securego_gosec_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/securego/gosec/v2",
        sum = "h1:anHKLS/ApTYU6NZkKa/5cQqqcbKZURjvc+MtR++S4EQ=",
        version = "v2.9.1",
    )

    go_repository(
        name = "com_github_sergi_go_diff",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sergi/go-diff",
        sum = "h1:we8PVUC3FE2uYfodKH/nBHMSetSfHDR6scGdBi+erh0=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_shazow_go_diff",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/shazow/go-diff",
        sum = "h1:W65qqJCIOVP4jpqPQ0YvHYKwcMEMVWIzWC5iNQQfBTU=",
        version = "v0.0.0-20160112020656-b6b7b6733b8c",
    )

    go_repository(
        name = "com_github_shirou_gopsutil_v3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/shirou/gopsutil/v3",
        sum = "h1:flTg1DrnV/UVrBqjLgVgDJzx6lf+91rC64/dBHmO2IA=",
        version = "v3.21.10",
    )
    go_repository(
        name = "com_github_shogo82148_go_shuffle",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/shogo82148/go-shuffle",
        sum = "h1:VI1u+o2KZPZ5AhuPpXY0JBdpQPnkTx6Dd5XJhK/9MYE=",
        version = "v0.0.0-20170808115208-59829097ff3b",
    )
    go_repository(
        name = "com_github_shopify_goreferrer",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Shopify/goreferrer",
        sum = "h1:KkH3I3sJuOLP3TjA/dfr4NAY8bghDwnXiU7cTKxQqo0=",
        version = "v0.0.0-20220729165902-8cddb4f5de06",
    )
    go_repository(
        name = "com_github_shopify_logrus_bugsnag",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/Shopify/logrus-bugsnag",
        sum = "h1:UrqY+r/OJnIp5u0s1SbQ8dVfLCZJsnvazdBP5hS4iRs=",
        version = "v0.0.0-20171204204709-577dee27f20d",
    )
    go_repository(
        name = "com_github_shopspring_decimal",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/shopspring/decimal",
        sum = "h1:abSATXmQEYyShuxI4/vyW3tV1MrKAJzCZ/0zLUXYbsQ=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_shurcool_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/shurcooL/go",
        sum = "h1:MZM7FHLqUHYI0Y/mQAt3d2aYa0SiNms/hFqC9qJYolM=",
        version = "v0.0.0-20180423040247-9e1955d9fb6e",
    )
    go_repository(
        name = "com_github_shurcool_go_goon",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/shurcooL/go-goon",
        sum = "h1:llrF3Fs4018ePo4+G/HV/uQUqEI1HMDjCeOf2V6puPc=",
        version = "v0.0.0-20170922171312-37c2f522c041",
    )

    go_repository(
        name = "com_github_shurcool_sanitized_anchor_name",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/shurcooL/sanitized_anchor_name",
        sum = "h1:PdmoCO6wvbs+7yrJyMORt4/BmY5IYyJwS/kOiWx8mHo=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_sirupsen_logrus",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sirupsen/logrus",
        sum = "h1:trlNQbNUG3OdDrDil03MCb1H2o9nJ1x4/5LYw7byDE0=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_github_sivchari_tenv",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sivchari/tenv",
        sum = "h1:FdTpgRlTue5eb5nXIYgS/lyVXSjugU8UUVDwhP1NLU8=",
        version = "v1.4.7",
    )

    go_repository(
        name = "com_github_smartystreets_assertions",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/smartystreets/assertions",
        sum = "h1:zE9ykElWQ6/NYmHa3jpm/yHnI4xSofP+UP6SpjHcSeM=",
        version = "v0.0.0-20180927180507-b2de0cb4f26d",
    )
    go_repository(
        name = "com_github_smartystreets_goconvey",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/smartystreets/goconvey",
        sum = "h1:fv0U8FUIMPNf1L9lnHLvLhgicrIVChEkdzIKYqbNC9s=",
        version = "v1.6.4",
    )
    go_repository(
        name = "com_github_soheilhy_cmux",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/soheilhy/cmux",
        sum = "h1:jjzc5WVemNEDTLwv9tlmemhC73tI08BNOIGwBOo10Js=",
        version = "v0.1.5",
    )
    go_repository(
        name = "com_github_sonatard_noctx",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sonatard/noctx",
        sum = "h1:VC1Qhl6Oxx9vvWo3UDgrGXYCeKCe3Wbw7qAWL6FrmTY=",
        version = "v0.0.1",
    )
    go_repository(
        name = "com_github_sourcegraph_go_diff",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sourcegraph/go-diff",
        sum = "h1:hmA1LzxW0n1c3Q4YbrFgg4P99GSnebYa3x8gr0HZqLQ=",
        version = "v0.6.1",
    )

    go_repository(
        name = "com_github_spaolacci_murmur3",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/spaolacci/murmur3",
        sum = "h1:qLC7fQah7D6K1B0ujays3HV9gkFtllcxhzImRR7ArPQ=",
        version = "v0.0.0-20180118202830-f09979ecbc72",
    )
    go_repository(
        name = "com_github_spf13_afero",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/spf13/afero",
        sum = "h1:j49Hj62F0n+DaZ1dDCvhABaPNSGNkt32oRFxI33IEMw=",
        version = "v1.9.2",
    )
    go_repository(
        name = "com_github_spf13_cast",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/spf13/cast",
        sum = "h1:s0hze+J0196ZfEMTs80N7UlFt0BDuQ7Q+JDnHiMWKdA=",
        version = "v1.4.1",
    )

    go_repository(
        name = "com_github_spf13_cobra",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/spf13/cobra",
        sum = "h1:o94oiPyS4KD1mPy2fmcYYHHfCxLqYjJOhGsCHFZtEzA=",
        version = "v1.6.1",
    )
    go_repository(
        name = "com_github_spf13_jwalterweatherman",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/spf13/jwalterweatherman",
        sum = "h1:ue6voC5bR5F8YxI5S67j9i582FU4Qvo2bmqnqMYADFk=",
        version = "v1.1.0",
    )

    go_repository(
        name = "com_github_spf13_pflag",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/spf13/pflag",
        sum = "h1:iy+VFUOCP1a+8yFto/drg2CJ5u0yRoB7fZw3DKv/JXA=",
        version = "v1.0.5",
    )
    go_repository(
        name = "com_github_spf13_viper",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/spf13/viper",
        sum = "h1:yR6EXjTp0y0cLN8OZg1CRZmOBdI88UcGkhgyJhu6nZk=",
        version = "v1.9.0",
    )

    go_repository(
        name = "com_github_spyzhov_ajson",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/spyzhov/ajson",
        sum = "h1:1MDIlPc6x0zjNtpa7tDzRAyFAvRX+X8ZsvtYz5lZg6A=",
        version = "v0.7.1",
    )
    go_repository(
        name = "com_github_ssgelm_cookiejarparser",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ssgelm/cookiejarparser",
        sum = "h1:cRdXauUbOTFzTPJFaeiWbHnQ+tRGlpKKzvIK9PUekE4=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_ssgreg_nlreturn_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ssgreg/nlreturn/v2",
        sum = "h1:X4XDI7jstt3ySqGU86YGAURbxw3oTDPK9sPEi6YEwQ0=",
        version = "v2.2.1",
    )

    go_repository(
        name = "com_github_stackexchange_wmi",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/StackExchange/wmi",
        sum = "h1:VIkavFPXSjcnS+O8yTq7NI32k0R5Aj+v39y29VYDOSA=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_stefanberger_go_pkcs11uri",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/stefanberger/go-pkcs11uri",
        sum = "h1:lIOOHPEbXzO3vnmx2gok1Tfs31Q8GQqKLc8vVqyQq/I=",
        version = "v0.0.0-20201008174630-78d3cae3a980",
    )
    go_repository(
        name = "com_github_stoewer_go_strcase",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/stoewer/go-strcase",
        sum = "h1:Z2iHWqGXH00XYgqDmNgQbIBxf3wrNq0F3feEy0ainaU=",
        version = "v1.2.0",
    )

    go_repository(
        name = "com_github_stretchr_objx",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/stretchr/objx",
        sum = "h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=",
        version = "v0.5.0",
    )
    go_repository(
        name = "com_github_stretchr_testify",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/stretchr/testify",
        sum = "h1:w7B6lhMri9wdJUVmEZPGGhZzrYTPvgJArz7wNPgYKsk=",
        version = "v1.8.1",
    )
    go_repository(
        name = "com_github_subosito_gotenv",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/subosito/gotenv",
        sum = "h1:Slr1R9HxAlEKefgq5jn9U+DnETlIUa6HfgEzj0g5d7s=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_sylvia7788_contextcheck",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/sylvia7788/contextcheck",
        sum = "h1:MsiVqROAdr0efZc/fOCt0c235qm9XJqHtWwM+2h2B04=",
        version = "v1.0.4",
    )

    go_repository(
        name = "com_github_tchap_go_patricia",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tchap/go-patricia",
        sum = "h1:JvoDL7JSoIP2HDE8AbDH3zC8QBPxmzYe32HHy5yQ+Ck=",
        version = "v2.2.6+incompatible",
    )
    go_repository(
        name = "com_github_tchap_go_patricia_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tchap/go-patricia/v2",
        sum = "h1:6rQp39lgIYZ+MHmdEq4xzuk1t7OdC35z/xm0BGhTkes=",
        version = "v2.3.1",
    )

    go_repository(
        name = "com_github_tdakkota_asciicheck",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tdakkota/asciicheck",
        sum = "h1:HxLVTlqcHhFAz3nWUcuvpH7WuOMv8LQoCWmruLfFH2U=",
        version = "v0.0.0-20200416200610-e657995f937b",
    )
    go_repository(
        name = "com_github_tdewolff_minify_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tdewolff/minify/v2",
        sum = "h1:kejsHQMM17n6/gwdw53qsi6lg0TGddZADVyQOz1KMdE=",
        version = "v2.12.4",
    )
    go_repository(
        name = "com_github_tdewolff_parse_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tdewolff/parse/v2",
        sum = "h1:KCkDvNUMof10e3QExio9OPZJT8SbdKojLBumw8YZycQ=",
        version = "v2.6.4",
    )

    go_repository(
        name = "com_github_tenntenn_modver",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tenntenn/modver",
        sum = "h1:2klLppGhDgzJrScMpkj9Ujy3rXPUspSjAcev9tSEBgA=",
        version = "v1.0.1",
    )
    go_repository(
        name = "com_github_tenntenn_text_transform",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tenntenn/text/transform",
        sum = "h1:f+jULpRQGxTSkNYKJ51yaw6ChIqO+Je8UqsTKN/cDag=",
        version = "v0.0.0-20200319021203-7eef512accb3",
    )
    go_repository(
        name = "com_github_tetafro_godot",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tetafro/godot",
        sum = "h1:BVoBIqAf/2QdbFmSwAWnaIqDivZdOV0ZRwEm6jivLKw=",
        version = "v1.4.11",
    )
    go_repository(
        name = "com_github_timakin_bodyclose",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/timakin/bodyclose",
        sum = "h1:ig99OeTyDwQWhPe2iw9lwfQVF1KB3Q4fpP3X7/2VBG8=",
        version = "v0.0.0-20200424151742-cb6215831a94",
    )

    go_repository(
        name = "com_github_tinylib_msgp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tinylib/msgp",
        sum = "h1:gWmO7n0Ys2RBEb7GPYB9Ujq8Mk5p2U08lRnmMcGy6BQ=",
        version = "v1.1.2",
    )
    go_repository(
        name = "com_github_tklauser_go_sysconf",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tklauser/go-sysconf",
        sum = "h1:JeUVdAOWhhxVcU6Eqr/ATFHgXk/mmiItdKeJPev3vTo=",
        version = "v0.3.9",
    )
    go_repository(
        name = "com_github_tklauser_numcpus",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tklauser/numcpus",
        sum = "h1:ILuRUQBtssgnxw0XXIjKUC56fgnOrFoQQ/4+DeU2biQ=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_github_tmc_grpc_websocket_proxy",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tmc/grpc-websocket-proxy",
        sum = "h1:uruHq4dN7GR16kFc5fp3d1RIYzJW5onx8Ybykw2YQFA=",
        version = "v0.0.0-20201229170055-e5319fda7802",
    )
    go_repository(
        name = "com_github_tomarrell_wrapcheck_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tomarrell/wrapcheck/v2",
        sum = "h1:mU4H9KsqqPZUALOUbVOpjy8qNQbWLoLI9fV68/1tq30=",
        version = "v2.4.0",
    )
    go_repository(
        name = "com_github_tomasen_realip",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tomasen/realip",
        sum = "h1:fb190+cK2Xz/dvi9Hv8eCYJYvIGUTN2/KLq1pT6CjEc=",
        version = "v0.0.0-20180522021738-f0c99a92ddce",
    )
    go_repository(
        name = "com_github_tommy_muehle_go_mnd_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/tommy-muehle/go-mnd/v2",
        sum = "h1:1t0f8Uiaq+fqKteUR4N9Umr6E99R+lDnLnq7PwX2PPE=",
        version = "v2.4.0",
    )

    go_repository(
        name = "com_github_uber_jaeger_client_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/uber/jaeger-client-go",
        sum = "h1:D6wyKGCecFaSRUpo8lCVbaOOb6ThwMmTEbhRwtKR97o=",
        version = "v2.30.0+incompatible",
    )
    go_repository(
        name = "com_github_uber_jaeger_lib",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/uber/jaeger-lib",
        sum = "h1:td4jdvLcExb4cBISKIpHuGoVXh+dVKhn2Um6rjCsSsg=",
        version = "v2.4.1+incompatible",
    )
    go_repository(
        name = "com_github_ugorji_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ugorji/go",
        sum = "h1:/68gy2h+1mWMrwZFeD1kQialdSzAb432dtpeJ42ovdo=",
        version = "v1.1.7",
    )
    go_repository(
        name = "com_github_ugorji_go_codec",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ugorji/go/codec",
        sum = "h1:YPXUKf7fYbp/y8xloBqZOw2qaVggbfwMlI8WM3wZUJ0=",
        version = "v1.2.7",
    )
    go_repository(
        name = "com_github_ultraware_funlen",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ultraware/funlen",
        sum = "h1:5ylVWm8wsNwH5aWo9438pwvsK0QiqVuUrt9bn7S/iLA=",
        version = "v0.0.3",
    )
    go_repository(
        name = "com_github_ultraware_whitespace",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ultraware/whitespace",
        sum = "h1:If7Va4cM03mpgrNH9k49/VOicWpGoG70XPBFFODYDsg=",
        version = "v0.0.4",
    )

    go_repository(
        name = "com_github_urfave_cli",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/urfave/cli",
        sum = "h1:gsqYFH8bb9ekPA12kRo0hfjngWQjkJPlN9R0N78BoUo=",
        version = "v1.22.2",
    )
    go_repository(
        name = "com_github_urfave_cli_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/urfave/cli/v2",
        sum = "h1:qph92Y649prgesehzOrQjdWyxFOp/QVM+6imKHad91M=",
        version = "v2.3.0",
    )
    go_repository(
        name = "com_github_urfave_negroni",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/urfave/negroni",
        sum = "h1:kIimOitoypq34K7TG7DUaJ9kq/N4Ofuwi1sjz0KipXc=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_uudashr_gocognit",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/uudashr/gocognit",
        sum = "h1:rrSex7oHr3/pPLQ0xoWq108XMU8s678FJcQ+aSfOHa4=",
        version = "v1.0.5",
    )

    go_repository(
        name = "com_github_valyala_bytebufferpool",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/valyala/bytebufferpool",
        sum = "h1:GqA5TC/0021Y/b9FG4Oi9Mr3q7XYx6KllzawFIhcdPw=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_valyala_fasthttp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/valyala/fasthttp",
        sum = "h1:CRq/00MfruPGFLTQKY8b+8SfdK60TxNztjRMnH0t1Yc=",
        version = "v1.40.0",
    )
    go_repository(
        name = "com_github_valyala_fasttemplate",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/valyala/fasttemplate",
        sum = "h1:TVEnxayobAdVkhQfrfes2IzOB6o+z4roRkPF52WA1u4=",
        version = "v1.2.1",
    )
    go_repository(
        name = "com_github_valyala_quicktemplate",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/valyala/quicktemplate",
        sum = "h1:LUPTJmlVcb46OOUY3IeD9DojFpAVbsG+5WFTcjMJzCM=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_github_valyala_tcplisten",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/valyala/tcplisten",
        sum = "h1:rBHj/Xf+E1tRGZyWIWwJDiRY0zc1Js+CV5DqwacVSA8=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_github_vbatts_tar_split",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/vbatts/tar-split",
        sum = "h1:Via6XqJr0hceW4wff3QRzD5gAk/tatMw/4ZA7cTlIME=",
        version = "v0.11.2",
    )
    go_repository(
        name = "com_github_vektah_gqlparser_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/vektah/gqlparser/v2",
        sum = "h1:Yjzp66g6oVq93Jihbi0qhGnf/6zIWjcm8H6gA27zstE=",
        version = "v2.4.6",
    )

    go_repository(
        name = "com_github_viki_org_dnscache",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/viki-org/dnscache",
        sum = "h1:EVObHAr8DqpoJCVv6KYTle8FEImKhtkfcZetNqxDoJQ=",
        version = "v0.0.0-20130720023526-c70c1f23c5d8",
    )

    go_repository(
        name = "com_github_vishvananda_netlink",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/vishvananda/netlink",
        sum = "h1:+UB2BJA852UkGH42H+Oee69djmxS3ANzl2b/JtT1YiA=",
        version = "v1.1.1-0.20210330154013-f5de75959ad5",
    )
    go_repository(
        name = "com_github_vishvananda_netns",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/vishvananda/netns",
        sum = "h1:p4VB7kIXpOQvVn1ZaTIVp+3vuYAXFe3OJEvjbUYJLaA=",
        version = "v0.0.0-20210104183010-2eb08e3e575f",
    )
    go_repository(
        name = "com_github_vmihailenco_msgpack_v5",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/vmihailenco/msgpack/v5",
        sum = "h1:5gO0H1iULLWGhs2H5tbAHIZTV8/cYafcFOr9znI5mJU=",
        version = "v5.3.5",
    )
    go_repository(
        name = "com_github_vmihailenco_tagparser_v2",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/vmihailenco/tagparser/v2",
        sum = "h1:y09buUbR+b5aycVFQs/g70pqKVZNBmxwAhO7/IwNM9g=",
        version = "v2.0.0",
    )

    go_repository(
        name = "com_github_xanzy_ssh_agent",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/xanzy/ssh-agent",
        sum = "h1:wUMzuKtKilRgBAD1sUb8gOwwRr2FGoBVumcjoOACClI=",
        version = "v0.3.0",
    )

    go_repository(
        name = "com_github_xeipuuv_gojsonpointer",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/xeipuuv/gojsonpointer",
        sum = "h1:zGWFAtiMcyryUHoUjUJX0/lt1H2+i2Ka2n+D3DImSNo=",
        version = "v0.0.0-20190905194746-02993c407bfb",
    )
    go_repository(
        name = "com_github_xeipuuv_gojsonreference",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/xeipuuv/gojsonreference",
        sum = "h1:EzJWgHovont7NscjpAxXsDA8S8BMYve8Y5+7cuRE7R0=",
        version = "v0.0.0-20180127040603-bd5ef7bd5415",
    )
    go_repository(
        name = "com_github_xeipuuv_gojsonschema",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/xeipuuv/gojsonschema",
        sum = "h1:LhYJRs+L4fBtjZUfuSZIKGeVu0QRy8e5Xi7D17UxZ74=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_github_xiang90_probing",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/xiang90/probing",
        sum = "h1:eY9dn8+vbi4tKz5Qo6v2eYzo7kUS51QINcR5jNpbZS8=",
        version = "v0.0.0-20190116061207-43a291ad63a2",
    )
    go_repository(
        name = "com_github_xlab_treeprint",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/xlab/treeprint",
        sum = "h1:G/1DjNkPpfZCFt9CSh6b5/nY4VimlbHF3Rh4obvtzDk=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_github_xo_terminfo",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/xo/terminfo",
        sum = "h1:QldyIu/L63oPpyvQmHgvgickp1Yw510KJOqX7H24mg8=",
        version = "v0.0.0-20210125001918-ca9a967f8778",
    )

    go_repository(
        name = "com_github_xordataexchange_crypt",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/xordataexchange/crypt",
        sum = "h1:ESFSdwYZvkeru3RtdrYueztKhOBCSAAzS4Gf+k0tEow=",
        version = "v0.0.3-0.20170626215501-b2862e3d0a77",
    )

    go_repository(
        name = "com_github_yashtewari_glob_intersection",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yashtewari/glob-intersection",
        sum = "h1:6gJvMYQlTDOL3dMsPF6J0+26vwX9MB8/1q3uAdhmTrg=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_yeya24_promlinter",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yeya24/promlinter",
        sum = "h1:goWULN0jH5Yajmu/K+v1xCqIREeB+48OiJ2uu2ssc7U=",
        version = "v0.1.0",
    )
    go_repository(
        name = "com_github_yosssi_ace",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yosssi/ace",
        sum = "h1:tUkIP/BLdKqrlrPwcmH0shwEEhTRHoGnc1wFIWmaBUA=",
        version = "v0.0.5",
    )

    go_repository(
        name = "com_github_yudai_gojsondiff",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yudai/gojsondiff",
        sum = "h1:27cbfqXLVEJ1o8I6v3y9lg8Ydm53EKqHXAOMxEGlCOA=",
        version = "v1.0.0",
    )
    go_repository(
        name = "com_github_yudai_golcs",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yudai/golcs",
        sum = "h1:BHyfKlQyqbsFN5p3IfnEUduWvb9is428/nNb5L3U01M=",
        version = "v0.0.0-20170316035057-ecda9a501e82",
    )
    go_repository(
        name = "com_github_yudai_pp",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yudai/pp",
        sum = "h1:Q4//iY4pNF6yPLZIigmvcl7k/bPgrcTPIFIcmawg5bI=",
        version = "v2.0.1+incompatible",
    )

    go_repository(
        name = "com_github_yuin_goldmark",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yuin/goldmark",
        sum = "h1:fVcFKWvrslecOb/tg+Cc05dkeYx540o0FuFt3nUVDoE=",
        version = "v1.4.13",
    )
    go_repository(
        name = "com_github_yvasiyarov_go_metrics",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yvasiyarov/go-metrics",
        sum = "h1:+lm10QQTNSBd8DVTNGHx7o/IKu9HYDvLMffDhbyLccI=",
        version = "v0.0.0-20140926110328-57bccd1ccd43",
    )
    go_repository(
        name = "com_github_yvasiyarov_gorelic",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yvasiyarov/gorelic",
        sum = "h1:hlE8//ciYMztlGpl/VA+Zm1AcTPHYkHJPbHqE6WJUXE=",
        version = "v0.0.0-20141212073537-a9bba5b9ab50",
    )
    go_repository(
        name = "com_github_yvasiyarov_newrelic_platform_go",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/yvasiyarov/newrelic_platform_go",
        sum = "h1:ERexzlUfuTvpE74urLSbIQW0Z/6hF9t8U4NsJLaioAY=",
        version = "v0.0.0-20140908184405-b21fdbd4370f",
    )
    go_repository(
        name = "com_github_ziutek_mymysql",
        build_file_proto_mode = "disable_global",
        importpath = "github.com/ziutek/mymysql",
        sum = "h1:GB0qdRGsTwQSBVYuVShFBKaXSnSnYYC2d9knnE1LHFs=",
        version = "v1.5.4",
    )

    go_repository(
        name = "com_gitlab_gitlab_org_gitaly_v15",
        build_file_proto_mode = "disable_global",
        importpath = "gitlab.com/gitlab-org/gitaly/v15",
        sum = "h1:acswTBIZg8SjGWcHOycBB7e44Ow/DxjC62+lHcA19s8=",
        version = "v15.5.2",
    )

    go_repository(
        name = "com_gitlab_gitlab_org_labkit",
        build_file_proto_mode = "disable_global",
        importpath = "gitlab.com/gitlab-org/labkit",
        sum = "h1:Vm3NAMZ8RqAunXlvPWby3GJ2R35vsYGP6Uu0YjyMIlY=",
        version = "v1.16.0",
    )

    go_repository(
        name = "com_gitlab_gitlab_org_security_products_analyzers_common_v2",
        build_file_proto_mode = "disable_global",
        importpath = "gitlab.com/gitlab-org/security-products/analyzers/common/v2",
        sum = "h1:e4drZgHymJ2uFu5hgEb42QKkOM2X+Q4AL0yn0FyzSv4=",
        version = "v2.24.0",
    )
    go_repository(
        name = "com_gitlab_gitlab_org_security_products_analyzers_report_v3",
        build_file_proto_mode = "disable_global",
        importpath = "gitlab.com/gitlab-org/security-products/analyzers/report/v3",
        sum = "h1:TxuJhvlHxdyEtO9u3esOCkPp4GBjU6Bu8mWzwx3muSM=",
        version = "v3.7.1",
    )
    go_repository(
        name = "com_gitlab_gitlab_org_security_products_analyzers_ruleset",
        build_file_proto_mode = "disable_global",
        importpath = "gitlab.com/gitlab-org/security-products/analyzers/ruleset",
        sum = "h1:IMkx99kHby0P4ALeVOgCaWEccx/wF/CYvqOAmhG6dJU=",
        version = "v1.0.0",
    )

    go_repository(
        name = "com_google_cloud_go",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go",
        sum = "h1:DNtEKRBAAzeS4KyIory52wWHuClNaXJ5x1F7xa4q+5Y=",
        version = "v0.105.0",
    )
    go_repository(
        name = "com_google_cloud_go_accessapproval",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/accessapproval",
        sum = "h1:a3HOTMJn4UJkd9tP5RsbLulZuTo/0BQNBiCyF+M843I=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_accesscontextmanager",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/accesscontextmanager",
        sum = "h1:k571J6qjHx6aX3Xc2TTuWjnm445gvwqLM9MnU5vl5sM=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_aiplatform",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/aiplatform",
        sum = "h1:QqHZT1IMldf/daXoSnkJWBIqGBsw50X+xP6HSVzLRPo=",
        version = "v1.24.0",
    )
    go_repository(
        name = "com_google_cloud_go_analytics",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/analytics",
        sum = "h1:NKw6PpQi6V1O+KsjuTd+bhip9d0REYu4NevC45vtGp8=",
        version = "v0.12.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigateway",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/apigateway",
        sum = "h1:6ttKca2YZfWXvrQpqtV30sgoweTFiNi4KUJq8NQX9zA=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_apigeeconnect",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/apigeeconnect",
        sum = "h1:MoA5qGxCEdoTuVptVKTfrVPis5cOtuGiA3I+fQ7Surg=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_appengine",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/appengine",
        sum = "h1:QhFG5LFVrX37NLkfUpQ3dAR9Gznup0f6PcdPK7UuNgo=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_google_cloud_go_area120",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/area120",
        sum = "h1:TCMhwWEWhCn8d44/Zs7UCICTWje9j3HuV6nVGMjdpYw=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_artifactregistry",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/artifactregistry",
        sum = "h1:AwrnS6BpiBZ8qvxEq5IQi8r3+qgI4cMoaFbzxCdrCy4=",
        version = "v1.8.0",
    )

    go_repository(
        name = "com_google_cloud_go_asset",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/asset",
        sum = "h1:sMxKnJq59UyrQuA7BZxw5wvMUng/2eLdge5a1OW7p3E=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_assuredworkloads",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/assuredworkloads",
        sum = "h1:v61vhw4RR2n7jd3CcSNJPEK/OJwXemwD58QHO3nH5g4=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_automl",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/automl",
        sum = "h1:mxglGBZ1vDAwoKvNFWnzQElAhpBimmA+MAwiblV0J/U=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_baremetalsolution",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/baremetalsolution",
        sum = "h1:nOCMDfI/vxCyKnVvA4oo172lMBHnTNhyhfBL3/HXb/s=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_batch",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/batch",
        sum = "h1:wvVjJvoIBCYPEz5vMRkVMi79fEeHQpDiCvVkdRsx/D8=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_beyondcorp",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/beyondcorp",
        sum = "h1:dPVGvvuSwN7RBYJ/su6BCc7Ymhd6dmL2zuuVF9v2teE=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_google_cloud_go_bigquery",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/bigquery",
        sum = "h1:JuTk8po4bCKRwObdT0zLb1K0BGkGHJdtgs2GK3j2Gws=",
        version = "v1.42.0",
    )
    go_repository(
        name = "com_google_cloud_go_billing",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/billing",
        sum = "h1:3ANOolJbkUAIo4yDRmKNXpUg9LCJjj1gZWaDe8fbLtQ=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_binaryauthorization",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/binaryauthorization",
        sum = "h1:Dtsg+OM81q68H8I6wuhxpECaaCL+EQwRzh1z3Qk/BPU=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_certificatemanager",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/certificatemanager",
        sum = "h1:LqkQLIx//smOvPcRlc7t06iIMMBR5XKa9hDmxYl/+70=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_channel",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/channel",
        sum = "h1:JXMJPd1FsOmnpAeBp+J/zPPLungzqpsNxtgcfdeVJmg=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_cloudbuild",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/cloudbuild",
        sum = "h1:nRHAdoqDcgxS8f5/J24+Bom4KPEhlFUB47MF67Qi42Y=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_clouddms",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/clouddms",
        sum = "h1:GY1TOZU4fXUfYwm9eUQo3Uv1ocQmDMqLthGeCMMRVyM=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_cloudtasks",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/cloudtasks",
        sum = "h1:Rk6r+yV11IinyZDG5w42VPX8HwBUWKdgiY2kflJ1mzY=",
        version = "v1.7.0",
    )

    go_repository(
        name = "com_google_cloud_go_compute",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/compute",
        sum = "h1:gKVJMEyqV5c/UnpzjjQbo3Rjvvqpr9B1DFSbJC4OXr0=",
        version = "v1.12.1",
    )
    go_repository(
        name = "com_google_cloud_go_compute_metadata",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/compute/metadata",
        sum = "h1:efOwf5ymceDhK6PKMnnrTHP4pppY5L22mle96M1yP48=",
        version = "v0.2.1",
    )
    go_repository(
        name = "com_google_cloud_go_contactcenterinsights",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/contactcenterinsights",
        sum = "h1:hlcZ1oUQ+uD0koFgPqx+MPjM2022heTA3e5OhyGQYfM=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_container",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/container",
        sum = "h1:sju38eJg66Q1SYfjSOb/b0a5GvrCs1cnN3lCsLG4M+A=",
        version = "v1.6.0",
    )

    go_repository(
        name = "com_google_cloud_go_containeranalysis",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/containeranalysis",
        sum = "h1:2824iym832ljKdVpCBnpqm5K94YT/uHTVhNF+dRTXPI=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_datacatalog",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/datacatalog",
        sum = "h1:vYBwR8Sy0jVv6AIWCz37ylpDU7IQm2KgexqzOZePIEc=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataflow",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/dataflow",
        sum = "h1:CW3541Fm7KPTyZjJdnX6NtaGXYFn5XbFC5UcjgALKvU=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataform",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/dataform",
        sum = "h1:fnwkyzCVcPI/TmBheGgpmK2h+hWUIDHcZBincHRhrQ0=",
        version = "v0.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_datafusion",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/datafusion",
        sum = "h1:KSWJ5rZNzYodtUpPIG8fHTMiFealGSMnEH8buPYRxRU=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_google_cloud_go_datalabeling",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/datalabeling",
        sum = "h1:dp8jOF21n/7jwgo/uuA0RN8hvLcKO4q6s/yvwevs2ZM=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataplex",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/dataplex",
        sum = "h1:70G+NfNVqLZ1Xur5HUSo+mhsOKHrEaLhDoAVta0zC6Q=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_dataproc",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/dataproc",
        sum = "h1:6l5Np4HcJL56xytV2WrLkGcX4dD9fWupNLorwwFBVZY=",
        version = "v1.7.0",
    )

    go_repository(
        name = "com_google_cloud_go_dataqna",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/dataqna",
        sum = "h1:gx9jr41ytcA3dXkbbd409euEaWtofCVXYBvJz3iYm18=",
        version = "v0.6.0",
    )

    go_repository(
        name = "com_google_cloud_go_datastore",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/datastore",
        sum = "h1:/May9ojXjRkPBNVrq+oWLqmWCkr4OU5uRY29bu0mRyQ=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_google_cloud_go_datastream",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/datastream",
        sum = "h1:jXIZI27wPR+dZ+V0c2eLg3dtIS9ZIVDRxTAYX+w6Sss=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_deploy",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/deploy",
        sum = "h1:xJVwJtJcUPW0opljx4hmox2zPhwVT6uSJNf/2Uw2YmQ=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_google_cloud_go_dialogflow",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/dialogflow",
        sum = "h1:hz1Sz2AXYnKsKcfg4qSbaNQT83SocTUFpHa+FiQtXqs=",
        version = "v1.18.0",
    )
    go_repository(
        name = "com_google_cloud_go_dlp",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/dlp",
        sum = "h1:XMxKeYF9PluIiGFOpWIy9Vbm0p7wZiz+PmSQENCh+yw=",
        version = "v1.6.0",
    )

    go_repository(
        name = "com_google_cloud_go_documentai",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/documentai",
        sum = "h1:z6qGScSr29w19ZdMkoWnE94BwkRqEYS1zo1XnarlLoY=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_domains",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/domains",
        sum = "h1:pu3JIgC1rswIqi5romW0JgNO6CTUydLYX8zyjiAvO1c=",
        version = "v0.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_edgecontainer",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/edgecontainer",
        sum = "h1:hd6J2n5dBBRuAqnNUEsKWrp6XNPKsaxwwIyzOPZTokk=",
        version = "v0.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_essentialcontacts",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/essentialcontacts",
        sum = "h1:AOuaOcGaL1Z1Hy5L4WjZHGYonI3jyNmjVyxB3UsbfX8=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_eventarc",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/eventarc",
        sum = "h1:yzOkQ4EakJ5ZHjqPVMevvQx2T9ZDX9PAB1kY1GpLivU=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_filestore",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/filestore",
        sum = "h1:EOLyFUmQTPu1BFaPNDp33D4aVCQuM66agwGOIoMe+uw=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_firestore",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/firestore",
        sum = "h1:dMIWvm+3O0E3DM7kcZPH0FBQ94Xg/OMkdTNDaY9itbI=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_functions",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/functions",
        sum = "h1:aMns6pr/uwRSCUO4C4aZRwl3Ii2EVytOz1qb3/+Cl6c=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_gaming",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/gaming",
        sum = "h1:NNdzlx3QIQJ3n07YA7hSM0mZfvjLk37IPsQFFPdJGzY=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkebackup",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/gkebackup",
        sum = "h1:i87EQ2xU6wGLmGtx61eZPgmG34UkFj/4sQdZJVm77cI=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_google_cloud_go_gkeconnect",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/gkeconnect",
        sum = "h1:zAcvDa04tTnGdu6TEZewaLN2tdMtUOJJ7fEceULjguA=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkehub",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/gkehub",
        sum = "h1:JTcTaYQRGsVm+qkah7WzHb6e9sf1C0laYdRPn9aN+vg=",
        version = "v0.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_gkemulticloud",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/gkemulticloud",
        sum = "h1:YmwWEfTvLDlzloheLf0khs0Zm0n5f1UkmyWg/9TkyBM=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_gsuiteaddons",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/gsuiteaddons",
        sum = "h1:ARyw0SXZ8pL32f73qq+rvfBRUJhjaeaPVKU5ezCWXw0=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_iam",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/iam",
        sum = "h1:nsqQC88kT5Iwlm4MeNGTpfMWddp6NB/UOLFTH6m1QfQ=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_iap",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/iap",
        sum = "h1:+qMHlGyocHc3WFAK5U+COzR8/yNocwRbtfpj1/ROyks=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_ids",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/ids",
        sum = "h1:QG0XBQPjp+EAhTMyIbzy3BdE6z50jUPjjb8t4pnFU+M=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_google_cloud_go_iot",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/iot",
        sum = "h1:thTIvR3CovF4c2mY8H9A/vantRJADy2ih5/S1FIZzcE=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_kms",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/kms",
        sum = "h1:uc58n3b/n/F2yDMJzHMbXORkJSh3fzO4/+jju6eR7Zg=",
        version = "v1.5.0",
    )

    go_repository(
        name = "com_google_cloud_go_language",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/language",
        sum = "h1:332UUOZdk3m496ZBCoR/IVcdpKB7HHXHJ/MRmV5oYog=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_lifesciences",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/lifesciences",
        sum = "h1:tIqhivE2LMVYkX0BLgG7xL64oNpDaFFI7teunglt1tI=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_longrunning",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/longrunning",
        sum = "h1:y50CXG4j0+qvEukslYFBCrzaXX0qpFbBzc3PchSu/LE=",
        version = "v0.1.1",
    )
    go_repository(
        name = "com_google_cloud_go_managedidentities",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/managedidentities",
        sum = "h1:DjXiirnpfYFZVznSjZp56nRVtEMkRqOkPSzv3Kdog+8=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_mediatranslation",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/mediatranslation",
        sum = "h1:qAJzpxmEX+SeND10Y/4868L5wfZpo4Y3BIEnIieP4dk=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_memcache",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/memcache",
        sum = "h1:l0/5n8tYsccs7WKs6Io4i4TpEbYfjjNbZFAiUkFiNF8=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_metastore",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/metastore",
        sum = "h1:tWzZ+6StVuoMhUldBB6nxCnxJZvt6Lvcl2KkH2NDzuE=",
        version = "v1.7.0",
    )

    go_repository(
        name = "com_google_cloud_go_monitoring",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/monitoring",
        sum = "h1:zK12aN/jzLRzrRXopEOUaG6kvuF6WJmlv1mXn1L7a6E=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkconnectivity",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/networkconnectivity",
        sum = "h1:dOuHLqgOBPfeOg3pxFTFh26r4U7STEt2cpwrY4tXOHs=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_networkmanagement",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/networkmanagement",
        sum = "h1:BVv0NDQYy+dwOgd/zeNns/WQplmpsLStWoQqkI72Bso=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_google_cloud_go_networksecurity",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/networksecurity",
        sum = "h1:qDEX/3sipg9dS5JYsAY+YvgTjPR63cozzAWop8oZS94=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_notebooks",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/notebooks",
        sum = "h1:zGHnJXlskHACBeiXXIe2gC5nMPos38ofT8UeLt3GNfQ=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_optimization",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/optimization",
        sum = "h1:rdvdT5Ha+ji5U2TMTCg+xpBSBE6Hc5f5vRejf1sO88o=",
        version = "v1.1.0",
    )
    go_repository(
        name = "com_google_cloud_go_orchestration",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/orchestration",
        sum = "h1:VPpBBR/MdrFiMXNoGG+IliEzSmF+sgTRmlPb3SA4+x8=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_orgpolicy",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/orgpolicy",
        sum = "h1:RZNIHYRfKVhK/ZHmpVAjpltK7QovtAOwTSntf9pl1xE=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_google_cloud_go_osconfig",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/osconfig",
        sum = "h1:Bu+46U3mzvQ1X+fAX+f2eZaiOZw3LekD9sze/i4f7+w=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_oslogin",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/oslogin",
        sum = "h1:hTBDUARxmYv/qcVPI3DikVNuhCrleuJsLUjlrETKhVo=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_phishingprotection",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/phishingprotection",
        sum = "h1:OrwHLSRSZyaiOt3tnY33dsKSedxbMzsXvqB21okItNQ=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_policytroubleshooter",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/policytroubleshooter",
        sum = "h1:sBeOCvjJFNT/YRtXpWAX08ww70fcRHKqjVGULZ4Vr0E=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_privatecatalog",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/privatecatalog",
        sum = "h1:Vz86uiHCtNGm1DeC32HeG2VXmOq5JRYA3VRPf8ZEcSg=",
        version = "v0.6.0",
    )

    go_repository(
        name = "com_google_cloud_go_profiler",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/profiler",
        sum = "h1:R6y/xAeifaUXxd2x6w+jIwKxoKl8Cv5HJvcvASTPWJo=",
        version = "v0.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_pubsub",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/pubsub",
        sum = "h1:9cH52jizPUVSSrSe+J16RC9wB0QI7i/cfuCm5UUCcIk=",
        version = "v1.5.0",
    )
    go_repository(
        name = "com_google_cloud_go_recaptchaenterprise_v2",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/recaptchaenterprise/v2",
        sum = "h1:JSN8PT1XCAv2Gx01QkVC0xIjnabwawVPCp8l7VPBU6Q=",
        version = "v2.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommendationengine",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/recommendationengine",
        sum = "h1:6w+WxPf2LmUEqX0YyvfCoYb8aBYOcbIV25Vg6R0FLGw=",
        version = "v0.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_recommender",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/recommender",
        sum = "h1:A50LDNE2+IJuvTCLHOSZM3/J5vZsmd4RcUm2Blh04m0=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_redis",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/redis",
        sum = "h1:oMeAQhZG9FJbsnUQGCZYDapWNPx2flDrYxxrAi3WdWU=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcemanager",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/resourcemanager",
        sum = "h1:lJLboEWuUMexWSpW48DRMnDM6LYbtu2pN5q9YBV3fM8=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_resourcesettings",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/resourcesettings",
        sum = "h1:Tr87MO4TbLCKkh/+8Vrx+NL0bxWiQNLOJL8dan11L7k=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_retail",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/retail",
        sum = "h1:wIvFcq5dgo8mnrmEFVf//yUiLzqUHGF1fyqIwiN+R6Q=",
        version = "v1.10.0",
    )
    go_repository(
        name = "com_google_cloud_go_run",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/run",
        sum = "h1:ojwQUZe8P1xAcI8H91HJC/NkxNcKmb9/a9CZpkkZJzE=",
        version = "v0.2.0",
    )

    go_repository(
        name = "com_google_cloud_go_scheduler",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/scheduler",
        sum = "h1:iwMNmT8b2LhkcsHYYDCIxrBAVCklbbbbXvbfVLMLvX4=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_secretmanager",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/secretmanager",
        sum = "h1:4wYWL2t10q+xUtFFS0QuWlqwQguMrwC6FDpjtMM6cUI=",
        version = "v1.8.0",
    )

    go_repository(
        name = "com_google_cloud_go_security",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/security",
        sum = "h1:o9frPOtXW2f4zMlw4SYPE42LRz/hhrYVWtAEUkPvyA4=",
        version = "v1.9.0",
    )
    go_repository(
        name = "com_google_cloud_go_securitycenter",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/securitycenter",
        sum = "h1:Fg6hwy3Of6Q9IvBC/P0pppkUvKHmkRNy++K4/g+PZAw=",
        version = "v1.15.0",
    )
    go_repository(
        name = "com_google_cloud_go_servicecontrol",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/servicecontrol",
        sum = "h1:jyC94OAVsoHBDfpInk66a5ZQJZPgRzSWHQX3OFfnNqE=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_google_cloud_go_servicedirectory",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/servicedirectory",
        sum = "h1:6lhsfUhTmHsiCEdmJd8mBvOO9FdDeZSzUOlGl14x/bw=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_servicemanagement",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/servicemanagement",
        sum = "h1:uqHy7K4df2Z+WwJUhaGMXOdwhZu03fCHbD7xRP3+B8M=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_serviceusage",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/serviceusage",
        sum = "h1:meJ9kF0A9F9f7ztBuyI7tMPThCmulzCI6tAiYtLmgZU=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_shell",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/shell",
        sum = "h1:5kEngDPazPYCyXJyDSgxu8cm4rhC8zgw90IxebjAFaM=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_spanner",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/spanner",
        sum = "h1:mvgDB+f4CfKXebTIi36aqY02eVG4nU/Z3KV6stQaYCc=",
        version = "v1.7.0",
    )
    go_repository(
        name = "com_google_cloud_go_speech",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/speech",
        sum = "h1:j0SMkg35JKbVq185AR4mBR2q0ARFa/CO3g7Mkg35ZY0=",
        version = "v1.8.0",
    )

    go_repository(
        name = "com_google_cloud_go_storage",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/storage",
        sum = "h1:YOO045NZI9RKfCj1c5A/ZtuuENUc8OAW+gHdGnDgyMQ=",
        version = "v1.27.0",
    )
    go_repository(
        name = "com_google_cloud_go_storagetransfer",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/storagetransfer",
        sum = "h1:Xwtox/X2Jl3ariFevzZjzMlsvp3znzvryYZVxMwfZ9g=",
        version = "v1.5.0",
    )

    go_repository(
        name = "com_google_cloud_go_talent",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/talent",
        sum = "h1:+hCSx3Ua8RomkTi8LpZ7dJdSVciWRod7rUlqUDnRVJA=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_texttospeech",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/texttospeech",
        sum = "h1:m0TASV/2VF1lVmXe+VDIur6rbf2gJeAKvnm27lBCg4g=",
        version = "v1.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_tpu",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/tpu",
        sum = "h1:AH8Q8iwRW96Jt7EsuxpxOoVMBjuo+bau9yYxdO6xPY8=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_trace",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/trace",
        sum = "h1:YmCusGVDe0scOkHKi+y/GgBZlwQjHybtQHV12AfJEgo=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_translate",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/translate",
        sum = "h1:uxcD/O79CBRTXN0pOaETVvgyP8tM6biRqG0aVXOWWbs=",
        version = "v1.3.0",
    )
    go_repository(
        name = "com_google_cloud_go_video",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/video",
        sum = "h1:sTtr2Bn0c229zODjrFmRhu5x/FHT921aBBPG51IxExE=",
        version = "v1.8.0",
    )

    go_repository(
        name = "com_google_cloud_go_videointelligence",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/videointelligence",
        sum = "h1:uwarOUvggPELA1Ihen1sDBeyLd2fbIBoiphlMgkFg24=",
        version = "v1.8.0",
    )
    go_repository(
        name = "com_google_cloud_go_vision_v2",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/vision/v2",
        sum = "h1:V7mNTkweq+C+1GquSQiCszoMvurbzoZkqtGcbkzdaO8=",
        version = "v2.4.0",
    )
    go_repository(
        name = "com_google_cloud_go_vmmigration",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/vmmigration",
        sum = "h1:pGVbZnqcHs61mOckeZuatk5aCujGpQlG4EnF2+3jvOE=",
        version = "v1.2.0",
    )
    go_repository(
        name = "com_google_cloud_go_vpcaccess",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/vpcaccess",
        sum = "h1:/PHFPPMOO9M7X00pq0Tn+/ScVDVJ2AATn+78tM7is1M=",
        version = "v1.4.0",
    )

    go_repository(
        name = "com_google_cloud_go_webrisk",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/webrisk",
        sum = "h1:5Y961bETa0Nja4Ej/NdDDz1rhCKf1zwAo22lzJy23Qs=",
        version = "v1.6.0",
    )
    go_repository(
        name = "com_google_cloud_go_websecurityscanner",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/websecurityscanner",
        sum = "h1:rkU8GGJ3GY/FoyhAM/tyZA4eNSoNmcIoY+oQkPGBhSc=",
        version = "v1.3.0",
    )

    go_repository(
        name = "com_google_cloud_go_workflows",
        build_file_proto_mode = "disable_global",
        importpath = "cloud.google.com/go/workflows",
        sum = "h1:las2VwuMdEytyyDhTLEKORtig48dhOViDd1U4RSIsr8=",
        version = "v1.8.0",
    )

    go_repository(
        name = "com_shuralyov_dmitri_gpu_mtl",
        build_file_proto_mode = "disable_global",
        importpath = "dmitri.shuralyov.com/gpu/mtl",
        sum = "h1:VpgP7xuJadIUuKccphEpTJnWhS2jkQyMt6Y7pJCD7fY=",
        version = "v0.0.0-20190408044501-666a987793e9",
    )

    go_repository(
        name = "dev_gocloud",
        build_file_proto_mode = "disable_global",
        importpath = "gocloud.dev",
        sum = "h1:4rM/SVL0lLs+rhC0Gmc+gt/82DBpb7nbpIZKXXnfMXg=",
        version = "v0.26.0",
    )

    go_repository(
        name = "in_gopkg_alecthomas_kingpin_v2",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/alecthomas/kingpin.v2",
        sum = "h1:jMFz6MfLP0/4fUyZle81rXUoxOBFi19VUFKVDOQfozc=",
        version = "v2.2.6",
    )
    go_repository(
        name = "in_gopkg_check_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/check.v1",
        sum = "h1:Hei/4ADfdWqJk1ZMxUNpqntNwaWcugrBjAiHlqqRiVk=",
        version = "v1.0.0-20201130134442-10cb98267c6c",
    )
    go_repository(
        name = "in_gopkg_cheggaaa_pb_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/cheggaaa/pb.v1",
        sum = "h1:n1tBJnnK2r7g9OW2btFH91V92STTUevLXYFb8gy9EMk=",
        version = "v1.0.28",
    )
    go_repository(
        name = "in_gopkg_datadog_dd_trace_go_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/DataDog/dd-trace-go.v1",
        sum = "h1:DkD0plWEVUB8v/Ru6kRBW30Hy/fRNBC8hPdcExuBZMc=",
        version = "v1.32.0",
    )
    go_repository(
        name = "in_gopkg_errgo_v2",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/errgo.v2",
        sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
        version = "v2.1.0",
    )
    go_repository(
        name = "in_gopkg_fsnotify_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/fsnotify.v1",
        sum = "h1:xOHLXZwVvI9hhs+cLKq5+I5onOuwQLhQwiu63xxlHs4=",
        version = "v1.4.7",
    )
    go_repository(
        name = "in_gopkg_gcfg_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/gcfg.v1",
        sum = "h1:m8OOJ4ccYHnx2f4gQwpno8nAX5OGOh7RLaaz0pj3Ogs=",
        version = "v1.2.3",
    )

    go_repository(
        name = "in_gopkg_inf_v0",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/inf.v0",
        sum = "h1:73M5CoZyi3ZLMOyDlQh031Cx6N9NDJ2Vvfl76EDAgDc=",
        version = "v0.9.1",
    )
    go_repository(
        name = "in_gopkg_ini_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/ini.v1",
        sum = "h1:Dgnx+6+nfE+IfzjUEISNeydPJh9AXNNsWbGP9KzCsOA=",
        version = "v1.67.0",
    )

    go_repository(
        name = "in_gopkg_natefinch_lumberjack_v2",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/natefinch/lumberjack.v2",
        sum = "h1:1Lc07Kr7qY4U2YPouBjpCLxpiyxIVoxqXgkXLknAOE8=",
        version = "v2.0.0",
    )
    go_repository(
        name = "in_gopkg_neurosnap_sentences_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/neurosnap/sentences.v1",
        sum = "h1:v7ElyP020iEZQONyLld3fHILHWOPs+ntzuQTNPkul8E=",
        version = "v1.0.6",
    )
    go_repository(
        name = "in_gopkg_resty_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/resty.v1",
        sum = "h1:CuXP0Pjfw9rOuY6EP+UvtNvt5DSqHpIxILZKT/quCZI=",
        version = "v1.12.0",
    )
    go_repository(
        name = "in_gopkg_square_go_jose_v2",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/square/go-jose.v2",
        sum = "h1:7odma5RETjNHWJnR32wx8t+Io4djHE1PqxCFx3iiZ2w=",
        version = "v2.5.1",
    )
    go_repository(
        name = "in_gopkg_tomb_v1",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/tomb.v1",
        sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
        version = "v1.0.0-20141024135613-dd632973f1e7",
    )
    go_repository(
        name = "in_gopkg_warnings_v0",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/warnings.v0",
        sum = "h1:wFXVbFY8DY5/xOe1ECiWdKCzZlxgshcYVNkBHstARME=",
        version = "v0.1.2",
    )
    go_repository(
        name = "in_gopkg_yaml_v2",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/yaml.v2",
        sum = "h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=",
        version = "v2.4.0",
    )
    go_repository(
        name = "in_gopkg_yaml_v3",
        build_file_proto_mode = "disable_global",
        importpath = "gopkg.in/yaml.v3",
        sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
        version = "v3.0.1",
    )
    go_repository(
        name = "io_etcd_go_bbolt",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/bbolt",
        sum = "h1:/ecaJf0sk1l4l6V4awd65v2C3ILy7MSj+s/x1ADCIMU=",
        version = "v1.3.6",
    )
    go_repository(
        name = "io_etcd_go_etcd",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/etcd",
        sum = "h1:/RwRVN9EdXAVtdHxP7Ndn/tfmM9/goiwU0QTnLBgS4w=",
        version = "v0.0.0-20200513171258-e048e166ab9c",
    )
    go_repository(
        name = "io_etcd_go_etcd_api_v3",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/etcd/api/v3",
        sum = "h1:OHVyt3TopwtUQ2GKdd5wu3PmmipR4FTwCqoEjSyRdIc=",
        version = "v3.5.4",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_pkg_v3",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/etcd/client/pkg/v3",
        sum = "h1:lrneYvz923dvC14R54XcA7FXoZ3mlGZAgmwhfm7HqOg=",
        version = "v3.5.4",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v2",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/etcd/client/v2",
        sum = "h1:Dcx3/MYyfKcPNLpR4VVQUP5KgYrBeJtktBwEKkw08Ao=",
        version = "v2.305.4",
    )
    go_repository(
        name = "io_etcd_go_etcd_client_v3",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/etcd/client/v3",
        sum = "h1:p83BUL3tAYS0OT/r0qglgc3M1JjhM0diV8DSWAhVXv4=",
        version = "v3.5.4",
    )
    go_repository(
        name = "io_etcd_go_etcd_pkg_v3",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/etcd/pkg/v3",
        sum = "h1:V5Dvl7S39ZDwjkKqJG2BfXgxZ3QREqqKifWQgIw5IM0=",
        version = "v3.5.4",
    )
    go_repository(
        name = "io_etcd_go_etcd_raft_v3",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/etcd/raft/v3",
        sum = "h1:YGrnAgRfgXloBNuqa+oBI/aRZMcK/1GS6trJePJ/Gqc=",
        version = "v3.5.4",
    )
    go_repository(
        name = "io_etcd_go_etcd_server_v3",
        build_file_proto_mode = "disable_global",
        importpath = "go.etcd.io/etcd/server/v3",
        sum = "h1:CMAZd0g8Bn5NRhynW6pKhc4FRg41/0QYy3d7aNm9874=",
        version = "v3.5.4",
    )
    go_repository(
        name = "io_k8s_api",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/api",
        sum = "h1:Q1v5UFfYe87vi5H7NU0p4RXC26PPMT8KOpr1TLQbCMQ=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_apiextensions_apiserver",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/apiextensions-apiserver",
        sum = "h1:bfI4KS31w2f9WM1KLGwnwuVlW3RSRPuIsfNF/3HzR0k=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_apimachinery",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/apimachinery",
        sum = "h1:7o9ium4uyUOM76t6aunP0nZuex7gDf8VGwkR5RcJnQc=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_apiserver",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/apiserver",
        sum = "h1:m7+xGuG5+KYAnEsqaFtDyWMkmMMEOFYlu+NlWv5qSBI=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_cli_runtime",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/cli-runtime",
        sum = "h1:Zs7P7l7db/5J+KDePOVtDlArAa9pZXaDinGWGZl0aM8=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_client_go",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/client-go",
        sum = "h1:oB4Dyl8d6UbfDHD8Bv8evKylzs3BXzzufLiO27xuPs0=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_code_generator",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/code-generator",
        sum = "h1:BEH+wDi90bGyrYcY4abGtUqaOX7G94RRrEu8l+SvIeo=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_component_base",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/component-base",
        sum = "h1:UrsxciGdrCY03ULT1h/S/gXFCOPnLhUVwSyx+hM/zq4=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_component_helpers",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/component-helpers",
        sum = "h1:Ldwi2U50KohMVDKBcVT3zDpIOKFP4bPEy/4Dj4NsoYU=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_cri_api",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/cri-api",
        sum = "h1:0DHL/hpTf4Fp+QkUXFefWcp1fhjXr9OlNdY9X99c+O8=",
        version = "v0.23.1",
    )
    go_repository(
        name = "io_k8s_gengo",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/gengo",
        sum = "h1:TT1WdmqqXareKxZ/oNXEUSwKlLiHzPMyB0t8BaFeBYI=",
        version = "v0.0.0-20211129171323-c02415ce4185",
    )

    go_repository(
        name = "io_k8s_klog_v2",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/klog/v2",
        sum = "h1:atnLQ121W371wYYFawwYx1aEY2eUfs4l3J72wtgAwV4=",
        version = "v2.80.1",
    )
    go_repository(
        name = "io_k8s_kube_openapi",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/kube-openapi",
        sum = "h1:MQ8BAZPZlWk3S9K4a9NCkIFQtZShWqoha7snGixVgEA=",
        version = "v0.0.0-20220803162953-67bda5d908f1",
    )
    go_repository(
        name = "io_k8s_kubectl",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/kubectl",
        sum = "h1:HnWJziEtmsm4JaJiKT33kG0kadx68MXxUE8UEbXnN4U=",
        version = "v0.25.3",
    )

    go_repository(
        name = "io_k8s_metrics",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/metrics",
        sum = "h1:fp5RuALkbwI3UbKITdNYu6sa3LF4JPANR/ofq3oe+Fg=",
        version = "v0.25.3",
    )
    go_repository(
        name = "io_k8s_sigs_apiserver_network_proxy_konnectivity_client",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/apiserver-network-proxy/konnectivity-client",
        sum = "h1:LYqFq+6Cj2D0gFfrJvL7iElD4ET6ir3VDdhDdTK7rgc=",
        version = "v0.0.33",
    )
    go_repository(
        name = "io_k8s_sigs_cli_utils",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/cli-utils",
        sum = "h1:zCUitt54f0/MYj/ajVFnG6XSXMhpZ72O/3RewIchW8w=",
        version = "v0.34.0",
    )
    go_repository(
        name = "io_k8s_sigs_controller_runtime",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/controller-runtime",
        sum = "h1:tUsRCSJVM1QQOOeViGeX3GMT3dQF1eePPw6sEE3xSlg=",
        version = "v0.13.1",
    )

    go_repository(
        name = "io_k8s_sigs_json",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/json",
        sum = "h1:iXTIw73aPyC+oRdyqqvVJuloN1p0AC/kzH07hu3NE+k=",
        version = "v0.0.0-20220713155537-f223a00ba0e2",
    )
    go_repository(
        name = "io_k8s_sigs_kustomize_api",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/kustomize/api",
        sum = "h1:7YM7gW3kYBwtKvoY216ZzY+8hM+lV53LUayghNRJ0vM=",
        version = "v0.12.1",
    )

    go_repository(
        name = "io_k8s_sigs_kustomize_kustomize_v4",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/kustomize/kustomize/v4",
        sum = "h1:cDW6AVMl6t/SLuQaezMET8hgnadZGIAr8tUrxFVOrpg=",
        version = "v4.5.7",
    )
    go_repository(
        name = "io_k8s_sigs_kustomize_kyaml",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/kustomize/kyaml",
        sum = "h1:Qz53EAaFFANyNgyOEJbT/yoIHygK40/ZcvU3rgry2Tk=",
        version = "v0.13.9",
    )
    go_repository(
        name = "io_k8s_sigs_structured_merge_diff_v4",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/structured-merge-diff/v4",
        sum = "h1:PRbqxJClWWYMNV1dhaG4NsibJbArud9kFxnAMREiWFE=",
        version = "v4.2.3",
    )
    go_repository(
        name = "io_k8s_sigs_yaml",
        build_file_proto_mode = "disable_global",
        importpath = "sigs.k8s.io/yaml",
        sum = "h1:a2VclLzOGrwOHDiV8EfBGhvjHvP46CtW5j6POvhYGGo=",
        version = "v1.3.0",
    )
    go_repository(
        name = "io_k8s_utils",
        build_file_proto_mode = "disable_global",
        importpath = "k8s.io/utils",
        sum = "h1:L14f2LWkOxG2rYsuSA3ltQnnST1vMfek/GUk+VemxD4=",
        version = "v0.0.0-20221101230645-61b03e2f6476",
    )
    go_repository(
        name = "io_nhooyr_websocket",
        build_file_proto_mode = "disable_global",
        importpath = "nhooyr.io/websocket",
        sum = "h1:usjR2uOr/zjjkVMy0lW+PPohFok7PCow5sDjLgX4P4g=",
        version = "v1.8.7",
    )
    go_repository(
        name = "io_opencensus_go",
        build_file_proto_mode = "disable_global",
        importpath = "go.opencensus.io",
        sum = "h1:gqCw0LfLxScz8irSi8exQc7fyQ0fKQU/qnC/X8+V/1M=",
        version = "v0.23.0",
    )

    go_repository(
        name = "io_opencensus_go_contrib_exporter_stackdriver",
        build_file_proto_mode = "disable_global",
        importpath = "contrib.go.opencensus.io/exporter/stackdriver",
        sum = "h1:a9+GZPUe+ONKUwULjlEOucMMG0qfSCCenlji0Nhqbys=",
        version = "v0.13.10",
    )

    go_repository(
        name = "io_opentelemetry_go_contrib",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/contrib",
        sum = "h1:ubFQUn0VCZ0gPwIoJfBJVpeBlyRMxu8Mm/huKWYd9p0=",
        version = "v0.20.0",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_google_golang_org_grpc_otelgrpc",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc",
        sum = "h1:PRXhsszxTt5bbPriTjmaweWUsAnJYeWBhUMLRetUgBU=",
        version = "v0.36.4",
    )
    go_repository(
        name = "io_opentelemetry_go_contrib_instrumentation_net_http_otelhttp",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp",
        sum = "h1:aUEBEdCa6iamGzg6fuYxDA8ThxvOG240mAvWDU+XLio=",
        version = "v0.36.4",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/exporters/otlp",
        sum = "h1:PTNgq9MRmQqqJY0REVbZFvwkYOA85vbdQU/nVfxDyqg=",
        version = "v0.20.0",
    )

    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_internal_retry",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/internal/retry",
        sum = "h1:X2GndnMCsUPh6CiY2a+frAbNsXaPLbB0soHRYhAZ5Ig=",
        version = "v1.11.1",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace",
        sum = "h1:MEQNafcNCB0uQIti/oHgU7CZpUMYQ7qigBwMVKycHvc=",
        version = "v1.11.1",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_exporters_otlp_otlptrace_otlptracegrpc",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc",
        sum = "h1:MFAyzUPrTwLOwCi+cltN0ZVyy4phU41lwH+lyMyQTS4=",
        version = "v1.7.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_metric",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/metric",
        sum = "h1:xQAyl7uGEYvrLAiV/09iTJlp1pZnQ9Wl793qbVvED1E=",
        version = "v0.33.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_oteltest",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/oteltest",
        sum = "h1:YVfA0ByROYqTwOxqHVZYZExzEpfZor+MU1rU+ip2v9Q=",
        version = "v0.19.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/sdk",
        sum = "h1:F7KmQgoHljhUuJyA+9BiU+EkJfyX5nVVF4wyzWZpKxs=",
        version = "v1.11.1",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_export_metric",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/sdk/export/metric",
        sum = "h1:c5VRjxCXdQlx1HjzwGdQHzZaVI82b5EbBgOu2ljD92g=",
        version = "v0.20.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_sdk_metric",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/sdk/metric",
        sum = "h1:7ao1wpzHRVKf0OQ7GIxiQJA6X7DLX9o14gmVon7mMK8=",
        version = "v0.20.0",
    )
    go_repository(
        name = "io_opentelemetry_go_otel_trace",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/otel/trace",
        sum = "h1:ofxdnzsNrGBYXbP7t7zpUK281+go5rF7dvdIZXF8gdQ=",
        version = "v1.11.1",
    )
    go_repository(
        name = "io_opentelemetry_go_proto_otlp",
        build_file_proto_mode = "disable_global",
        importpath = "go.opentelemetry.io/proto/otlp",
        sum = "h1:IVN6GR+mhC4s5yfcTbmzHYODqvWAp3ZedA2SJPI1Nnw=",
        version = "v0.19.0",
    )
    go_repository(
        name = "io_rsc_binaryregexp",
        build_file_proto_mode = "disable_global",
        importpath = "rsc.io/binaryregexp",
        sum = "h1:HfqmD5MEmC0zvwBuF187nq9mdnXjXsSivRiXN7SmRkE=",
        version = "v0.2.0",
    )

    go_repository(
        name = "io_rsc_quote_v3",
        build_file_proto_mode = "disable_global",
        importpath = "rsc.io/quote/v3",
        sum = "h1:9JKUTTIUgS6kzR9mK1YuGKv6Nl+DijDNIc0ghT58FaY=",
        version = "v3.1.0",
    )
    go_repository(
        name = "io_rsc_sampler",
        build_file_proto_mode = "disable_global",
        importpath = "rsc.io/sampler",
        sum = "h1:7uVkIFmeBqHfdjD+gZwtXXI+RODJ2Wc4O7MPEh/QiW4=",
        version = "v1.3.0",
    )

    go_repository(
        name = "land_oras_oras_go",
        build_file_proto_mode = "disable_global",
        importpath = "oras.land/oras-go",
        sum = "h1:yoKosVIbsPoFMqAIFHTnrmOuafHal+J/r+I5bdbVWu4=",
        version = "v1.2.0",
    )

    go_repository(
        name = "net_starlark_go",
        build_file_proto_mode = "disable_global",
        importpath = "go.starlark.net",
        sum = "h1:+FNtrFTmVw0YZGpBGX56XDee331t6JAXeK2bcyhLOOc=",
        version = "v0.0.0-20200306205701-8dd3e2ee1dd5",
    )

    go_repository(
        name = "org_bitbucket_creachadair_shell",
        build_file_proto_mode = "disable_global",
        importpath = "bitbucket.org/creachadair/shell",
        sum = "h1:reJflDbKqnlnqb4Oo2pQ1/BqmY/eCWcNGHrIUO8qIzc=",
        version = "v0.0.6",
    )

    go_repository(
        name = "org_golang_google_api",
        build_file_proto_mode = "disable_global",
        importpath = "google.golang.org/api",
        sum = "h1:JxJl2qQ85fRMPNvlZY/enexbxpCjLwGhZUtgfGeQ51I=",
        version = "v0.102.0",
    )
    go_repository(
        name = "org_golang_google_appengine",
        build_file_proto_mode = "disable_global",
        importpath = "google.golang.org/appengine",
        sum = "h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=",
        version = "v1.6.7",
    )

    go_repository(
        name = "org_golang_google_genproto",
        build_file_proto_mode = "disable_global",
        importpath = "google.golang.org/genproto",
        sum = "h1:QgY/XxIAIeccR+Ca/rDdKubLIU9rcJ3xfy1DC/Wd2Oo=",
        version = "v0.0.0-20221027153422-115e99e71e1c",
    )
    go_repository(
        name = "org_golang_google_grpc",
        build_file_proto_mode = "disable_global",
        importpath = "google.golang.org/grpc",
        sum = "h1:DS/BukOZWp8s6p4Dt/tOaJaTQyPyOoCcrjroHuCeLzY=",
        version = "v1.50.1",
    )
    go_repository(
        name = "org_golang_google_grpc_cmd_protoc_gen_go_grpc",
        build_file_proto_mode = "disable_global",
        importpath = "google.golang.org/grpc/cmd/protoc-gen-go-grpc",
        sum = "h1:M1YKkFIboKNieVO5DLUEVzQfGwJD30Nv2jfUgzb5UcE=",
        version = "v1.1.0",
    )
    go_repository(
        name = "org_golang_google_protobuf",
        build_file_proto_mode = "disable_global",
        importpath = "google.golang.org/protobuf",
        sum = "h1:d0NfwRgPtno5B1Wa6L2DAG+KivqkdutMf1UhdNx175w=",
        version = "v1.28.1",
    )
    go_repository(
        name = "org_golang_x_crypto",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/crypto",
        sum = "h1:fmw3UbQh+nxngCAHrDCCztao/kbYFnWjoqop8dHx05A=",
        version = "v0.0.0-20220926161630-eccd6366d1be",
    )
    go_repository(
        name = "org_golang_x_exp",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/exp",
        sum = "h1:FR+oGxGfbQu1d+jglI3rCkjAjUnhRSZcUxr+DqlDLNo=",
        version = "v0.0.0-20200331195152-e8c3332aa8e5",
    )
    go_repository(
        name = "org_golang_x_image",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/image",
        sum = "h1:+qEpEAPhDZ1o0x3tHzZTQDArnOixOzGD9HUJfcg0mb4=",
        version = "v0.0.0-20190802002840-cff245a6509b",
    )
    go_repository(
        name = "org_golang_x_lint",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/lint",
        sum = "h1:VLliZ0d+/avPrXXH+OakdXhpJuEoBZuwh1m2j7U6Iug=",
        version = "v0.0.0-20210508222113-6edffad5e616",
    )
    go_repository(
        name = "org_golang_x_mobile",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/mobile",
        sum = "h1:4+4C/Iv2U4fMZBiMCc98MG1In4gJY5YRhtpDNeDeHWs=",
        version = "v0.0.0-20190719004257-d2bd2a29d028",
    )
    go_repository(
        name = "org_golang_x_mod",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/mod",
        sum = "h1:6zppjxzCulZykYSLyVDYbneBfbaBIQPYMevg0bEwv2s=",
        version = "v0.6.0-dev.0.20220419223038-86c51ed26bb4",
    )
    go_repository(
        name = "org_golang_x_net",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/net",
        sum = "h1:hZ/3BUoy5aId7sCpA/Tc5lt8DkFgdVS2onTpJsZ/fl0=",
        version = "v0.1.0",
    )
    go_repository(
        name = "org_golang_x_oauth2",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/oauth2",
        sum = "h1:nt+Q6cXKz4MosCSpnbMtqiQ8Oz0pxTef2B4Vca2lvfk=",
        version = "v0.0.0-20221014153046-6fdb5e3db783",
    )
    go_repository(
        name = "org_golang_x_sync",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/sync",
        sum = "h1:wsuoTGHzEhffawBOhz5CYhcrV4IdKZbEyZjBMuTp12o=",
        version = "v0.1.0",
    )
    go_repository(
        name = "org_golang_x_sys",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/sys",
        sum = "h1:kunALQeHf1/185U1i0GOB/fy1IPRDDpuoOOqRReG57U=",
        version = "v0.1.0",
    )
    go_repository(
        name = "org_golang_x_term",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/term",
        sum = "h1:g6Z6vPFA9dYBAF7DWcH6sCcOntplXsDKcliusYijMlw=",
        version = "v0.1.0",
    )
    go_repository(
        name = "org_golang_x_text",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/text",
        sum = "h1:BrVqGRd7+k1DiOgtnFvAkoQEWQvBc25ouMJM6429SFg=",
        version = "v0.4.0",
    )
    go_repository(
        name = "org_golang_x_time",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/time",
        sum = "h1:xYY+Bajn2a7VBmTM5GikTmnK8ZuX8YgnQCqZpbBNtmA=",
        version = "v0.1.0",
    )
    go_repository(
        name = "org_golang_x_xerrors",
        build_file_proto_mode = "disable_global",
        importpath = "golang.org/x/xerrors",
        sum = "h1:H2TDz8ibqkAF6YGhCdN3jS9O0/s90v0rJh3X/OLHEUk=",
        version = "v0.0.0-20220907171357-04be3eba64a2",
    )
    go_repository(
        name = "org_gonum_v1_gonum",
        build_file_proto_mode = "disable_global",
        importpath = "gonum.org/v1/gonum",
        sum = "h1:CCXrcPKiGGotvnN6jfUsKk4rRqm7q09/YbKb5xCEvtM=",
        version = "v0.8.2",
    )
    go_repository(
        name = "org_mozilla_go_mozlog",
        build_file_proto_mode = "disable_global",
        importpath = "go.mozilla.org/mozlog",
        sum = "h1:rKyWXYDfrVOpMFBion4Pmx5sJbQreQNXycHvm4KwJSg=",
        version = "v0.0.0-20170222151521-4bb13139d403",
    )

    go_repository(
        name = "org_mozilla_go_pkcs7",
        build_file_proto_mode = "disable_global",
        importpath = "go.mozilla.org/pkcs7",
        sum = "h1:A/5uWzF44DlIgdm/PQFwfMkW0JX+cIcQi/SwLAmZP5M=",
        version = "v0.0.0-20200128120323-432b2356ecb1",
    )
    go_repository(
        name = "org_uber_go_atomic",
        build_file_proto_mode = "disable_global",
        importpath = "go.uber.org/atomic",
        sum = "h1:ECmE8Bn/WFTYwEW/bpKD3M8VtR/zQVbavAoalC1PYyE=",
        version = "v1.9.0",
    )
    go_repository(
        name = "org_uber_go_automaxprocs",
        build_file_proto_mode = "disable_global",
        importpath = "go.uber.org/automaxprocs",
        sum = "h1:e1YG66Lrk73dn4qhg8WFSvhF0JuFQF0ERIp4rpuV8Qk=",
        version = "v1.5.1",
    )

    go_repository(
        name = "org_uber_go_goleak",
        build_file_proto_mode = "disable_global",
        importpath = "go.uber.org/goleak",
        sum = "h1:xqgm/S+aQvhWFTtR0XK3Jvg7z8kGV8P4X14IzwN3Eqk=",
        version = "v1.2.0",
    )
    go_repository(
        name = "org_uber_go_multierr",
        build_file_proto_mode = "disable_global",
        importpath = "go.uber.org/multierr",
        sum = "h1:dg6GjLku4EH+249NNmoIciG9N/jURbDG+pFlTkhzIC8=",
        version = "v1.8.0",
    )
    go_repository(
        name = "org_uber_go_tools",
        build_file_proto_mode = "disable_global",
        importpath = "go.uber.org/tools",
        sum = "h1:0mgffUl7nfd+FpvXMVz4IDEaUSmT1ysygQC7qYo7sG4=",
        version = "v0.0.0-20190618225709-2cfd321de3ee",
    )

    go_repository(
        name = "org_uber_go_zap",
        build_file_proto_mode = "disable_global",
        importpath = "go.uber.org/zap",
        sum = "h1:OjGQ5KQDEUawVHxNwQgPpiypGHOxo2mNZsOqTak4fFY=",
        version = "v1.23.0",
    )
    go_repository(
        name = "sh_helm_helm_v3",
        build_file_proto_mode = "disable_global",
        importpath = "helm.sh/helm/v3",
        sum = "h1:uTnNlYx8QcTSNA4ZJ50Llwife4CSohUY4ehumyVf2QE=",
        version = "v3.10.1",
    )

    go_repository(
        name = "tools_gotest_v3",
        build_file_proto_mode = "disable_global",
        importpath = "gotest.tools/v3",
        sum = "h1:4AuOwCGf4lLR9u3YOe2awrHygurzhO/HeQ6laiA6Sx0=",
        version = "v3.0.3",
    )
    go_repository(
        name = "xyz_gomodules_jsonpatch_v2",
        build_file_proto_mode = "disable_global",
        importpath = "gomodules.xyz/jsonpatch/v2",
        sum = "h1:4pT439QV83L+G9FkcCriY6EkpcK6r6bK+A5FBUMI7qY=",
        version = "v2.2.0",
    )

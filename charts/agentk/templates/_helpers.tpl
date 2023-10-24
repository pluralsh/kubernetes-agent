{{/*
Expand the name of the chart.
*/}}
{{- define "agentk.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "agentk.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "agentk.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "agentk.labels" -}}
helm.sh/chart: {{ include "agentk.chart" . }}
{{ include "agentk.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- if .Values.additionalLabels }}
{{ toYaml .Values.additionalLabels }}
{{- end -}}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "agentk.selectorLabels" -}}
app.kubernetes.io/name: {{ include "agentk.name" . }}
{{- end }}

{{/*
Secret Name
*/}}
{{- define "agentk.secretName" -}}
{{- if .Values.config.secretName }}
{{- .Values.config.secretName }}
{{- else }}
{{- printf "%s-token" (include "agentk.fullname" .) -}}
{{- end }}
{{- end }}

{{/*
Common annotations
*/}}
{{- define "agentk.annotations" -}}
{{- $ := index . 0 }}
{{- $annotations := index . 1 }}
{{- $observabilityEnabled := ($.Values.config.observability).enabled }}
{{- $token := $.Values.config.token }}
{{- /*
If observability is disabled, always remove the prometheus annotations to avoid Prometheus scraping a closed port
*/ -}}
{{- if (not $observabilityEnabled) }}
{{- $annotations := unset $annotations "prometheus.io/path" }}
{{- $annotations := unset $annotations "prometheus.io/port" }}
{{- $annotations := unset $annotations "prometheus.io/scrape" }}
{{- end }}
{{- with $token }}
{{ printf "checksum/token: %s" (. | sha256sum) }}
{{- end }}
{{- with $annotations }}
{{ toYaml $annotations }}
{{- end }}
{{- if and (not $token) (not $annotations) }}
{{ printf "{}" -}}
{{- end }}
{{- end }}

{{/*
Observability TLS Secret Name
*/}}
{{- define "agentk.observabilitySecretName" -}}
{{- $name := (((.Values.config.observability).tls).secret).name }}
{{- if $name }}
{{-   $name }}
{{- else }}
{{-   printf "%s-observability" (include "agentk.fullname" .) -}}
{{- end }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "agentk.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "agentk.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

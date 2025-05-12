{{/* vim: set filetype=helm: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "microservice-demo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "microservice-demo.fullname" -}}
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
{{- define "microservice-demo.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "microservice-demo.labels" -}}
helm.sh/chart: {{ include "microservice-demo.chart" . }}
{{ include "microservice-demo.selectorLabelsBase" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Base selector labels (used by multiple components if needed)
*/}}
{{- define "microservice-demo.selectorLabelsBase" -}}
app.kubernetes.io/name: {{ include "microservice-demo.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Selector labels for Order Service
*/}}
{{- define "microservice-demo.orderService.selectorLabels" -}}
{{- include "microservice-demo.selectorLabelsBase" . }}
app.kubernetes.io/component: order-service
{{- end }}

{{/*
Selector labels for Inventory Service
*/}}
{{- define "microservice-demo.inventoryService.selectorLabels" -}}
{{- include "microservice-demo.selectorLabelsBase" . }}
app.kubernetes.io/component: inventory-service
{{- end }}


{{/*
Create the name of the service account to use
*/}}
{{- define "microservice-demo.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (printf "%s-%s" .Release.Name .Values.serviceAccount.name) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}


{{/*
Get OTel Collector Endpoint Host (without port).
*/}}
{{- define "microservice-demo.otelCollectorHost" -}}
{{- $defaultHost := "otel-collector.observability.svc.cluster.local" -}}
{{- $host := $defaultHost -}}
{{- if .Values.otel.collector.endpoint -}}
  {{- $endpoint := .Values.otel.collector.endpoint -}}
  {{- $endpointClean := trim (trimPrefix "http://" (trimPrefix "https://" $endpoint)) -}}
  {{- if contains ":" $endpointClean -}}
    {{- $parts := split ":" $endpointClean -}}
    {{- if eq (kindOf $parts) "map" -}}
      {{- /* Access host using map key "_0", provide fallback */}}
      {{- $host = index $parts "_0" | default $endpointClean -}}
    {{- else -}}
      {{- /* Fallback to standard list behavior */}}
      {{- $host = first $parts | default $endpointClean -}}
    {{- end -}}
  {{- else -}}
    {{- $host = $endpointClean -}}
  {{- end -}}
{{- end -}}
{{- print $host -}}
{{- end -}}


{{/*
Get OTel Collector Endpoint Port.
*/}}
{{- define "microservice-demo.otelCollectorPort" -}}
{{- $defaultPort := "4317" -}}
{{- $port := $defaultPort -}}
{{- if .Values.otel.collector.endpoint -}}
  {{- $endpoint := .Values.otel.collector.endpoint -}}
  {{- $endpointClean := trim (trimPrefix "http://" (trimPrefix "https://" $endpoint)) -}}
  {{- if contains ":" $endpointClean -}}
    {{- $parts := split ":" $endpointClean -}}
    {{- if eq (kindOf $parts) "map" -}}
       {{- /* Access port using map key "_1", provide fallback with default */}}
       {{- $port = index $parts "_1" | default $defaultPort -}}
    {{- else -}}
      {{- /* Fallback to standard list behavior */}}
      {{- if gt (len $parts) 1 -}}
        {{- $port = last $parts | default $defaultPort -}}
      {{- else -}}
        {{- $port = $defaultPort -}}
      {{- end -}}
    {{- end -}}
  {{- else -}}
    {{- $port = $defaultPort -}}
  {{- end -}}
{{- end -}}
{{- print $port -}}
{{- end -}}

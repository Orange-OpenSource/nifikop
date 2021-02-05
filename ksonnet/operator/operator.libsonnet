// WARNING: This does NOT deploy the CRDs!

local k = (import 'github.com/grafana/jsonnet-libs/ksonnet-util/kausal.libsonnet');

k {
  // Namespace for the manifests
  nifikop: {
    local container = $.core.v1.container,
    local deployment = $.apps.v1.deployment,
    local policyRule = $.rbac.v1.policyRule,

    // FIXME this is too permissive {{{
    rbac:
      local rbac = if $._config.nifikop.namespace_scoped then $.util.namespacedRBAC else $.util.rbac;
      rbac($._config.nifikop.sa_name, [

        policyRule.withApiGroups([''])
        + policyRule.withResources([
          'pods',
          'services',
          'endpoints',
          'persistentvolumeclaims',
          'events',
          'configmaps',
          'secrets',
        ])
        + policyRule.withVerbs(['create', 'delete', 'get', 'list', 'patch', 'update', 'watch']),

        policyRule.withApiGroups(['policy'])
        + policyRule.withResources(['poddisruptionbudgets'])
        + policyRule.withVerbs(['create', 'delete', 'get', 'list', 'patch', 'update', 'watch']),

        policyRule.withApiGroups(['apps'])
        + policyRule.withResources([
          'deployments',
          'replicasets',
          'statefulsets',
        ])
        + policyRule.withVerbs(['create', 'delete', 'get', 'list', 'patch', 'update', 'watch']),

        policyRule.withApiGroups(['coordination.k8s.io'])
        + policyRule.withResources([
          'leases',
        ])
        + policyRule.withVerbs(['create', 'delete', 'get', 'list', 'patch', 'update', 'watch']),

        policyRule.withApiGroups(['apps'])
        + policyRule.withResources([
          'infikop/finalizers',
        ])
        + policyRule.withVerbs(['update']),

        policyRule.withApiGroups(['nifi.orange.com'])
        + policyRule.withResources([
          'nifiusers',
          'nifiusergroups',
          'nificlusters',
          'nifidataflows',
          'nifiregistryclients',
          'nifiparametercontexts',
        ])
        + policyRule.withVerbs(['create', 'delete', 'get', 'list', 'patch', 'update', 'watch', 'deletecollection']),

        policyRule.withApiGroups(['nifi.orange.com'])
        + policyRule.withResources([
          'nifiusers/status',
          'nifiusergroups/status',
          'nificlusters/status',
          'nifidataflows/status',
          'nifiregistryclients/status',
          'nifiparametercontexts/status',
        ])
        + policyRule.withVerbs(['get', 'update', 'patch']),

        policyRule.withApiGroups(['cert-manager.io'])
        + policyRule.withResources(
          ['issuers', 'certificates']
          + if $._config.nifikop.namespace_scoped then [] else ['clusterissuers']
        )
        + policyRule.withVerbs(['create', 'delete', 'get', 'list', 'patch', 'update', 'watch']),
      ])
      + {
        local mixin = { metadata+: { labels+: $._config.nifikop.labels } },
        role+: mixin,
        cluster_role_binding+: mixin,
        service_account+: mixin,
      },
    // }}}

    container::
      container.new('nifikop-manager', $._images.nifikop)
      + container.withCommand(['/manager'])
      + container.withArgs([
        '--leader-elect',
        '--cert-manager-enabled=%(certmanager_enabled)s' % $._config.nifikop,
      ])
      + container.withPorts([
        $.core.v1.containerPort.new('http-health', 8081),
        $.core.v1.containerPort.new('http-metrics', 9710),
      ])
      + container.withEnvMap({
        WATCH_NAMESPACE: std.join(',', $._config.nifikop.namespaces),
        OPERATOR_NAME: $._config.nifikop.name,
      })
      + container.withEnvMixin([
        k.core.v1.envVar.fromFieldPath('POD_NAME', 'metadata.name'),
      ])
      // Readieness Probe
      + container.readinessProbe.httpGet.withPath('/readyz')
      + container.readinessProbe.httpGet.withPort(8081)
      + container.readinessProbe.withInitialDelaySeconds(15)
      + container.readinessProbe.withPeriodSeconds(20)
      // Liveness Probe
      + container.livenessProbe.httpGet.withPath('/healthz')
      + container.livenessProbe.httpGet.withPort(8081)
      + container.livenessProbe.withInitialDelaySeconds(5)
      + container.livenessProbe.withPeriodSeconds(10)
      // Resource Requests
      + $.util.resourcesRequests('100m', '128Mi')
      + $.util.resourcesLimits('500m', '512Mi'),


    deploy:
      deployment.new(
        name=$._config.nifikop.deploy_name,
        replicas=$._config.nifikop.deploy_replicas,
        containers=[self.container],
        podLabels=$._config.nifikop.labels,
      )
      + deployment.metadata.withLabelsMixin($._config.nifikop.labels)
      + deployment.spec.template.metadata.withAnnotations({ 'sidecar.istio.io/inject': 'false' })
      + deployment.spec.template.spec.securityContext.withRunAsUser(1000)
      + deployment.spec.template.spec.withServiceAccountName($._config.nifikop.sa_name)
      + $.util.antiAffinity,

  },

}

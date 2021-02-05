{
  _config+:: {
    namespace: error 'namespace has to be defined',

    // NifiKop operator deployment
    nifikop+: {
      local nifikop = self,

      name: 'nifikop',

      // define listen namespaces for creating clusters inside here
      namespace_scoped: true,  // helps lower CPU load in largish clusters, reduce blast radius
      namespaces: if self.namespace_scoped then [$._config.namespace] else [],

      // CertManager configuration
      certmanager_enabled: true,
      certmanager_cluster_scoped: true,

      // Resource names
      sa_name: self.name,
      deploy_name: self.name,
      deploy_replicas: 2,  // HA mode, we actually don't need 3 as leader election is handled by kubernetes
      labels: { app: nifikop.name },

    },
  },
}

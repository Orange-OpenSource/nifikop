---
id: 3_node_config
title: Node configuration
sidebar_label: Node configuration
---

NodeConfig defines the node configuration

```yaml
   default_group:
      # provenanceStorage allow to specify the maximum amount of data provenance information to store at a time
      # https://nifi.apache.org/docs/nifi-docs/html/administration-guide.html#write-ahead-provenance-repository-properties
      provenanceStorage: "10 GB"
      #RunAsUser define the id of the user to run in the Nifi image
      # +kubebuilder:validation:Minimum=1
      runAsUser: 1000
      # Set this to true if the instance is a node in a cluster.
      # https://nifi.apache.org/docs/nifi-docs/html/administration-guide.html#basic-cluster-setup
      isNode: true
      # Docker image used by the operator to create the node associated
      # https://hub.docker.com/r/apache/nifi/
#      image: "apache/nifi:1.11.2"
      # nodeAffinity can be specified, operator populates this value if new pvc added later to node
      # https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#node-affinity
#      nodeAffinity:
      # imagePullPolicy define the pull policy for NiFi cluster docker image
      imagePullPolicy: IfNotPresent
      # storageConfigs specifies the node related configs
      storageConfigs:
        # Name of the storage config, used to name PV to reuse into sidecars for example.
        - name: provenance-repository
          # Path where the volume will be mount into the main nifi container inside the pod.
          mountPath: "/opt/nifi/provenance_repository"
          # Kubernetes PVC spec
          # https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/#create-a-persistentvolumeclaim
          pvcSpec:
            accessModes:
              - ReadWriteOnce
            storageClassName: "standard"
            resources:
              requests:
                storage: 10Gi
        - mountPath: "/opt/nifi/nifi-current/logs"
          name: logs
          pvcSpec:
            accessModes:
              - ReadWriteOnce
            storageClassName: "standard"
            resources:
              requests:
                storage: 10Gi
```

## NodeConfig

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|provenanceStorage|string|provenanceStorage allow to specify the maximum amount of data provenance information to store at a time: [write-ahead-provenance-repository-properties](https://nifi.apache.org/docs/nifi-docs/html/administration-guide.html#write-ahead-provenance-repository-properties)|No|"8 GB"|
|runAsUser|int64|define the id of the user to run in the Nifi image|No|1000|
|isNode|boolean|Set this to true if the instance is a node in a cluster: [basic-cluster-setup](https://nifi.apache.org/docs/nifi-docs/html/administration-guide.html#basic-cluster-setup)|No|true|
|image|string| Docker image used by the operator to create the node associated. [Nifi docker registry](https://hub.docker.com/r/apache/nifi/)|No|""|
|imagePullPolicy|[PullPolicy](https://godoc.org/k8s.io/api/core/v1#PullPolicy)| define the pull policy for NiFi cluster docker image.)|No|""|
|nodeAffinity|string| operator populates this value if new pvc added later to node [node-affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#node-affinity)|No|nil|
|storageConfigs|\[  \][StorageConfig](#storageconfig)|specifies the node related configs.|No|nil|
|serviceAccountName|string|specifies the serviceAccount used for this specific node.|No|"default"|
|resourcesRequirements|[ResourceRequirements](https://godoc.org/k8s.io/api/core/v1#ResourceRequirements)| works exactly like Container resources, the user can specify the limit and the requests through this property [manage-compute-resources-container](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/).|No|nil|
|imagePullSecrets|\[  \][LocalObjectReference](https://godoc.org/k8s.io/api/core/v1#TypedLocalObjectReference)|specifies the secret to use when using private registry.|No|nil|
|nodeSelector|map\[string\]string|nodeSelector can be specified, which set the pod to fit on a node [nodeselector](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#nodeselector)|No|nil|
|tolerations|\[  \][Toleration](https://godoc.org/k8s.io/api/core/v1#Toleration)|tolerations can be specified, which set the pod's tolerations [taint-and-toleration](https://kubernetes.io/docs/concepts/configuration/taint-and-toleration/#concepts).|No|nil|
|nodeAnnotations|map\[string\]string|Additionnal annotation to attach to the pod associated [annotations](https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/#syntax-and-character-set).|No|nil|

## StorageConfig

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|name|string|Name of the storage config, used to name PV to reuse into sidecars for example.|Yes| - |
|mountPath|string|Path where the volume will be mount into the main nifi container inside the pod.|Yes| - |
|pvcSpec|[PersistentVolumeClaimSpec](https://godoc.org/k8s.io/api/core/v1#PersistentVolumeClaimSpec)|Kubernetes PVC spec. [create-a-persistentvolumeclaim](https://kubernetes.io/docs/tasks/configure-pod-container/configure-persistent-volume-storage/#create-a-persistentvolumeclaim).|Yes| - |
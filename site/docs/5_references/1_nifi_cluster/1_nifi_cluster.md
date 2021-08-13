---
id: 1_nifi_cluster
title: NiFi cluster
sidebar_label: NiFi cluster
---

`NifiCluster` describes the desired state of the NiFi cluster we want to setup through the operator.

```yaml
apiVersion: nifi.orange.com/v1alpha1
kind: NifiCluster
metadata:
  name: simplenifi
spec:
  service:
    headlessEnabled: true
  zkAddress: "zookeepercluster-client.zookeeper:2181"
  zkPath: "/simplenifi"
  clusterImage: "apache/nifi:1.11.3"
  oneNifiNodePerNode: false
  nodeConfigGroups:
    default_group:
      isNode: true
      storageConfigs:
        - mountPath: "/opt/nifi/nifi-current/logs"
          name: logs
          pvcSpec:
            accessModes:
              - ReadWriteOnce
            storageClassName: "standard"
            resources:
              requests:
                storage: 10Gi
      serviceAccountName: "default"
      resourcesRequirements:
        limits:
          cpu: "2"
          memory: 3Gi
        requests:
          cpu: "1"
          memory: 1Gi
  nodes:
    - id: 1
      nodeConfigGroup: "default_group"
    - id: 2
      nodeConfigGroup: "default_group"
  propagateLabels: true
  nifiClusterTaskSpec:
    retryDurationMinutes: 10
  listenersConfig:
    internalListeners:
      - type: "http"
        name: "http"
        containerPort: 8080
      - type: "cluster"
        name: "cluster"
        containerPort: 6007
      - type: "s2s"
        name: "s2s"
        containerPort: 10000
  externalServices:
    - name: "clusterip"
      spec:
        type: ClusterIP
        portConfigs:
          - port: 8080
            internalListenerName: "http"
      serviceAnnotations:
        toto: tata
```

## NifiCluster

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|metadata|[ObjectMetadata](https://godoc.org/k8s.io/apimachinery/pkg/apis/meta/v1#ObjectMeta)|is metadata that all persisted resources must have, which includes all objects users must create.|No|nil|
|spec|[NifiClusterSpec](#nificlusterspec)|defines the desired state of NifiCluster.|No|nil|
|status|[NifiClusterStatus](#nificlusterstatus)|defines the observed state of NifiCluster.|No|nil|

## NifiClusterSpec

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|service|[ServicePolicy](#servicepolicy)| defines the policy for services owned by NiFiKop operator. |No| - |
|pod|[PodPolicy](#podpolicy)| defines the policy for pod owned by NiFiKop operator. |No| - |
|zkAddress|string| specifies the ZooKeeper connection string in the form hostname:port where host and port are those of a Zookeeper server.|Yes|""|
|zkPath|string| specifies the Zookeeper chroot path as part of its Zookeeper connection string which puts its data under same path in the global ZooKeeper namespace.|Yes|"/"|
|initContainerImage|string|  can override the default image used into the init container to check if ZoooKeeper server is reachable.. |Yes|"busybox"|
|initContainers|\[ \]string| defines additional initContainers configurations. |No|\[ \]|
|clusterImage|string| can specify the whole nificluster image in one place. |No|""|
|oneNifiNodePerNode|boolean|if set to true every nifi node is started on a new node, if there is not enough node to do that it will stay in pending state. If set to false the operator also tries to schedule the nifi node to a unique node but if the node number is insufficient the nifi node will be scheduled to a node where a nifi node is already running.|Yes| nil |
|propagateLabels|boolean| - |Yes|false|
|managedAdminUsers|\[  \][ManagedUser](#managedusers)| contains the list of users that will be added to the managed admin group (with all rights). |No|[]|
|managedReaderUsers|\[  \][ManagedUser](#managedusers)| contains the list of users that will be added to the managed admin group (with all rights). |No|[]|
|readOnlyConfig|[ReadOnlyConfig](./2_read_only_config.md)| specifies the read-only type Nifi config cluster wide, all theses will be merged with node specified readOnly configurations, so it can be overwritten per node.|No| nil |
|nodeConfigGroups|map\[string\][NodeConfig](./3_node_config.md)| specifies multiple node configs with unique name|No| nil |
|nodes|\[  \][Node](./3_node_config.md)| specifies the list of cluster nodes, all node requires an image, unique id, and storageConfigs settings|Yes| nil 
|disruptionBudget|[DisruptionBudget](#disruptionbudget)| defines the configuration for PodDisruptionBudget.|No| nil |
|ldapConfiguration|[LdapConfiguration](#ldapconfiguration)| specifies the configuration if you want to use LDAP.|No| nil |
|nifiClusterTaskSpec|[NifiClusterTaskSpec](#nificlustertaskspec)| specifies the configuration of the nifi cluster Tasks.|No| nil |
|listenersConfig|[ListenersConfig](./6_listeners_config.md)| specifies nifi's listener specifig configs.|Yes| - |
|sidecarConfigs|\[  \][Container](https://godoc.org/k8s.io/api/core/v1#Container)|Defines additional sidecar configurations. [Check documentation for more informations]|
|externalServices|\[  \][ExternalServiceConfigs](./7_external_service_config.md)| specifies settings required to access nifi externally.|No| - |

## NifiClusterStatus

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|nodesState|map\[string\][NodeState](./5_node_state.md)|Store the state of each nifi node.|No| - |
|State|[ClusterState](#clusterstate)|Store the state of each nifi node.|Yes| - |
|rootProcessGroupId|string|contains the uuid of the root process group for this cluster.|No| - |

## ServicePolicy

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|headlessEnabled|boolean| specifies if the cluster should use headlessService for Nifi or individual services using service per nodes may come an handy case of service mesh.|Yes|false|
|annotations|map\[string\]string|Annotations specifies the annotations to attach to services the NiFiKop operator creates|No|-|

## PodPolicy

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|annotations|map\[string\]string|Annotations specifies the annotations to attach to pods the NiFiKop operator creates|No|-|

## ManagedUsers

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|create|bool| if set to true, will create a podDisruptionBudget.|No| - |
|name|string|name field is use to name the NifiUser resource, if not identity is provided it will be used to name the user on NiFi cluster side.|Yes| - |

## DisruptionBudget

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|identity|string|identity field is use to define the user identity on NiFi cluster side, it use full when the user's name doesn't suite with Kubernetes resource name.|No| - |
|budget|string| the budget to set for the PDB, can either be static number or a percentage.|Yes| - |

## LdapConfiguration

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|enabled|boolean| if set to true, we will enable ldap usage into nifi.properties configuration.|No| false |
|managerDN|string| The DN of the manager that is used to bind to the LDAP server to search for users|No| "" |
[managerPassword|string| The password of the manager that is used to bind to the LDAP server to search for users|No| "" |
|authStrategy|string|How the connection to the LDAP server is authenticated. Possible values are ANONYMOUS, SIMPLE, LDAPS, or START_TLS.|No|"START_TLS"|
|tls|LdapTlsConfig(#TlsLdapConfig)|TLS Configuration|No| nil |
|referralStrategy|string|Strategy for handling referrals. Possible values are FOLLOW, IGNORE, THROW.|No|"FOLLOW"|
|connectTimeout|int|Duration of connect timeout (secs).|No|10|
|readTimeout|int|Duration of read timeout (secs).|No|10|
|url|string| space-separated list of URLs of the LDAP servers (i.e. ldap://${hostname}:${port}).|No| "" |
|identityStrategy|String|Strategy to identify users. Possible values are USE_DN and USE_USERNAME. USE_DN will use the full DN of the user entry if possible. USE_USERNAME will use the username the user logged in with.|No|"USE_DN"|
|authExpiration|int|The duration of how long the user authentication is valid for. If the user never logs out, they will be required to log back in following this duration (hours)|No|12|
|pageSize|int|The page size when retrieving users and groups. If not specified, no paging is performed.|No|""|
|searchBase|string| base DN for searching for users (i.e. CN=Users,DC=example,DC=com).|No| "" |
|searchFilter|string| Filter for searching for users against the 'User Search Base'. (i.e. sAMAccountName={0}). The user specified name is inserted into '{0}'.|No| "" |
|ldapSync|bool|If set to true, nifi will sync users and group from ldap database|No| false |
|userSync|LdapSyncSpec(#LdapSyncSpec User)|Ldap User Sync Spec|Yes| nil |
|groupSync|LdapSyncSpec(#LdapSyncSpec Group)|Ldap Group Sync Spec|Yes| nil |

## LdapTlsConfig

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|keystore|LdapKeystore(#LdapKeystore)|TLS LDAP Keystore|No| nil |
|clientAuth|string|Client authentication policy when connecting to LDAP using LDAPS or START_TLS. Possible values are REQUIRED, WANT, NONE|No|NONE|
|protocol|string|Protocol to use when connecting to LDAP using LDAPS or START_TLS. (i.e. TLS, TLSv1.1, TLSv1.2, etc).|No| "TLS" |
|shutdownGracefully|bool|Specifies whether the TLS should be shut down gracefully before the target context is closed. Defaults to false.|No|false|

## LdapKeystore

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|secretName|string|SecretName should contain ca certs|Yes| "" |
|password|string|Password for the Keystore and Truststore that is used when connecting to LDAP using LDAPS or START_TLS.|Yes| "" |
|type|string|Type of the Keystore and Truststore that is used when connecting to LDAP using LDAPS or START_TLS (i.e. JKS or PKCS12).|No| "JKS" |

## LdapSyncSpec User

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|searchBase|string|Base DN for searching for users (i.e. ou=users,o=nifi). Required to search users.|Yes| "" |
|searchFilter|string|Filter for searching for users against the 'User Search Base'.|No| "" |
|searchScope|string|Search scope for searching users (ONE_LEVEL, OBJECT, or SUBTREE). Required if searching users.|No|"ONE_LEVEL"|
|objectClass|string|Object class for identifying users (i.e. person). Required if searching users.|No| "person" |
|nameAttr|string|Attribute to use to extract user identity (i.e. cn). Optional. If not set, the entire DN is used.|No|""|
|groupAttr|string|Attribute to use to define group membership (i.e. memberof). Optional. If not set group membership will not be calculated through the users. Will rely on group membership being defined through Group Member Attribute if set. The value of this property is the name of the attribute in the user ldap entry that associates them with a group. The value of that user attribute could bea dn or group name for instance. What value is expected is configured in the User Group Name Attribute - Referenced Group Attribute. |No| "" |
|referencedAttr|string|If blank, the value of the attribute defined in User Group Name Attribute is expected to be the full dn of the group. If not blank, this property will define the attribute of the group ldap entry that the value of the attribute defined in User Group Name Attribute is referencing (i.e. name). Use of this property requires that Group Search Base is also configured.|No| "" |

## LdapSyncSpec Group

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|searchBase|string|Base DN for searching for groups (i.e. ou=groups,o=nifi). Required to search groups.|Yes| "" |
|searchFilter|string|Filter for searching for groups against the 'Group Search Base'.|No| "" |
|searchScope|string|Search scope for searching groups (ONE_LEVEL, OBJECT, or SUBTREE). Required if searching groups.||"ONE_LEVEL"|
|objectClass|string|Object class for identifying groups (i.e. groupOfNames). Required if searching groups.|No|"group"|
|nameAttr|string|Attribute to use to extract group name (i.e. cn). Optional. If not set, the entire DN is used.|No| |
|groupAttr|string|Attribute to use to define group membership (i.e. member). Optional. If not set group membership will not be calculated through the groups. Will rely on group membership being defined through 'User Group Name Attribute' if set. The value of this property is the name of the attribute in the group ldap entry that associates them with a user. The value of that group attribute could be a dn or memberUid for instance. What value is expected is configured in the 'Group Member Attribute - Referenced User Attribute'. (i.e. member: cn=User 1,ou=users,o=nifi vs. memberUid: user1)|No| "" |
|referencedAttr|string|If blank, the value of the attribute defined in 'Group Member Attribute' is expected to be the full dn of the user. If not blank, this property will define the attribute of the user ldap entry that the value of the attribute defined in 'Group Member Attribute' is referencing (i.e. uid). Use of this property requires that 'User Search Base' is also configured. (i.e. member: cn=User 1,ou=users,o=nifi vs. memberUid: user1)|No| "" |

## NifiClusterTaskSpec

|Field|Type|Description|Required|Default|
|-----|----|-----------|--------|--------|
|retryDurationMinutes|int| describes the amount of time the Operator waits for the task.|Yes| 5 |

## ClusterState

|Name|Value|Description|
|-----|----|------------|
|NifiClusterInitializing|ClusterInitializing|states that the cluster is still in initializing stage|
|NifiClusterInitialized|ClusterInitialized|states that the cluster is initialized|
|NifiClusterReconciling|ClusterReconciling|states that the cluster is still in reconciling stage|
|NifiClusterRollingUpgrading|ClusterRollingUpgrading|states that the cluster is rolling upgrading|
|NifiClusterRunning|ClusterRunning|states that the cluster is in running state|

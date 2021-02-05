# NifiKop Tanka templates
> This is an early stage

Follow the [tanka documentation](https://tanka.dev/) for setting up your environment.

Example `main.jsonnet` file:
```jsonnet
(import 'github.com/orangeopensource/nifikop/ksonnet/operator/main.libsonnet')
+ {
  _config+:: {
    namespace: 'my-namespace',
  },
}
```

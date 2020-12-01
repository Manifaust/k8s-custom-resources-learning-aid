Kubebuilder talks explains
* What are CRDs, types, and controllers
* How do you start building CRDs and controllers quickly

# Prerequisites
* Install [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/)

# Debugging
If you run into issues compare your code with the [book's example](https://github.com/kubernetes-sigs/kubebuilder/tree/book-v2/docs/book/src/cronjob-tutorial/testdata/project) on GitHub. Make sure you're on the `book-v2` (or whatever version you're reading) branch when looking at GitHub.

When looking at those examples, you need to remember to substitute `tutorial.kubebuilder.io/project/api/v1` with your package.

# Section 1.9 bugs
In section 1.9 of the book, we're shown how to deploy the CronJob example. But this chapter has some bugs or out of date issues.

## k8s API changes
The book asks you to run `make install` with the CronJob example. However, because of changes in kubernetes, that won't work. You might run into a failure like:
```sh
* spec.validation.openAPIV3Schema.properties[spec]...properties[protocol].default: Required value: this property is in x-kubernetes-list-map-keys, so it must have a default or be a required property
make: *** [install] Error 1

```

This bug is caused by changes in the k8s API which necessitates changes in controller-tools. There' an issue and [workaround](https://github.com/kubernetes-sigs/kubebuilder/issues/1466#issuecomment-712444882) in the kubebuilder git that you need to implement to get pass this error.

## Disabling webhooks
The book recommends disabling headhooks for local testing by running `make run ENABLE_WEBHOOKS=false` but that's not enough. You might run into an error like this:
```sh
ERROR	setup	problem running manager	{"error": ".../k8s-webhook-server/serving-certs/tls.crt: no such file or directory"}
```

The temporary workaround is in a [PR](https://github.com/kubernetes-sigs/kubebuilder/pull/1862/files). It invovles modifying `main.go`.
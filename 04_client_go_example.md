# client-go
`client-go` is a tool for accessing the k8s api in golang. We can use it to get information about pods, deployments, and any other types of resources.

These examples assume you're running `minikube` locally, but the `client-go` repo has instructions to connect to other clusters.

## Download
The quickest way to try client go is to clone the `client-go` repo and try the examples.

```sh
$ git clone https://github.com/kubernetes/client-go
$ cd client-go
```

## Authentication
By default, `client-go` uses the current context in the kubeconfig to access your cluster. You can use `kubectl` to check whether your current context is configured to point to a cluster.

```sh
$ kubectl get pods
NAME                          READY   STATUS    RESTARTS   AGE
hello-node-7567d9fdc9-w8xb8   1/1     Running   1          60m
```

Use the *out-of-cluster* example to test that `client-go` can authenticate with your cluster.

```sh
$ cd examples/out-of-cluster-client-configuration
$ go run main.go
There are 10 pods in the cluster
Pod example-xxxxx in namespace default not found
There are 10 pods in the cluster
Pod example-xxxxx in namespace default not found
There are 10 pods in the cluster
...
```

## Basic usage
Out of the box, the code inside the example `main.go` looks for a pod named `example-xxxx`. Lets modify the code to look for a pod that's already deployed in our cluster. We'll use the pod `hello-node-7567d9fdc9-w8xb8` which we saw when we ran `kubectl get pods` earlier, but your pod name will be different. After changing the pod name in `main.go`, run the example again and it should find it.

```sh
% go run main.go
There are 10 pods in the cluster
Found pod hello-node-7567d9fdc9-w8xb8 in namespace default
```

## Exercising the k8s API
Next we'll try a different demo which creates, updates, and deletes a deployment using `client-go`. When you run the example, it will:
1. Create a deployment `demo-deployment` with 2 replicas using image `nginx:1.12`
2. Update the deployment to have 1 replica, and switch to using image `nginx:1.13`
3. List deployments
4. Delete the deployment

The example will wait for user input between each step.

```sh
% cd ../create-update-delete-deployment
% go run main.go
Creating deployment...
Created deployment "demo-deployment".
-> Press Return key to continue.

Updating deployment...
Updated deployment...
-> Press Return key to continue.

Listing deployments in namespace "default":
 * demo-deployment (1 replicas)
 * hello-node (1 replicas)
-> Press Return key to continue.

Deleting deployment...
Deleted deployment.
```

In between each step, we can use `kubectl describe` to confirm the changes we expect are happening.

```sh
% kubectl describe deployment demo
Name:                   demo-deployment
Namespace:              default
...
Replicas:               2 desired | 2 updated | 2 total | 2 available | 0 unavailable
...
Pod Template:
  Labels:  app=demo
  Containers:
   web:
    Image:        nginx:1.12
...
```


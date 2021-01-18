_This chapter points you introduces Tekton and points you to the tutorial. I mostly just list out all gotchas I found while going through the tutorial._

# A k8s native CI/CD
Tekton is a k8s native CI/CD tool. It is a good example of providing CustomResourceDefinitions as a way of extending the k8s API to provide useful functionality.

To start playing around with Tekton, users just have to install their CRDs and controllers, which is very easy to do. These CRDs allow the user to create Tekton objects such pipelines and tasks using `kubectl apply` and yaml files. There is a Tekton CLI `tkn` but users don't need to rely on that.

Since it was built from the ground up to rely on the k8s platform, Tekton can rely on existing k8s concepts instead of inventing its own:
* It doesn't provide its own auth and roles implementation, it relies on the concepts roles and rolebindings that already exist in k8s.
* It doesn't have its own secrets implementation.
* Tekton pipelines and tasks execute by creating deployments and pods, it doesn't need to invent its own worker concept.

# Getting familiar with Tekton
Going through the [Tekon tutorial](https://github.com/tektoncd/pipeline/blob/master/docs/tutorial.md) is a good way to try out the experience of installing CRDs and creating custom resources. The tutorial has some outdated elements, so the rest of this doc will highlight those issues so the tutorial can go more smoothly.

The following sections cover:
1. Concepts
2. Debugging
3. Tutorial

# Tekton Concepts
Read this [IBM article](https://developer.ibm.com/devpractices/devops/articles/introduction-to-tekton-architecture-and-design/) about the components and concepts of Tekton.

# Debugging
When you create a *taskrun* or *pipelinerun* in Tekton you can use the Tekton CLI to see the logs. When something goes wrong, usually knowing which stage of the task that's broken and seeing its logs will give you a big hint.

```sh
tkn taskrun describe my-task-run
tkn taskrun logs my-task-run
```

From my experience the most common issues that a novice will run into with this tutorial invovles accessing the image registry and providing Tekton the permissions to take action.

# Tutorial
Try to go through the [Tekon tutorial](https://github.com/tektoncd/pipeline/blob/master/docs/tutorial.md) from top to bottom. It walks through installing Tekon, setting up tasks, and creating pipelines.

## Create dockerhub secret
The example task will push to a registry. It's easiest to use your Dockerhub account.

```sh
kubectl create secret \
  docker-registry \
  regcred \
  --docker-server=https://index.docker.io/v1/ \
  --docker-username=<DOCKERHUB_USERNAME> \
  --docker-password=<DOCKERHUB_PASSWORD> \
  --docker-email <DOCKERHUB_EMAIL>
```

## Creating tasks
In the section of the tutorial where they ask you to create a git-resource, change the `revision` property to `v0.32.0` because the Dockerfile in `master` doesn't build as easily anymore.

```yml
apiVersion: tekton.dev/v1alpha1
kind: PipelineResource
metadata:
  name: skaffold-git
spec:
  type: git
  params:
    - name: revision
      value: v0.32.0 # this is different than example
    - name: url
      value: https://github.com/GoogleContainerTools/skaffold
```

## Creating Pipelines
In the part where it asks you to create a `clusterrole`, in order for the script to work in `zsh`, instead of using `--verb=*` you should use `--verb="*"`. Also, add `services` as one of the list of resources this clusterrole applies to, like this:

```sh
kubectl create clusterrole tutorial-role \
  --verb="*" \
  --resource=deployments,deployments.apps,services
```

Another edit you have to do: when the tutorial asks you to configure a new task named `deploy-using-kubectl`, you have to pass `validate=true` to the call to `kubectl`, like this:

```yml
apiVersion: tekton.dev/v1beta1
kind: Task
metadata:
  name: deploy-using-kubectl
...
  steps:
...
    - name: run-kubectl
      image: lachlanevenson/k8s-kubectl
      command: ["kubectl"]
      args:
        - "apply"
        - "-f"
        - "$(params.path)"
        - "validate=false" 
```

## Delete Your Pipelines and Tasks
To delete all the k8s objects you creating, run `kubectl delete -f` for all that files you ran `kubectl apply -f` with.

## Delete Tekton
```sh
$ kubectl delete --filename https://storage.googleapis.com/tekton-releases/pipeline/latest/release.yaml
```

# Exercises
* Complete the [kpack tutorial](https://github.com/pivotal/kpack/blob/master/docs/tutorial.md).
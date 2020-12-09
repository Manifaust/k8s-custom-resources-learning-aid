# Custom Resources and Custom Controllers
Custom resources and custom controllers are concepts in k8s that allow developers to extend the k8s and provide those extensions to users. You can think.  

## Custom Resource Definitions (CRDs) and Custom Resources
Developers are not limited to using defaults k8s types such as pods, deployments, and replicasets in their k8s clusters. They can define their own types such using custom resource definitions (CRDs) and provide them to users to install. Once installed, users can then create custom resources from CRDs. Creates custom resources are stored in the k8s cluster next to regular resources.

Soon we'll see some examples from Tekton, which is a project that provides CRDs for CI/CD concepts such as tasks and pipelines. After installing CRDs, users can create their own instances of tasks and pipelines and custom them.

As another example, say we want to create a CRD and controller that allow users to receive weather reports every day. Our CRD will declare which fields are required for creating such reports, such as weather location, time of reminder, and destination email address to send the reminder email. After the user installs our CRD they can write a custom resource yaml:

```yaml
# my-weather-reminder.yml
apiVersion: weather.com/v1alpha1
kind: Reminder
metadata:
  name: reminder-toronto
spec:
  city: Toronto, ON
  time: '0700'
  email: foo@example.com

```

Then they can run `kubectl apply -f my-weather-reminder.yml` to create that custom resource. Once created, a `Reminder` object will now exist inside the cluster and can be searched for, examined, and deleted just like other k8s objects.

To use a programming language analogy, think of a CRD as the type declaration/definition, and custom resources as instances of that type. In k8s, types belong to *groups*, are *versioned*, and have *names*, so often times they're referred to as GroupVersionKind or GVK. The programming analogy only goes so far because while custom resources can hold properties and state and user intent, it doesn't have any logic. The actual work happens inside controllers. 

## Controllers
A controller is a program that run in the background, paying attention to the creation/update/deletion of custom resources of a specific CRD and takes action based on what it sees. While it's in charge of one type of resource, it can access the entire k8s API and manipulate any kind of resources.

Using our weather reminder example again, to support the `Reminder` CRD, we need to create a controller that pay attention to the existance of user created `Reminder` resources. We might implement it to run periodically to list all the `Reminder` resources and see if it's the right time to check the weather for one of them. When the time comes, it will create a pod to run a script to query the weather and send an email.

# Further Concepts
## Declarative
[k8s concepts](https://kube.academy/lessons/kubernetes-concepts)

## Spec and Status
From kubebuilder book:
> It’s a controller’s job to ensure that, for any given object, the actual state of the world (both the cluster state, and potentially external state like running containers for Kubelet or loadbalancers for a cloud provider) matches the desired state in the object. [...] We call this process reconciling.

## Reconcilliation
[What is reconcilliation](https://speakerdeck.com/thockin/kubernetes-what-is-reconciliation)

# Kubernetes API

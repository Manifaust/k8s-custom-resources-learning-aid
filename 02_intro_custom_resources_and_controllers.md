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

# Recurring Concepts
At this point, it's worth talking about certain related concepts/patterns that come up again and again in k8s. You'll be expect to recognize and understand these patterns when you're using using k8s software. Furthermore, you'll need to implement these patterns when developing k8s native resources and controllers. These patterns are baked into how k8s itself works and plays a large part in the user's mental model when they approach a new piece of k8s software. 

## Declarative vs. Imperative
When a k8s user wants to take action, whether it's the creation of a deployment of the running of a Tekton task, the convention they follow should not be to write a script filled with commands to execute (i.e. the imperative model).

They expect to follow the common k8s pattern of writing a yaml file (custom resource) that describes the object or task they want to create (i.e. the declarative model). They just want to state their intention, or end result, so as much as possible the custom resource and controller should abstract away all the steps it takes to realize the user's vision.

This is one reason why there's so much yaml when working with k8s.

## Spec and Status
The concept of user intent appears when you design custom resource definitions as well. As you'll see later, custom types have have many properties and conventions, one of which is the separation of two properties: *spec* and *status*.

Spec is the property of a custom resource that holds user intent. It's basically all the fields that the user configured, which represents the future state they want.

Status is the property that holds the current state of the system. It's supposed to be continuously updated by the controller, and can be read by the user. By reading the state of a custom resource, the user (and other controllers) knows whether an object or action has started, failed, or is in progress. 

Corollary to that point, it's perfectly normal that there exists some time delay between before the user expressing their intent and the controller (and other processes) completing the necessary actions to realize that intent. This asynchronous, non-blocking way of working is how k8s resources are expected to be implemented. We can contrast this to some traditional *API-driven* designs where once the user sends a request, the server may not send a response until the relevant transaction is done and the user's desire is met.

## Reconcilliation
Now we're ready to describe what a controller loop looks like:
1. List the custom resources the controller is responsible for.
2. Check what the spec (user intention) of those resources.
3. Check the status (current state) of those resources.
4. Take the appropriate action to move the current state closer to the desired state.
5. Update the current state.
6. Repeat

This is pattern is called **reconcilliation** and this type of loop appears in many places in k8s. You will be expected to recognize this pattern when you see it, and when you design your own custom controller, you'll be expected to implement this pattern.

Check out the slide deck [What is reconcilliation](https://speakerdeck.com/thockin/kubernetes-what-is-reconciliation) a more detailed illustration.

# Implementing a Controller
So far we've talked abstractly about the controller taking action and executing commands. But the only concrete examples of kicking off tasks have all involved using `kubectl` on the terminal.

We won't talk about it in detail here, but I just want to finish this section by saying that controllers are just programs (Golang, Ruby, etc. applications) you write that run in the background somewhere. The first controller tutorials I reference later on will get you running an example Golang-written controller in the terminal, and hooking onto a local minikube. In production, the controller will be running in the target k8s cluster itself.

There are Golang packages that providing helpful functionality focused on implementing controllers. Your custom controller code is expected to use these packages and implement certain functions and interfaces from these packages. Furthermore, k8s has a broad API where you can send http CRUD requests to take specific action. In a later section we'll use `client-go` to demonstrate how easy it is to use the API, and it's one of several packages that help you take action or retrieve information. Not only will you be sending the API requests, the controller you implement will respond to requests as well.

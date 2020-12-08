*The basic objective of custom resources and controllers is to extend the k8s API (and provide these extensions to the user). Here are my thoughts about what that means and what it's important.*

# Extending the k8s API
Extending the k8s API means doing more than just containerizing an app, and creating a deployment for that app. It includes:
* Creating your own k8s object and concepts. k8s comes with a default set of objects like pods and replicasets, but you can create your own.
* Using the k8s API to understand what objects exists in the cluster, to observe the state of the system, and to take action based on that information.
* Using the k8s API to make changes to the system, by creating objects, updating configuration.

# Why is it Important to Learn?
The reason we want to learn to extend the k8s API is because it's an important part of building products on top of k8s. It's about what's the best way to reach customers, to take advantage of efficiencies and patterns in technology, and talking the same language as people we want to work with.

## k8s as a Common Platform
Eventually, every body will be running k8s, and a lot of people will be building for k8s, from customers, to vendors, to other teams with collaborate with. Since they are running k8s, we know they have a platform that support fundamental cloud native dependencies like secrets, running/scaling workloads, service discovery, and so on. Because these building blocks are there, we can take advantage of them instead of building everything from scratch, if we're willing to work in a k8s native way.

## Customers want k8s
Often times executives and VPs running enterprise companies make decisions about where they'll spend money based on where they think techology are heading. If they know that we're leaders in the area/domain that they're headed, they'll more likely ask for our opinions about how to solve difficult problem, and how technology can help.

A lot of our customer are investing in k8s, so they want to know that our products will work nicely with the ecosystem they're building. By letting them know that the products we're buliding are designed to be k8s native, it increases their confidence in using those products.

## k8s is a Business Priority
Because of the above points, it's a priority to use k8s as the platform on which we build our own platform. That means more than just deploying our applications on top of k8s. It also means:
* Understanding the k8s operator and developer personas and the problem they run into.
* Knowing the gaps in the k8s experience and where we can help create solutions.
* Understand the OSS and commerical k8s ecosystem.
* Allowing customers and other teams to use your software in a way that follow k8s conventions.

In a k8s-focused organization, it benefits everyone to collaborate and communicate at a higher abstraction level. You will be expected to have a good understanding of the k8s building blocks and boundaries and interfaces. Even if you're not using these patterns, you still need to be knowledgeable enough to form opinions on the topic in order to influence the direction your product is going, and to know how best to work with other teams.

# What does "k8s native" mean?
I mentioned several times that it's important that we build software in a k8s native way (and this study aid helps you learn how to do that), I didn't explain what that actually means. Well *k8s native* is a made up phrase that doesn't have an exact definition. In general, when people say Kubernetes native, they mean the relevant software:
* Is aware that it's running inside k8s and so takes advantage of the k8s API and libraries. It's not just a generic containerized app or service.
* Follows k8s community conventions and patterns when it comes to how it is deployed, configured, and how it works.
* Supports interoperability with the k8s ecosystem, including OSS projects, and other k8s native projects.
* Utilizes k8s bulding blocks like pods and deployments instead of reinventing the wheel.

In general, successful k8s software tends to understand and play nicely with the existing community. The phrase *k8s native* might only be useful as a contrast against a *incongruous* type of software that is inflexible, and wants to cram existing ways of working without fitting into a new model of work.

It's worth pointing out that the community is growing and evolving too, so new patterns will emerge to supplant old patterns. In the mean time, it's useful to learn what the existing conventions are, and hopefully this learning aid help you do that.

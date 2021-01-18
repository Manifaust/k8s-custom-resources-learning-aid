_This section guides you through making a very basic controller using kubebuilder to make things easier._

# Kubebuilder
[Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) is a framework for creating k8s controllers. It includes libraries and scripts to get people up and running quickly. The project is maintained by the k8s SIG API Machinery.

There are other frameworks too, such as [Operator SDK](https://github.com/operator-framework/operator-sdk) (which shares the same underlying library [controller-runtime](https://github.com/kubernetes-sigs/controller-runtime)). You certainly don't need kubebuilder to create a controller but we're focusing on it because it's quicker to get started and it's maintained by a k8s SIG. 

Developers are still figuring out how to make implementing controller easier so these frameworks are not perfectly. You'll find that some tools and processes are a bit clunky, and you'll likely run into more boilerplate code than you want. It's all because this is such a new area of growth in k8s development, and it's part of the fun!

# Install Kubebuilder
First install [kustomize](https://kubectl.docs.kubernetes.io/installation/kustomize/), which a dependency of kubebuilder.

After that, use the snippet below to install kubebuilder. At the time of writing, the newest version of kubebuilder is `3.0.0-alpha.1`. Version 3 is new enough that the kubebuilder book - which is the focus of the next chapter - hasn't been updated to incorporate it. For this reason we're sticking with an older version `2.3.1`. Once the book is updated, or you become very comfortable with kubebuilder, then you should read the migration docs for version 3.

```bash
os=$(go env GOOS)
arch=$(go env GOARCH)

# download kubebuilder and extract it to tmp
curl -L "https://go.kubebuilder.io/dl/2.3.1/${os}/${arch}" | tar -xz -C /tmp/

# move to a long-term location and put it on your path
# (you'll need to set the KUBEBUILDER_ASSETS env var if you put it somewhere else)
sudo mv "/tmp/kubebuilder_2.3.1_${os}_${arch}" /usr/local/kubebuilder
export PATH=$PATH:/usr/local/kubebuilder/bin

# Remember to modify your PATH in .bashrc, .bash_profile, or .zshrc 
```

# Initialize Your Project

When you initialize a kubebuilder project, the framewok will create an entire golang application for you with many files ready to be modified. The `init` command take several flags, for now we'll specify a _domain_ which makes up part of the _group_ name of your CRD. It will also need a _repo_, which becomes the golang module path for your new project.

```bash
$ mkdir my-weather-app-crd
$ cd my-weather-app-crd
$ kubebuilder init --domain example.com --repo github.com/Manifaust/k8s-custom-resources-learning-aid/examples/weather-app
```

If you browse the folders and files of the resulting golang project, you'll see it has not created a CRD nor controller for you yet. But everything is set up for one to be created. One thing you'll see is that the `go.mod` includes libraries from `apimachinery`, `client-go`, and `controller-runtime`. These libraries makes it easier to implement controllers and CRDs.

```bash
$ less config/rbac/role_binding.yaml
module github.com/Manifaust/k8s-custom-resources-learning-aid/examples/weather-app

go 1.13

require (
        k8s.io/apimachinery v0.17.2
        k8s.io/client-go v0.17.2
        sigs.k8s.io/controller-runtime v0.5.0
)
```

To create a CRD and controller, use the `kubebuilder create` command. Here, `weather-app` makes up part of the group name (combined with the *domain* from earlier). `v1alpha1` is our CRD's version. And `CheckWeather` is the CRD's *kind*. You'll learn more about these concepts in the kubebuilder book.

```bash
$ kubebuilder create api \
              --group weather-app \
              --version v1alpha1 \
              --kind CheckWeather
Create Resource [y/n]
y
Create Controller [y/n]
y
Writing scaffold for you to edit...
api/v1alpha1/checkweather_types.go
controllers/checkweather_controller.go
...
```

Here are some of the files that's been generated:
## `api/v1alpha1/checkweather_types.go`
This file defines properties of the _spec_ and _status_ of your CRD. You'll customize this file later. For now, just take note the that spec contains only one field `Foo`, and the status contains no fields yet.

## `config/samples/weather-app_v1alpha1_checkweather.yaml`
```bash
$ cat config/samples/weather-app_v1alpha1_checkweather.yaml
apiVersion: weather-app.example.com/v1alpha1
kind: CheckWeather
metadata:
  name: checkweather-sample
spec:
  # Add fields here
  foo: bar
```

An example custom resource based on your CRD. This is just like any other resource yaml such as those for pods, replicasets, or deployments. Your users will write and apply yamls like this to create custom resources based on your definition. Like any other resource yaml, the `apiVersion` property specifies the desired group, in this case `weather-app.example.com`, and the desired version, in this case `v1alpha1`. The `kind` of our CRD is `CheckWeather`, which is what we provided in `kubebuilder create` command. The `spec` has the `foo` property as declared in `checkweather_types.go`. When we add more spec fields in `checkweather_types.go`, users will be able to add more fields to the spec portion of their resource yaml.

## `controllers/checkweather_controller.go`
This is our controller, it's where we'll check on all our resources and implement reconcilliation (described in an earlier chapter). The `Reconcile` function is where all the action starts, it gets called any time relevant changes occur. We'll be modifying this function later to add custom functionality.

# Installing your CRD
Now we play the role of admin and install the CRD. We'll install it in a custom namespace as not to distrupt other parts of the cluster.

```bash
$ kubectl create namespace weather-testing
$ kubectl config set-context $(kubectl config current-context) --namespace=weather-testing

$ make install
go: creating new go.mod: module tmp
go: found sigs.k8s.io/controller-tools/cmd/controller-gen in sigs.k8s.io/controller-tools v0.2.5
/Users/tony/workspace/go/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
kustomize build config/crd | kubectl apply -f -
customresourcedefinition.apiextensions.k8s.io/checkweathers.weather-app.example.com configured

$ kubectl get crds
NAME                                    CREATED AT
checkweathers.weather-app.example.com   2021-01-16T23:59:38Z
```

If you have a custom kubeconfig that's not in the `$HOME/kube` then you can use the `KUBECONFIG` environment variable to tell kubectl where your kubeconfig is.

```bash
$ KUBECONFIG=<my_kubeconfig_location> make install
```

# Run the Controller
In a real production cluster, your controller is deployed in a regular k8s deployment. During our testing however, we'll instead just run the controller locally on our computer as a normal golang app.

```bash
$ make run
go: creating new go.mod: module tmp
go: found sigs.k8s.io/controller-tools/cmd/controller-gen in sigs.k8s.io/controller-tools v0.2.5
/Users/tony/workspace/go/bin/controller-gen object:headerFile="hack/boilerplate.go.txt" paths="./..."
go fmt ./...
go vet ./...
/Users/tony/workspace/go/bin/controller-gen "crd:trivialVersions=true" rbac:roleName=manager-role webhook paths="./..." output:crd:artifacts:config=config/crd/bases
go run ./main.go
2021-01-18T11:19:23.345-0500	INFO	controller-runtime.metrics	metrics server is starting to listen	{"addr": ":8080"}
...
```

# Create a Custom Resource
Now that the namespace has your CRD installed, and your controller is running, it's time to play the role of the user and create a custom resource.

Open another terminal and use `kubectl apply` with the custom resource example we pointed to earlier.

```bash
$ kubectl apply -f config/samples/weather-app_v1alpha1_checkweather.yaml
checkweather.weather-app.example.com/checkweather-sample created
$ kubectl describe CheckWeather checkweather-sample
Name:         checkweather-sample
Namespace:    weather-testing
...
API Version:  weather-app.example.com/v1alpha1
Kind:         CheckWeather
...
Spec:
  Foo:   bar
```

You created a custom resource!

Back in the controller terminal screen, you'll see this log to indicate that the controller encountered the new resource:

```bash
...
2021-01-16T19:22:38.245-0500	DEBUG	controller-runtime.controller	Successfully Reconciled	{"controller": "checkweather", "request": "weather-testing/checkweather-sample"}
```

# Customize your Controller
Here're the steps we're going to take to create custom functionality:
1. Add a new field to the CRD's spec to accept a city location. To do that we need to modify `checkweather_types.go`.
2. Update our controller to read the custom resource's _spec_, look for the new city field, and to call the weather API. We will modify `checkweather_controller.go`.
3. Once the response from the API comes back with weather info, we'll update the custom resource's _status_.
4. Update our example resource yaml to add an example city to test out.

# Support more Spec and Status Fields
Open up `api/v1alpha1/checkweather_types.go` and add the new city field by modifying the `CheckWeatherSpec` struct:

```go
type CheckWeatherSpec struct {
  City string `json:"city,omitempty"`
}
```

Also, add constants to represent the current state of our request:
```go
const (
  StatePending  = "PENDING"
  StateFinished = "FINISHED"
)
```

A custom resource's status is how the user will know when the request is complete and what the result will be. Modify the `CheckWeatherStatus` struct to look support the new state and temperature fields. Also we'll need to add the _marker_ `+kubebuilder:subresource:status` to help the framework generate a status _subresource_ for our resource:
```go
type CheckWeatherStatus struct {
  State       string `json:"state,omitempty"`
  Temperature int32  `json:"temperature"`
}

// +kubebuilder:subresource:status
```

Side note: You might have noticed that there are a lot of comments like the one above in the files that kubebuilder provides. Kubebuilder relies on some generation to provide its functionality. Comment _markers_ tell the code generator what to do.

If this is too confusing check out the end result in the [example from the learning-aid repo](https://github.com/Manifaust/k8s-custom-resources-learning-aid/tree/main/examples/my-weather-app-crd).

# Update the Controller Code

Update the controller to implement the Reconcile logic:
1. Get a CheckWeather resource.
2. Check its state from readings its status.
3. If the state is `FINISHED`, then don't do anything, and don't requeue.
4. If the state is `PENDING`, then read the target city from the spec, and check the weather. Write the current temperature to the status. Then set the state to `FINISHED`.
5. If errors occur, then print a log message and requeue.

Start by editing `controllers/checkweather_controller.go` and updating the `Reconcile` function to implement the logic we described.

```go
func (r *CheckWeatherReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
  ctx := context.Background()
  log := r.Log.WithValues("namespace", req.Namespace, "checkweather", req.Name)

  log.Info("Start of Reconcile")

  log.Info("Getting CheckWeather resource")
  var cw weatherappv1alpha1.CheckWeather
  if err := r.Get(ctx, req.NamespacedName, &cw); err != nil {
    log.Error(err, "unable to fetch CheckWeather")
    // we'll ignore not-found errors, since they can't be fixed by an immediate
    // requeue (we'll need to wait for a new notification), and we can get them
    // on deleted requests.

    return ctrl.Result{}, client.IgnoreNotFound(err)
  }

  if cw.Status.State == "" {
    cw.Status.State = weatherappv1alpha1.StatePending
  }

  switch cw.Status.State {
  case weatherappv1alpha1.StatePending:
    city := cw.Spec.City
    log.Info("Checking weather", "city", city)
    temp, err := currentTemp(city)
    if err != nil {
      log.Error(err, "Problem while getting temperature")
      return ctrl.Result{}, err
    }
    log.Info("Obtained temperature", "temp", temp)
    cw.Status.Temperature = int32(temp)
    cw.Status.State = weatherappv1alpha1.StateFinished
  case weatherappv1alpha1.StateFinished:
    log.Info("Work is complete, don't requeue")
    return ctrl.Result{}, nil
  default:

    return ctrl.Result{}, nil
  }

  log.Info("Updating the status of the resource")
  if err := r.Status().Update(ctx, &cw); err != nil {
    log.Error(err, "Error updating status")
    return ctrl.Result{}, err
  }

  // Don't requeue, future changes will trigger Reconcile
  return ctrl.Result{}, nil
}

func currentTemp(city string) (int, error) {
  return 999, nil
}

```

The function `currentTemp` is just a place holder for now to make things simpler. We'll implement the real thing later.

If all that was too hard to read then check out [the example from the learning-aid repo](https://github.com/Manifaust/k8s-custom-resources-learning-aid/tree/main/examples/my-weather-app-crd).

# Update the CheckWeather Sample Resource
Delete the resource we created earlier.

```bash
$ kubectl delete checkweather checkweather-sample
```

Update `config/samples/weather-app_v1alpha1_checkweather.yaml` to include the city property.

```yaml
apiVersion: weather-app.example.com/v1alpha1
kind: CheckWeather
metadata:
  name: checkweather-sample
spec:
  city: "Mexico City"
```

# Rebuild the Controller and Run it
Shut down the running controller, rebuild the controller and run it again:
```bash
$ make install
$ make run
```

# Create the Resource
In another window, apply the updated resource yaml.

```bash
$ kubectl apply -f config/samples/weather-app_v1alpha1_checkweather.yaml
```

If everything is working then in the first window with the running controller you should see the logs we added.

```bash
...
2021-01-17T20:17:09.846-0500	INFO	controllers.CheckWeather	Checking weather{"namespace": "weather-testing", "checkweather": "checkweather-sample", "city": "Mexico City"}
2021-01-17T20:17:09.846-0500	INFO	controllers.CheckWeather	Obtained temperature	{"namespace": "weather-testing", "checkweather": "checkweather-sample", "temp": 999}
2021-01-17T20:17:09.846-0500	INFO	controllers.CheckWeather	Updating the status of the resource	{"namespace": "weather-testing", "checkweather": "checkweather-sample"}
...
```

Also, you should be able to `kubectl describe` the resource to see the weather result in the status field.

```bash
$ kubectl describe checkweather checkweather-sample
Name:         checkweather-sample
Namespace:    weather-testing
...
Spec:
  City:  Mexico City
Status:
  State:        FINISHED
  Temperature:  999
```

# Implement the Weather Request (Optional)
If the fake temperature is not satisfying, then you can implement a real weather request. In order to do so, you will be using the OpenWeatherMap API, a free API for querying weather info. You will need to [sign up on their website](http://home.openweathermap.org/users/sign_up) in order to obtain an API key.

Once you have an API key, you'll be able to add actual weather request logic to `controllers/checkweather_controller.go`. Start by importing this [user-made OpenWeatherMap library](https://github.com/briandowns/openweathermap).

```go
import (
 ...
  owm "github.com/briandowns/openweathermap"
 ...
)
```

Then replace the placeholder code in the `currentTemp` function. 

```go
func currentTemp(city string) (int, error) {
  var apiKey = "<your API key>"

  w, err := owm.NewCurrent("C", "EN", apiKey) // celsius english
  if err != nil {
    return 0, err
  }

  if err := w.CurrentByName(city); err != nil {
    return 0, err
  }

  temp := int(w.Main.Temp)

  return temp, nil
}
```

The code will take the API key and some basic configuration to create a client, then send a request to the OpenWeatherMap server. Warning: in a real production scenario you should not hard-code your API key in the code, it's not secure nor easy to update. Furthermore, you should be careful what kind of network calls or other potentially slow requests create because a controller might have to scale to support handling many resources. With our current implementation, if our controller encounters a hundred CheckWeather resources, then a single slow OpenWeatherMap response could become a bottle neck for the whole service. It might be better to offload work to a pod or job, then check the status of the pod from the controller. The example in the kubebuilder book will show you how to do that.

To test it out, we once again delete our sample resource, perform `make install`, `make run`, and apply our sample resource again. When the controller has processed the request, it'll have set resource's status to have a real temperature.

```bash
$ kubectl describe checkweather checkweather-sample
Name:         checkweather-sample
...
  Temperature:  12
```

You've just implemented a service where users can query the weather using a k8s resource, can you believe it?

# Exercises
* Use `kubebuilder create` to create another controller and CRD, call it WeatherWarning. This resource will read the status of a specific WeatherCheck resource, and if the temperature is colder than a configured value, then it'll log a warning. Using a new CRD/controller to act on information from another resource is a common pattern for providing additional services while maintaining loose coupling between components.
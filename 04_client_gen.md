# Generating a client
Looking inside the `update-codegen.sh` script, we see that it calls another script inside the `vendor` directory named `generate-groups.sh` from the [code-generator](https://github.com/kubernetes/code-generator) project. The script takes in files from our project that describes the custom resource and group we're creating and uses them to generate the clients, listers, informers, and deepcopy functions that correspond with that resource.

Looking at the script in detail:

```sh
$ vendor/k8s.io/code-generator/generate-groups.sh all \
    github.com/manifaust/k8s-custom-resource-lessons/kcrl-client-go/pkg/generated \
    github.com/manifaust/k8s-custom-resource-lessons/kcrl-client-go/pkg/apis \
    samplecontroller:v1alpha1 \
    --output-base "${GOPATH}/src" \
    --go-header-file "hack/boilerplate.go.txt"
```

`gegnerate-groups.sh` take these parameters:
1.  `all` - tells the script to use all the generators.
2.  `.../pkg/generated` - the output directory to put the generated clientset.
3.  `.../pkg/apis` - the input directory to feed in four files that describe our group and resource which `code-generator` needs:
    1. `<group>/register.go`
    2. `<group>/<version>/doc.go`
    3. `<group>/<version>/register.go`
    4. `<group>/<version>/types.go`
4. `samplecontroller:v1alpha1` - the group and version for our resource. Since our input files are copied from the [sample-controller](https://github.com/kubernetes/sample-controller) project, we'll use the the project's group and version.
5. `--output-base "${GOPATH}/src"` - the base directory where we can find the the input directory and generate the output directory.
6. `--go-header-file "hack/boilerplate.go.txt"` - the copyright header to use.

Now that we know what `update-codegen.sh` does under the hood, we can call the script.

```sh
$ cd k8s-custom-resource-lessons
$ cd kcrl-client-go
$ ./hack/update-codegen.sh
Generating deepcopy funcs
Generating clientset for samplecontroller:v1alpha1 at k8s.io/sample-controller/pkg/generated/clientset
Generating listers for samplecontroller:v1alpha1 at k8s.io/sample-controller/pkg/generated/listers
Generating informers for samplecontroller:v1alpha1 at k8s.io/sample-controller/pkg/generated/informers
```

We can look at what's been generated.

```sh
% git status
...
Untracked files:
  (use "git add <file>..." to include in what will be committed)
	pkg/apis/samplecontroller/v1alpha1/zz_generated.deepcopy.go
	pkg/generated/
```

* `zz_generated.deepcopy.go` contains generated deepcopy functions.
* `generated` is a folder that contains the generated clientset, informers, and listers.
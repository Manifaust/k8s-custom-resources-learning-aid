module github.com/Manifaust/k8s-custom-resources-learning-aid/examples/weather-app

go 1.13

require (
	github.com/briandowns/openweathermap v0.16.0
	github.com/go-logr/logr v0.1.0
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	github.com/prometheus/common v0.4.1
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.2
	sigs.k8s.io/controller-runtime v0.5.0
)

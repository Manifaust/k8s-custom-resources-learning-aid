/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	weatherappv1alpha1 "github.com/Manifaust/k8s-custom-resources-learning-aid/examples/weather-app/api/v1alpha1"
	owm "github.com/briandowns/openweathermap"
)

// CheckWeatherReconciler reconciles a CheckWeather object
type CheckWeatherReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=weather-app.example.com,resources=checkweathers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=weather-app.example.com,resources=checkweathers/status,verbs=get;update;patch

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
		log.Info("Unrecognized state, don't requeue")
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

func (r *CheckWeatherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&weatherappv1alpha1.CheckWeather{}).
		Complete(r)
}

package cmd

import (
	"context"
	"fmt"
	opslevelv1 "github.com/opslevel/kubernetes-operator/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

var (
	scheme = runtime.NewScheme()
)

var rootCmd = &cobra.Command{
	Use:   "opslevel-operator",
	Short: "Opslevel Kubernetes Operator",
	Long:  `Opslevel Kubernetes Operator`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Creating Manager...")
		mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{
			Scheme: scheme,
		})
		cobra.CheckErr(err)
		fmt.Println("Creating Controller...")
		err = builder.ControllerManagedBy(mgr).For(&opslevelv1.ClusterIdentifier{}).Owns(&corev1.ConfigMap{}).Complete(&OpsLevelReconciler{Client: mgr.GetClient()})
		cobra.CheckErr(err)
		fmt.Println("Starting...")
		err = mgr.Start(signals.SetupSignalHandler())
		fmt.Println("Shutting down...")
		cobra.CheckErr(err)

		//cl, err := client.New(config.GetConfigOrDie(), client.Options{Scheme: opslevelv1.Scheme})
		//cobra.CheckErr(err)
		//
		//id := &opslevelv1.ClusterIdentifier{}
		//err = cl.Get(context.Background(), getKey("dev", ""), id)
		//cobra.CheckErr(err)
		//fmt.Println(id.Spec.Name)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(opslevelv1.AddToScheme(scheme))

	cobra.OnInitialize(initConfig)
}

func initConfig() {
	viper.SetEnvPrefix("OPSLEVEL")
}

func getKey(name string, namespace string) client.ObjectKey {
	return client.ObjectKey{
		Namespace: namespace,
		Name:      name,
	}
}

type OpsLevelReconciler struct {
	client.Client
}

func (r *OpsLevelReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	resource := &opslevelv1.ClusterIdentifier{}
	err := r.Get(ctx, req.NamespacedName, resource)
	if err != nil {
		return reconcile.Result{}, client.IgnoreNotFound(err)
	}

	finalizerName := "opslevel.com/finalizer"

	if resource.ObjectMeta.DeletionTimestamp.IsZero() {
		// Object is not being deleted
		if !controllerutil.ContainsFinalizer(resource, finalizerName) {
			controllerutil.AddFinalizer(resource, finalizerName)
			if err := r.Update(ctx, resource); err != nil {
				return reconcile.Result{}, err
			}
		}
	} else {
		// Object is being deleted
		if controllerutil.ContainsFinalizer(resource, finalizerName) {
			// Finalizer is present so handle deletion
			cfg := &corev1.ConfigMap{}
			if err := r.Get(ctx, getKey(resource.ObjectMeta.Name, "default"), cfg); err == nil {
				fmt.Println("Deleting Configmap...")
				if err := r.Delete(ctx, cfg); err != nil {
					return reconcile.Result{}, err
				}
			}

			controllerutil.RemoveFinalizer(resource, finalizerName)
			if err := r.Update(ctx, resource); err != nil {
				return reconcile.Result{}, err
			}
		}

		return reconcile.Result{}, nil
	}

	cfg := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.ObjectMeta.Name,
			Namespace: "default",
		},
		Data: map[string]string{
			"Found": resource.Spec.Name,
		},
	}

	if err := controllerutil.SetControllerReference(resource, cfg, r.Scheme()); err != nil {
		return reconcile.Result{}, err
	}

	if err := r.Get(ctx, getKey(resource.ObjectMeta.Name, "default"), cfg); err != nil {
		fmt.Println("Creating Configmap...")
		if err := r.Create(ctx, cfg); err != nil {
			return reconcile.Result{}, err
		}
	} else {
		cfg.Data["Found"] = resource.Spec.Name
		fmt.Println("Updating Configmap...")
		if err := r.Update(ctx, cfg); err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}

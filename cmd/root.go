package cmd

import (
	"context"
	"fmt"
	opslevelv1 "github.com/opslevel/kubernetes-operator/api/v1"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

var rootCmd = &cobra.Command{
	Use:   "opslevel-operator",
	Short: "Opslevel Kubernetes Operator",
	Long:  `Opslevel Kubernetes Operator`,
	Run: func(cmd *cobra.Command, args []string) {
		cl, err := client.New(config.GetConfigOrDie(), client.Options{})
		cobra.CheckErr(err)
		opslevelv1.AddToScheme(cl.Scheme())

		id := &opslevelv1.ClusterIdentifier{}
		err = cl.Get(context.Background(), getKey("dev", ""), id)
		cobra.CheckErr(err)
		fmt.Println(id.Spec.Name)
	},
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
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

package kubernetes

import (
	"bytes"
	"context"
	"fmt"

	"github.com/xonvanetta/kubernetes-git-sync/pkg/encryption/sops"
	"github.com/xonvanetta/kubernetes-git-sync/pkg/git"
	"github.com/xonvanetta/kubernetes-git-sync/pkg/kubernetes/annotations"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"

	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

type controller struct {
	client client.Client
}

func NewManager() (manager.Manager, error) {
	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		return nil, fmt.Errorf("could not create manager: %w", err)
	}

	err = builder.
		ControllerManagedBy(mgr).
		For(&corev1.Secret{}).
		Complete(&controller{
			client: mgr.GetClient(),
		})
	if err != nil {
		return nil, fmt.Errorf("could not create controller: %w", err)
	}
	return mgr, nil
}

//TODO handle deletion

func (c *controller) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	secret := &corev1.Secret{}
	err := c.client.Get(ctx, req.NamespacedName, secret)
	if err != nil {
		return reconcile.Result{}, err
	}

	if !annotations.IsEnabled(secret) {
		return reconcile.Result{}, nil
	}

	yaml, err := printYaml(secret)
	if err != nil {
		return reconcile.Result{}, err
	}

	secret.GetObjectKind()

	if annotations.HasAgeRecipients(secret) {
		ageRecipients := annotations.GetAgeRecipients(secret)
		yaml, err = sops.Encrypt(yaml, ageRecipients)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	authSecret := &corev1.Secret{}
	err = c.client.Get(ctx, annotations.GetGitSecret(secret), authSecret)
	if err != nil {
		return reconcile.Result{}, err
	}

	options, err := git.GetGitCloneOptions(authSecret.Data, secret)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to get clone options: %w", err)
	}

	err = git.Save(ctx, secret, yaml, options)
	return reconcile.Result{}, err
}

func printYaml(object runtime.Object) ([]byte, error) {
	buf := bytes.NewBuffer(nil)

	printer := printers.YAMLPrinter{}
	err := printer.PrintObj(object, buf)
	return buf.Bytes(), err
}

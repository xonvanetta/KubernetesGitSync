package annotations

import (
	"strings"

	"github.com/go-git/go-git/v5"
	"k8s.io/apimachinery/pkg/types"
)

const (
	Enabled           = "kubernetes-git-sync/enabled"
	SopsAgeRecipients = "kubernetes-git-sync/sops-age-recipients"
	GitSecret         = "kubernetes-git-sync/git-secret"
	GitFilepath       = "kubernetes-git-sync/git-filepath" //TODO add multiple paths
	GitBranch         = "kubernetes-git-sync/git-branch"   //TODO No work
	GitUrl            = "kubernetes-git-sync/git-url"
	GitCommitOptions  = "kubernetes-git-sync/git-commit-options"
)

type ObjectMeta interface {
	GetLabels() map[string]string
	GetAnnotations() map[string]string
}

func IsEnabled(object ObjectMeta) bool {
	return object.GetAnnotations()[Enabled] == "true"
}

func HasAgeRecipients(object ObjectMeta) bool {
	return object.GetAnnotations()[SopsAgeRecipients] != ""
}

func GetAgeRecipients(object ObjectMeta) string {
	return object.GetAnnotations()[SopsAgeRecipients]
}

func GetGitFilepath(object ObjectMeta) string {
	return object.GetAnnotations()[GitFilepath]
}

func GetGitBranch(object ObjectMeta) string {
	return object.GetAnnotations()[GitBranch]
}

func GetGitSecret(object ObjectMeta) types.NamespacedName {
	parts := strings.SplitN(object.GetAnnotations()[GitSecret], string(types.Separator), 2)
	if len(parts) != 2 {
		return types.NamespacedName{}
	}

	return types.NamespacedName{
		Namespace: parts[0],
		Name:      parts[1],
	}
}

func GetGitCloneOptions(object ObjectMeta) *git.CloneOptions {
	a := object.GetAnnotations()
	return &git.CloneOptions{
		URL: a[GitUrl],
	}
}

func GetGitCommitOptions(obj ObjectMeta) *git.CommitOptions {
	opts := &git.CommitOptions{}
	//a := object.GetAnnotations()

	return opts
}

package git

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"github.com/xonvanetta/kubernetes-git-sync/pkg/kubernetes/annotations"
	ssh2 "golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const DefaultUser = "git"

type Object interface {
	GetName() string
	GetNamespace() string
	GetObjectKind() schema.ObjectKind

	annotations.ObjectMeta
}

func openOrCloneRepository(ctx context.Context, path string, options *git.CloneOptions) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil && errors.Is(err, git.ErrRepositoryNotExists) {
		repo, err = git.PlainCloneContext(ctx, path, false, options)
		if err != nil {
			return nil, fmt.Errorf("failed to clone repository: %w", err)
		}
	}
	if err != nil {
		return nil, fmt.Errorf("failed to open repository on path: %w", err)
	}

	return repo, nil
}

func GetGitCloneOptions(data map[string][]byte, object Object) (*git.CloneOptions, error) {
	cloneOptions := annotations.GetGitCloneOptions(object)

	user := string(data["user"])
	if user == "" {
		user = DefaultUser
	}

	//Todo add validation for secret containing fields
	auth, err := ssh.NewPublicKeys(user, data["private-key"], string(data["password"]))
	if err != nil {
		return cloneOptions, fmt.Errorf("failed to create auth with public keys: %w", err)
	}

	//Todo implement this from the secret.
	auth.HostKeyCallback = ssh2.InsecureIgnoreHostKey()

	cloneOptions.Auth = auth

	return cloneOptions, nil
}

func Save(ctx context.Context, object Object, yaml []byte, cloneOptions *git.CloneOptions) error {
	tempPath := filepath.Join(os.TempDir(), object.GetNamespace(), object.GetName())
	repo, err := git.PlainCloneContext(ctx, tempPath, false, cloneOptions)
	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	//branchName := annotations.GetGitBranch(object)
	//repo.Branch(branchName)

	defer func() {
		//TODO reuse fs clone with openOrCloneRepository
		os.RemoveAll(tempPath)
	}()

	w, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to create worktree: %w", err)
	}

	//if branch := annotations.GetGitBranch(object); branch != "" {
	//	err = w.Checkout(&git.CheckoutOptions{
	//		Branch: plumbing.NewBranchReferenceName(branch),
	//		Create: true,
	//	})
	//	if err != nil {
	//		return fmt.Errorf("failed to checkout branch: %w", err)
	//	}
	//}

	filepath := annotations.GetGitFilepath(object)
	if filepath == "" {
		return fmt.Errorf("cannot save file to empty filepath")
	}

	file, err := w.Filesystem.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file in git: %w", err)
	}
	defer file.Close()

	_, err = file.Write(yaml)
	if err != nil {
		return fmt.Errorf("failed to write updated secret to path: %w", err)
	}

	_, err = w.Add(filepath)
	if err != nil {
		return fmt.Errorf("failed to add file to git: %w", err)
	}

	status, err := w.Status()
	if err != nil {
		return fmt.Errorf("failed to get git status: %w", err)
	}

	if status.IsClean() {
		return nil
	}

	_, err = w.Commit(generateCommitMessage(object), annotations.GetGitCommitOptions(object))
	if err != nil {
		return fmt.Errorf("failed to commit: %w", err)
	}

	pushOptions := &git.PushOptions{
		Auth: cloneOptions.Auth,
	}
	err = repo.PushContext(ctx, pushOptions)
	if err != nil {
		return fmt.Errorf("failed to git push: %w", err)
	}
	return nil
}

func generateCommitMessage(object Object) string {
	commitMessage := fmt.Sprintf(`chore: update kubernetes %s %s/%s`,
		object.GetObjectKind().GroupVersionKind().Kind, object.GetName(), object.GetNamespace())
	return commitMessage
}

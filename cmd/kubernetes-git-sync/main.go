package main

import (
	"log"

	"github.com/xonvanetta/kubernetes-git-sync/pkg/kubernetes"
	"github.com/xonvanetta/shutdown/pkg/shutdown"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func main() {
	ctx := shutdown.Context()

	logf.SetLogger(zap.New())

	//err := cmd.Execute()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//return

	//repo, err := git.PlainCloneContext(ctx, config.RepositoryPath, false, &config.GitCloneOptions)
	//if err != nil {
	//	log.Fatalf("failed to setup repo: %s", err)
	//}

	mgr, err := kubernetes.NewManager()
	if err != nil {
		log.Fatalf("failed to setup manager: %s", err)
	}

	if err := mgr.Start(ctx); err != nil {
		log.Fatalf("could not start manager: %s", err)
	}
}

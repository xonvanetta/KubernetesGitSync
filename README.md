# KubernetesGitSync
Securely sync your Kubernetes objects to Git.

### Idea
Monitor Kubernetes objects for changes based on annotations or labels. When changes are detected, write the updated YAML to a specified Git repository, then commit and push the changes. If the object is sensitive (e.g., annotated or a secret), encrypt it using SOPS before committing.

### Use Case
A cert-manager in Cluster A generates a wildcard TLS secret. This TLS secret should be propagated to all other clusters without generating new ones. GitOps can be used to sync this secret across clusters.

### Cons
Managing ownership with GitOps can be challenging, particularly regarding who owns which resources. For instance, the cluster that initially creates a secret should not have that secret managed by GitOps, as this could lead to conflicts over updates. It's important to test scenarios where this issue might arise to determine its impact in real-world situations.

### Note
This project is a work in progress, as I have found another way to solve my issue.
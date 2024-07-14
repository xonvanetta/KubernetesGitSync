# KubernetesGitSync
Securely syncing your Kubernetes secrets to Git.


### Idea
Check if changes in annotated/labeled object, if changes detected write yaml to given git repositry, commit and push it.
If object is sensitive ie annotated or secret then use encryption(SOPS) on it before commit. 

### Use case
Having a cert-manager in cluster A that would generate a wilcard tls secret, now this tls secret should be added to all other clusters and not generate a new one, gitops would sync it.
 

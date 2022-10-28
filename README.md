# Argocd Autopilot Api

## Description
An api wrapper for argocd-autopilot that will be used by backstage.io

## Example
``` bash
curl http://localhost:8080/run \                                                  
    --include \
    --header "Content-Type: application/json" \
    --request "POST" \
    --data '{"git-repo": "https://github.com/tony-mw/autotest-argo.git","git-token-path": "/Users/$(whoami)/.github_token","root-command": "argocd-autopilot","args": ["repo","bootstrap"]}'
```
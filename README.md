# keink

keink(represent for [KubeEdge](https://github.com/kubeedge/kubeedge) IN [kind](https://github.com/kubernetes-sigs/kind)) is a tool for running local KubeEdge clusters using Docker container "nodes".

**Notice**: 
- The keink project is currently in alpha state, functionalities may subject to change. Please feel free to share your thoughts and contribute.
- This project is developed based on [kind](https://github.com/kubernetes-sigs/kind).

## Integrate with [KubeEdge](https://github.com/kubeedge/kubeedge)

### Prerequisites
We need to put KubeEdge codes to `$GOPATH/github.com/kubeedge/kubeedge`.
And we need to ensure them checkout to the right release branch. For example, we can run the command like `git checkout v1.11.1` to checkout the right release branch.

When running `bin/keink build edge-image` command to build `kubeedge/node` image, which contains KubeEdge components `cloudcore`, `edgecore` and `keadm` based on the [`kindest/node`](https://hub.docker.com/r/kindest/node) image, keink will use the above KubeEdge source codes.

### Build KubeEdge customized node image and start KubeEdge cluster

Build keink from source code, build kubeedge/node image, and create cluster.
```shell
# build keink binary
make
# build kubeedge/node image
bin/keink build edge-image
# create KubeEdge cluster based on the k8s cluster
bin/keink create kubeedge --image kubeedge/node:latest --wait 120s
```

And you can also download `keink` official release from [github release page](https://github.com/kubeedge/keink/releases), and then you can just run `keink create kubeedge` command directly to start KubeEdge cluster with `kubeedge/node:${version_tag}` image that KubeEdge officially published.

And then you can see the contents as follows:
```shell
# bin/keink create kubeedge --image kubeedge/node:latest --wait 120s 
Creating cluster "kind" ...
 ✓ Ensuring node image (kubeedge/node:latest) 🖼
 ✓ Preparing nodes 📦 📦  
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
 ✓ Joining worker nodes 🚜 
 ✓ Waiting ≤ 2m0s for control-plane = Ready ⏳ 
 • Ready after 0s 💚
 ✓ Starting KubeEdge 📜
```


and you can see all the nodes are in `Ready` status, included `control-plane`, `worker` and `edge-node`(KubeEdge edgecore nodes).
```shell
# kubectl get node -owide
NAME                 STATUS   ROLES                  AGE    VERSION                                                   INTERNAL-IP   EXTERNAL-IP   OS-IMAGE       KERNEL-VERSION       CONTAINER-RUNTIME
kind-control-plane   Ready    control-plane,master   116s   v1.23.4                                                   172.18.0.2    <none>        Ubuntu 21.10   4.15.0-169-generic   containerd://1.5.10
kind-worker          Ready    agent,edge             50s    v1.23.15-kubeedge-v1.13.0-beta.0.7+fb22a08cd41a52-dirty   172.18.0.3    <none>        Ubuntu 21.10   4.15.0-169-generic   containerd://1.5.10
```

Deploy a nginx demo to edge-node
```shell
kubectl apply -f ./pod.yaml
```

```
# kubectl get pod -owide
NAME    READY   STATUS    RESTARTS   AGE   IP           NODE          NOMINATED NODE   READINESS GATES
nginx   1/1     Running   0          30s   10.244.1.2   kind-worker   <none>           <none>
```

And nginx pod will be successfully running on the edge node. Congratulations, KubeEdge cluster is running successfully using `keink`!


## Contributing

If you're interested in being a contributor and want to get involved in
developing the KubeEdge code, please see [CONTRIBUTING](./CONTRIBUTING.md) for
details on submitting patches and the contribution workflow.


## License

keink is under the Apache 2.0 license. See the [LICENSE](license) file for details.
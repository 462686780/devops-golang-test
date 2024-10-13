StatefulSet核心功能是：包括Pod管理、持久存储（Persistent Volumes）的分配、服务发现等。我可以尝试模拟一些核心功能，比如Pod的顺序创建和持久存储的分配。

ValidatingAdmissionWebhook这个webhook 将在 Kubernetes API 服务器上运行，并在资源（如 Pod）创建或更新之前对其进行验证。如果验证失败，webhook 将拒绝请求。

helm chart我不太熟悉,但他就是一个模板跟ansible playbook一样，我这里手写的，就没这个东西

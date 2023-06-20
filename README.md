
## zero-example-mqueue

---

下面有一个项目块`scheduler`；

`scheduler`是定时任务服务；依赖`asynq.NewTask`向redis发送一条job任务，由每个服务的`job`执行，要求所有服务都使用一个redis服务。
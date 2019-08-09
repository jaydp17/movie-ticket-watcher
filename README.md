# movie-ticket-watcher
Looks for a particular movie in a particular cinema &amp; notifies you as soon as it goes live

# Local development

```sh
# this will watch & rebuild all the lambda function
make watch

# it'll start an api gateway simulator on port 3000
make local-api

# it starts a cors proxy, because stupid SAM CLI does't support CORS in local 
make local-cors-proxy
```

## TODO

Refer to [this Pivotal Tracker board](https://www.pivotaltracker.com/n/projects/2375220) for tasks & their status 
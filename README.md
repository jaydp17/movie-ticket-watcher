<p align="center">
  <a href="https://movie-ticket-watcher.jaydp.com">
    <img alt="Movie Ticket Watcher" src="/img/logo.png" width="60" />
  </a>
</p>
<h1 align="center">
  Movie Ticket Watcher (backend)
</h1>

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

## 🤝 Contributing

Contributions, issues and feature requests are welcome.<br />
Feel free to check [issues page](https://github.com/jaydp17/movie-ticket-watcher/issues) if you want to contribute.

## :memo: License

Licensed under the [MIT License](./LICENSE).
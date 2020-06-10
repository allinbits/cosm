# Cosmos CLI

Code scaffolding tool for Cosmos SDK applications.

```
cosmos-cli app github.com/example/blog
```

Generates a Cosmos SDK application.

To launch the app run `go mod tidy && make && ./init.sh && blogd start`.

```
cosmos-cli type post title body
```

Generates a type `Post` with two fields: `title` and `body`.

To add a post run `blogcli tx blog create-post "My title" "This is a blog" --from=me`.

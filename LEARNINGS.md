# GO-ing for web development


## The Go programming language

[Go](https://tour.golang.org/welcome/1) is an open source programming language designed for building simple, fast, and reliable software.

Go is a battery included programming language it has a web server already built in in the `net/http` package containing all of the 
functionality of the HTTP protocol.


## My Learning references

- [Go Language Tour](https://tour.golang.org/welcome/1)
- [Go Web Examples](https://gowebexamples.com/hello-world/)
- [Project structure](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1)
- [Writting REST APIs in Go](https://www.ribice.ba/rest-api-go3/)
- [Routing in net/http](https://subscription.packtpub.com/book/application_development/9781786468666/1/ch01lvl1sec10/routing-in-net-http)
- [A Recap of request handling](https://www.alexedwards.net/blog/a-recap-of-request-handling)


## Project layout decision 

Went for a [domain approach](https://medium.com/@benbjohnson/standard-package-layout-7cdbc8391fc1) and in the future the intention is to scale into a mix of of [Gorsk](https://github.com/ribice/gorsk) (service like approach).

TLDR; 

Root package is for domain types, simple types like `User`, `Product`, `Post`.

> The root package should not depend on any other package in your application!

Since dependencies are not allowed in the root, then they are pushed to subpackages, these work as adapters between the domain and the
implementation.

Ex: 
  The `UserService` might be backed by a PostgreSQL database, then we would have a `postgres` subpackage that would represent a PostgreSQL
  implementation of the user.

```go
  type UserService struct {
    DB *sql.DB
  }

  func (s *UserService) User(id int) (*myapp.User, error) {
    var u myapp.User
    row := db.QueryRow(`SELECT id, name FROM users WHERE id = $1`, id)
    if row.Scan(&u.ID, &u.Name); err != nil {
      return nil, err
    }
    return &u, nil
  }
```

## The `net/http` package

[`net/http`](https://golang.org/pkg/net/http/)

#### The **handler interface**:

Handlers are responsible for writing response headers and bodies. Almost any object can be a handler, so long as it satisfies the http.Handler interface

```go
// Response/Request data
type Handler interface {
  ServeHTTP(http.ResponseWriter, *http.Request)
}
```

#### Creating a request handler:

Here, the `/hello` function is a valid handler because it implements `(w http.ResponseWriter, r *http.Request)`

```go
http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, you've requested: %s\n", r.URL.Path)
})
```

### Creating the `httpserver` package

Implements an http server from the `net/http` package and exposes methods to Run and Create instances.

#### Routing:

Decided to use the `http.DefaultServeMux` instead of third party libraries like `gorilla/mux` to learn more about the `net/http`.

The router will be responsible for matching the URL of each incoming request against a list of registered patterns and calling the corresponding handler.

#### Serving static assets:

After creating a basic file server and exposing it under /static/ users can retrieve assets like:

http://localhost:7777/static/bidoof.jpg


#### Exposing an endpoint:

In my early attempt to implement a service based structure I created the `auth` service that exposed a simple `/auth/register` endpoint
for testing purposes.

```go
package api

import (
	"goproject/auth"
	authTransport "goproject/auth/transport"
	"goproject/httpserver"
)

// Run the http server
func Run() {
	s := httpserver.New()

  // Auth endpoints receives a service and a router
	authTransport.NewHTTP(auth.Initialize(), s.Router)

	s.Run()
}
```


### Connecting to a database:
*TODO*


#### SQL Injection: 

- Altough using an ORM would most likely solve this, I need to be carefull when building SQL statements:

```go
func buildSql(email string) string {
  return fmt.Sprintf("SELECT * FROM users WHERE email='%s';", email)
}
```
For instance if email had the value of `'; DROP TABLE users;'` then the query would be:

```sql
SELECT * FROM users WHERE email=''; DROP TABLE users;'';
```

Yikes... user table is gone!

So `database/sql` is awere of this and knows what is valid SQL and what is nefarious, so if we would pass a '' it would escape them.

TLDR; 

> So the short version of this story is *always use the database/sql package to construct SQL statements and insert values into them: 
[source](https://www.calhoun.io/inserting-records-into-a-postgresql-database-with-gos-database-sql-package/)


### Adding Logging:
*TODO*

### Adding Auth:
*TODO*

### Adding Error handling:
*TODO*
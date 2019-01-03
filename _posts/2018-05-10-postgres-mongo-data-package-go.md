---
layout: post
title: "Handling PostgreSQL and MongoDB in one Go data package — excerpt from my book."
date: 2018-05-10 11:26:37 UTC
updated: 2018-05-10 11:26:37 UTC
comments: false
summary: "..."
---

This is some excerpts from my book “[Build a SaaS app in
Go](https://buildsaasappingo.com/)” currently available for pre-order at a
discounted price.

When I announced the availability of the preview chapter in December 2017 on
Reddit and Twitter, I got quite surprising reactions regarding my choice of
MongoDB for the database engine.

I was surprised at how MongoDB is treated in the Go community. Indeed MongoDB is
not the choice for all scenario. My main thoughts were that it is quicker to
change your models while building your MVP than an RDBMS. And the amount of code
needed to get things done is usually lower, less boilerplate.

I do agree though that most often than not a NoSQL database is not the right
option. [Roadmap](https://roadmap.space/), for example, is similar to Trello in
a sense where all information of the system is related to a “card.” Trello is
using MongoDB; I picked MongoDB for Roadmap based on the fact that I knew I was
not going to need more relations and complex queries.

Truth is for some queries it would be simpler to use a relational database. At
LeadFuze I started building with a NoSQL database. It worked nicely until the
time came to produce reports and more complex statistical queries and what not.
We migrated to a SQL database.

I have never been a massive fan of ORM (object-relational mapper). In this
chapter, we’re going to write a database engine agnostic data package that will
handle MongoDB or PostgreSQL out of the box.

Since what we’re building in this book will eventually become an open source
package this seems to be the right route to take. Plus it might be a good
exercise to compare how it feels to use one paradigm over another.

The following example is an ideal implementation of a function using our data
package:

```go
import "github.com/dstpierre/gosaas/data"
//...
func (u User) detail(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	db := ctx.Value(engine.ContextDB)
	id := ctx.Value(engine.ContextUserID)

	detail, err := db.Users.GetDetail(id)
	if err != nil 
		engine.Respond(w, r, http.StatusInternalServerError)
		return
	}
	engine.Respond(w, r, http.StatusOK, detail)
}
```
#### Simple model struct

The first challenge we need to fix is how we can pass the database connection
for two completely different types.

*Also, bare with me even if you are not caring about MongoDB or vice-versa, what
we will build in this chapter can be applied to other problems.*

We need to use [type alias](https://golang.org/doc/go1.9) to be able to abstract
the fact that the SQL type and Mongo type are not using the same interface.

```go
type DB struct {
	DatabaseName 
Connection   *model.Connection
}
```

Our `DB` struct has a database name field and a custom type alias field that is
defined in our model package.

We use [build tags](https://golang.org/pkg/go/build/) to link our type alias
either with the `database/sql`’s `sql.DB`type or MongoDB’s connection. At build
time the Go compiler will use either the `data/model/types_sql.go`file or the
`data/model/types_mgo.go`file to assign the right underlying connection type to
our type alias.

Here’s the two files:

**File: data/model/types_sql.go**

```go
// +build !mgo

package model

import "database/sql"

type Connection = sql.DB
type Key = int64
```

**File: data/model/types_mgo.go**

```go
// +build mgo

package model

import (
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Connection = mgo.Session
type Key = bson.ObjectId
```

The SQL file uses the `!mgo`tag. By default, the compiler will build the SQL
files “_sql_”.

Depending on how we’re going to build our program the type alias will pass
around the right database connection.

Another issue we have similar to the connection types being different is the
primary key type being an `int64`inside SQL (usually) and `bson.ObjectId`for
MongoDB. We’re using another type alias for this. We will use a type alias named
`Key`.

The goal is to create a shareable model package that both the SQL and MongoDB
implementation can reference. Since we created a `Key`type alias that received
either an `int64`or a `bson.ObjectId`type at compile time we can safely design
our models for both engine.

```go
package model

type User struct {
	ID    Key    `bson:"_id" json:"id"`
	Email 
	`bson:"email" json:"email"`
}
```

We use our `Key`type for the primary key type. We also define the metadata
information for the field name inside MongoDB and how this field look like when
we encode/decode to/from a JSON payload.

This is a simple example, but you might start to see how we will be able to pass
our model around with a nicely wrapped `Key`type that adjusts based on our build
tags.

#### A service interface

The main data package exposes `DB` that contains the `model.Connection` properly
aliasing the desired database connection and the services which are interfaces
and implemented into their own packages `data/pg` for Postgres and `data/mongo`
for Mongo.

**File: data/db.go**

```go
type UserServices interface {
	GetDetail(id model.Key) (*model.User, 
)
}

type DB struct {
	DatabaseName 
Connection   *model.Connection
	CopySession  

	Users UserServices
}
```

If you recall at the beginning of this chapter, we saw this line of code which
was how I wanted the data package to be called independently of if we choose SQL
or MongoDB.

```go
detail, err := db.Users.GetDetail(id)
```

The `Users`service field in the `DB`type is an interface that both the `pg`and
`mongo`package will implement.

Here’s the `mongo`implementation:

**File: data/mongo/users.go**

```go
package mongo

import (
	"github.com/dstpierre/gosaas/data/model"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type Users struct {
	DB *mgo.Database
}

func (u *Users) GetDetail(id model.Key) (*model.User, 
) {
	var user model.User
	where := bson.M{"_id": id}
	if err := u.DB.C("users").Find(where).One(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
```

Similarly to the SQL implementation we receive an `bson.ObjectId` via the
`idmodel.Key`as the primary key type that we perform a query on our users
collection.

> Note that in the book all unit tests are defined and lot more detailed, this
> blog post take some excerpts of the chapter and resume the usage of build tags
and interfaces to make this possible.

Inside the book, there’s a section on how to start a PostgreSQL and MongoDB
server via docker. All unit tests, an in-memory implementation for unit tests.
Integration tests and full SQL and Mongo implementation.

The source code is available with all pre-order, and the price of the book at
this time of writing is $17 and increases each time a new chapter is added (+$
two each chapter).

* [Golang](https://dominicstpierre.com/tagged/golang?source=post)
* [Database](https://dominicstpierre.com/tagged/database?source=post)
* [Postgresql](https://dominicstpierre.com/tagged/postgresql?source=post)
* [Mongodb](https://dominicstpierre.com/tagged/mongodb?source=post)


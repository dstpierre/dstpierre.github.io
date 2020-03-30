---
permalink: "/convert-dbml-to-go-data-access"
layout: post
title: "Convert .dbml Linq-to-SQL to Go, and parsing XML in Go"
date: 2020-03-30 07:55:40 UTC
updated: 2020-03-38 07:55:40 UTC
comments: false
summary: "..."
---

### Repo on GitHub

<div class="github-card" data-github="dstpierre/dbmltogo" data-width="400" data-height="" data-theme="default"></div>
<script src="//cdn.jsdelivr.net/github-cards/latest/widget.js"></script>

One of the many aspects I dig about Go is the backward compatibility they have. 
A 6-7-8 years old Go program is almost guaranty to compile with the latest Go 
compiler.

It's something that Microsoft and the .NET team should be more focused on IMHO. 

I've been writing C# code since 2001, and I still maintain legacy C# applications 
at a > 25-year-old fintech company as consultant.

There's 2-3 web applications that were written with Linq-to-SQL (.dbml file). The 
support of that feature is long gone in the .NET framework, and .NET Core does not 
have an option to handle those file as well.

Instead of migrating to .NET Core and rewrite those parts in something like 
`System.Data.SqlClient` or Entity Framework I decided to try prototyping what a 
fresh project could look like using Go.

### dbmltogo

I wrote a small [tool that parses a DBML file](https://github.com/dstpierre/dbmltogo) 
and generates Go boilerplate data access code targeting Microsoft SQL Server.

The idea is to have a complete starting point for all tables in a SQL Server 
database ready for CRUD operations.

This is how a specific `Bank` table could be used:

```go
banks := data.Banks{DB: db}
bank, err := banks.GetByID(123)
```

The plural name is where all the functions are attached:

```go
// Banks ...
type Banks struct {
	DB *sql.DB
}

// Create creates a new Bank
func (repo *Banks) Create(entity Bank) (id int, err error) {}
```

And the singular is the struct containing all the fields with proper types and 
JSON name:

```go
// Bank ...
type Bank struct {
	BankID         int            `json:"bankId"`
	BankNameFr     JSONNullString `json:"bankNameFr"`
	BankNameEn     JSONNullString `json:"bankNameEn"`
	CentralAddress JSONNullString `json:"centralAddress"`
	CentralPhone   JSONNullString `json:"centralPhone"`
	CentralFax     JSONNullString `json:"centralFax"`
	CommentA       JSONNullString `json:"commentA"`
	CommentB       JSONNullString `json:"commentB"`
	Short          JSONNullString `json:"short"`
	Fee            float64        `json:"fee"`
	DisplayOrder   JSONNullInt32  `json:"displayOrder"`
}
```

You pass a `*sql.DB` to the repository `struct` and you'll have a `Create`, 
`List`, `GetByID`, `Update` and `Delete` functions *ready* to be used.

### Parsing XML, generating Go code

Let's look at the XML representation of a tables in a typical DBML file:

```xml
<Table Name="dbo.Demos" Member="Demos">
    <Type Name="Demo">
      <Column Name="DemoID" Type="System.Int32" DbType="Int NOT NULL IDENTITY" IsPrimaryKey="true" IsDbGenerated="true" CanBeNull="false" />
			<Column Name="ParentID" Type="System.Int32" DbType="Int NOT NULL" CanBeNull="false" />
      <Column Name="Name" Type="System.String" DbType="NVarChar(50)" CanBeNull="true" />
      <Association Name="Parent_Demo" Member="Parent" ThisKey="ParentID" OtherKey="ParentID" Type="Parent" IsForeignKey="true" />
    </Type>
  </Table>
```

This is the structs we need to model this XML:

```go
type Table struct {
	XMLName xml.Name `xml:"Table"`
	Name    string   `xml:"Name,attr"`
	Member  string   `xml:"Member,attr"`
	Type    Type     `xml:"Type"`
}

type Type struct {
	XMLName      xml.Name      `xml:"Type"`
	Name         string        `xml:"Name,attr"`
	Columns      []Column      `xml:"Column"`
	Associations []Association `xml:"Association"`
}

type Column struct {
	XMLName       xml.Name `xml:"Column"`
	Name          string   `xml:"Name,attr"`
	CSharpType    string   `xml:"Type,attr"`
	DBType        string   `xml:"DbType,attr"`
	IsPrimaryKey  string   `xml:"IsPrimaryKey,attr"`
	IsDbGenerated string   `xml:"IsDbGenerated,attr"`
	CanBeNull     string   `xml:"CanBeNull,attr"`
}

type Association struct {
	XMLName  xml.Name `xml:"Association"`
	Name     string   `xml:"Name,attr"`
	Member   string   `xml:"Member,attr"`
	ThisKey  string   `xml:"ThisKey,attr"`
	OtherKey string   `xml:"OtherKey,attr"`
	Type     string   `xml:"Type,attr"`
}
```

We can parse the DBML file as follow:

```go
// error handling omitted
b, _ := ioutil.ReadFile(xmlfile)

var db Database
xml.Unmarshal(b, &db)
```

All the `encoding` packages in Go are just beautiful. Most packages are just 
amazing in fact ;).

Code generation is not complex when using the right tool.

At first I used `fmt.Sprintf` but I quickly realized that using `text/template` 
would be much cleaner.

This is an excerpt of the Go entity template:

```
package {{.PkgName}}

import (
	"database/sql"
	"time"
)

// {{.EntityName}} ...
type {{.EntityName}} struct {
{{range .Fields}}	{{.Name}} {{.Type}} `json:"{{.JSONName}}"`
{{end}}
}

// {{.MemberName}} ...
type {{.MemberName}} struct {
	DB *sql.DB
}

// Create creates a new {{.EntityName}}
func (repo *{{.MemberName}}) Create(entity {{.EntityName}}) (id {{.PKType}}, err error) {
	err = repo.DB.QueryRow(`
		INSERT INTO {{.TableName}}
		VALUES({{range .InsertFields}}@{{.Name}},
		{{end}}
		)
		SELECT SCOPE_IDENTITY()
	`, {{range .InsertFields}}sql.Named("{{.Name}}", entity.{{.Name}}),
		{{end}}
	).Scan(&id)
	return
}
```

It's using the Go `text/template` package. Here's an example:

### Converting the C# type to Go

The first challenge is to match the C# types to Go. I'm using a simple switch 
case for that, and I also created custom JSONNull types to have simpler marshaling 
to JSON instead of the SQL package default object.

```go
func gotype(c Column) (typ string) {
	if c.CanBeNull == "true" {
		switch c.CSharpType {
		case "System.Byte":
			typ = "JSONNullInt32"
		case "System.Int16":
			typ = "JSONNullInt32"
		case "System.Int64":
			typ = "JSONNullInt64"
		case "System.Guid",
			"System.String":
			typ = "JSONNullString"
		case "System.DateTime":
			typ = "JSONNullTime"
		case "System.Single",
			"System.Decimal",
			"System.Double":
			typ = "JSONNullFloat64"
		case "System.Data.Linq.Binary":
			typ = "[]byte"
		case "System.Int32":
			typ = "JSONNullInt32"
		case "System.Boolean":
			typ = "JSONNullBool"
		default:
			fmt.Println("unhandled type: ", c.CSharpType)
		}

		return
	}

	switch c.CSharpType {
	case "System.Byte":
		typ = "byte"
	case "System.Int16":
		typ = "int"
	case "System.Int64":
		typ = "int64"
	case "System.Guid":
		typ = "string"
	case "System.String":
		typ = "string"
	case "System.DateTime":
		typ = "time.Time"
	case "System.Decimal",
		"System.Single",
		"System.Double":
		typ = "float64"
	case "System.Data.Linq.Binary":
		typ = "[]byte"
	case "System.Int32":
		typ = "int"
	case "System.Boolean":
		typ = "bool"
	default:
		fmt.Println("unhandled type: ", c.CSharpType)
	}

	return
}
```

All the remaining needed information is already present in the XML file, we're 
ready to generate some Go code.

### It can also generate Elm modules

An experimental option is to generate Elm modules and api endpoint moduless.

My next step is to test from a real user screen if those are comfortable to use.

If they are I'll probably generate routes and handlers in Go for all types. More 
will be added later if I find this option to be the one I pursue.
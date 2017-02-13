# Gofield [![Build Status](https://travis-ci.org/enrichman/gofield.svg?branch=master)](https://travis-ci.org/enrichman/gofield)

Gofield will reduce your JSON serialized Go struct in a lighter object to keep your client happy!

You can use the "Facebook convention" in order to select the fields object: comma separated values, with curly brackets for the inner objects.

```go
fields := "foo,bar{inner_foo,inner_bar}"
lighterObj := gofield.Reduce(fatObj, fields)
```

i.e.:

if you have a JSON like

```json
{
    "name": "Enrico",
    "surname": "Candino",
    "age": 27,
    "city": "Rome",
    "email":[
        {
            "name": "work",
            "email": "mymail@work.it"
        },
        {
            "name": "personal",
            "email": "mymail@personal.it"
        }
    ]
}
```

you can reduce it like this selecting the fields `name,email{name},city`:

```go
gofield.Reduce(fatObj, "name,email{name},city")
```

```json
{
    "name": "Enrico",
    "city": "Rome",
    "email":[
        {
            "name": "work"
        },
        {
            "name": "personal"
        }
    ]
}
```

You can obviously nest as many levels as you want!

A sample server can be run on the `:8080` to test it. 

```shell
go run cmd/gofield-server/main.go
```

Just send your JSON with a POST on the `/reduce` endpoint, with a query parameter `fields`. :)

```
POST /reduce?fields=foo,bar{inner_bar}
```

Note:  
performance were not benchmarked, so do your homework before use it on production!

Issues and PR are welcome!

Developed By
--------

Enrico Candino - www.enricocandino.it

<a href="https://twitter.com/enrichmann">
  <img alt="Follow me on Twitter"
       src="http://icons.iconarchive.com/icons/danleech/simple/96/twitter-icon.png" />
</a>
<a href="https://plus.google.com/+EnricoCandino">
  <img alt="Follow me on Google+"
       src="http://icons.iconarchive.com/icons/danleech/simple/96/google-plus-icon.png" />
</a>
<a href="https://it.linkedin.com/in/enrico-candino-78995553">
  <img alt="Follow me on LinkedIn"
       src="http://icons.iconarchive.com/icons/danleech/simple/96/linkedin-icon.png" />
</a>


License
--------

    The MIT License (MIT)
    
    Copyright (c) 2017 Enrico Candino
    
    Permission is hereby granted, free of charge, to any person obtaining a copy
    of this software and associated documentation files (the "Software"), to deal
    in the Software without restriction, including without limitation the rights
    to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
    copies of the Software, and to permit persons to whom the Software is
    furnished to do so, subject to the following conditions:
    
    The above copyright notice and this permission notice shall be included in all
    copies or substantial portions of the Software.
    
    THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
    IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
    FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
    AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
    LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
    OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
    SOFTWARE.
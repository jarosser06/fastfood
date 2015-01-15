fastfood
========

Makes Cookbooks faster by pre-templating them

This tool is still under heavy development and tbh for the vast majority
of people out there it will be uninteresting and unuseful.

###Hacking
This tool uses gpm for dependency management so you need to
ensure that it is installed in order to build from source.

https://github.com/pote/gpm

###Known um ... Issues
Passing node attributes does work ... however you need to escape the
' chars in order to insure it populates the recipe correctly.

###Notes

Template Notes
```
If the variable in the template can be either a node attribute
or a string you can ensure the value will be formatted correctly
using the .QString method on the variable.

Example:
template |{ .QString .Name }| do

if .Name is a node attribute such as node['mycookbook']['myname']
then it will be written as such:

template node['mycookbook']['myname'] do

if .Name is a string such as templatename then it will be written
with single quotes around it:

template 'templatename' do
```

Create an empty cookbookk:
```shell
fastfood new mycookbook
```

Generate a set of recipes, etc where app is the provider:
```shell
fastfood gen app name:application1 type:python repo:git@github.com/cus/app1
```

Build/modify an existing cookbook from a config:
```shell
fastfood build ./mycookbook.json
```

Uses a json config file to allow for pre-templating cookbooks:
```json
{
  "name": "123456-customer",
  "providers": [
    {
      "provider": "django_app",
      "name": "application1",
      "type": "nginx",
      "repo": "git@github.com:customer/application1",
      "root": "/var/www"
    },
    {
      "name": "application2",
      "type": "nodejs",
      "repo": "github.com:customer/application2",
      "docroot": "/var/www/application2"
    }
  ]
}
```

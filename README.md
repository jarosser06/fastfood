fastfood
========

###External Dependencies
This tool uses gpm for dependency management so you need to
ensure that it is installed in order to build from source.

https://github.com/pote/gpm

###Notes

CLI Options
```shell
fastfood new 123456-customer

fastfood gen app application1 type:python repo:git@github.com/cus/app1

fastfood gen db mysql_master database:appdb
```

Set configs, global being the default and then allow for overrides.

Make Cookbooks faster by pre-templating them

Uses a json config file to allow for pre-templating cookbooks:

```json
{
  "id": "123456-customer",
  "cookbook_path": "/home/jim/cookbooks",
  "applications": [
    {
      "name": "application1",
      "type": "python",
      "repo": "github.com:customer/application1",
      "docroot": "/var/www/application1",
      "webserver": "nginx"
    },
    {
      "name": "application2",
      "type": "nodejs",
      "repo": "github.com:customer/application2",
      "docroot": "/var/www/application2"
    }
  ],
  "databases": [
    {
      "type": "mysql",
      "name": "customerdb",
      "user": "db_user"
    }
  ]
}
```

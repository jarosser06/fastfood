Contributing to Fastfood
========================

Fastfood now pulls templates, and data from a template bundle.  This
bundle includes a core manifest that has the available commands as well
as individual providers, each with their own manifest.  In order to add
to an existing provider you would need to modify that providers manifest.

For example, adding a new type to the data base provider involves creating the
template files and adding them to their respective locations and then adding the
new type to the types list:

```json
"types": {
  "cockroachdb": {
    "dependencies": [
      "cockroachdb",
      "lightson"
    ],
    "files": {
      "recipes/%s.rb": "recipes/cockroachdb",
      "libraries/cockroach_helpers.rb": "libraries/cockroach_helpers.rb",
      "test/unit/spec/%s_spec.rb": "test/unit/spec/coackroachdb"
    },
    "directories": [
      "libraries"
    ]
  }
}
```

Currently fastfood uses the default Golang templating language.  It is on the
roadmap to make the templating language a bit friendlier but for now anything
you need to know about the templates can be found [here](http://golang.org/pkg/text/template/).

TODO: Template example

####Note
If you need to ensure that directories other than just, "recipes", and
"test/unit/spec" are created for the cookbook in order to user your
fastfood provider then you will need to add them to the directories
list.  Even if the directory gets created by default you should absolutely
not assume that it still exists.


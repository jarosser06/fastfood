## 0.2.0 (Unreleased)

* Cookbooks built or modified with a json config have the config
  added to the cookbook if it does not already exist.
* Added a build command that allows for building cookbooks from
  a json file.
* Moved the manifest functionality from the fastfood/cmd package
* Added a force flag to the gen command
* Added trailing newline when appending to dependencies
* Added actual help to the new command
* Dependency appending no longer appends extra newlines if there
  are no dependencies.

## 0.1.2 (17 Dec 2014)

* Fastfood will now return an error if it cannot parse a provider
  manifest.

## 0.1.1 (16 Dec 2014)

* Fastfood no longer panics if no additional arguments are passed
  to the gen command.  It now presents a list of providers available.

## 0.1.0 (14 Dec 2014)

* Basic gen command that builds out Chef files (templates,
  recipes, tests, etc) based on a Template Pack
* Basic new command that will generate a new cookbook based on
  a set of template files located in the Template Pack

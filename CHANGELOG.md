## 0.3.0 (1 Feb 2015)

* Added support for extra dependency options
* Added support for tag, branch, and ref for berks dependencies
* Cleaned up option help to reduce confusion with path expansion
* Removed remaining references to old naming (stencil set/provider)
* Error message cleanup

## 0.2.2 (20 Jan 2015)

* Fixed bug that broke the new command

## 0.2.1 (8 Jan 2015)

* Fixed bad message when new cookbook is created

## 0.2.0 (8 Jan 2015)

* Basic support for berks dependency handling
* Reworked internal API to allow for future frameworks
* Added API versions to stencils and templatepack manifests
* Moved all stencils to stencils directory to avoid confusion
* Split fastfood util into seperate packages stringutil and fileutil
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

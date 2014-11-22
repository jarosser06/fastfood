#
# Cookbook Name:: |{.Cookbook.Name}|
# Recipe :: |{.Name}|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#
|{ if eq .Webserver "apache" }|
include_recipe 'apache2'
|{ else }|
node.default['nginx']['default_site_enabled'] = false
include_recipe 'nginx'
|{ end }|

  application |{ .QString .Name }| do
  path |{ .QString .Path }|
  owner |{ .QString .Owner }|
  group |{ .QString .Owner }|
  repository |{ .QString .Repo }|
  revision |{ .QString .Revision }|
end

|{ if eq .Webserver "apache" }|
|{template "ApacheSite" .}|
|{ else }|
|{template "NginxSite" .}|
|{end}|

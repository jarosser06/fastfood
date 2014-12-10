#
# Cookbook Name:: |{.Cookbook.Name}|
# Recipe :: |{.Options.Name}|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#
|{ if eq .Options.Webserver "apache" }|
include_recipe 'apache2'
|{ else }|
node.default['nginx']['default_site_enabled'] = false
include_recipe 'nginx'
|{ end }|

application |{ .QString .Options.Name }| do
  path ::File.join(|{ .QString .Options.Root }|, |{ .QString .Options.Name }|)
  owner |{ .QString .Options.Owner }|
  group |{ .QString .Options.Owner }|
  repository |{ .QString .Options.Repo }|
  revision |{ .QString .Options.Revision }|
end

|{ if eq .Options.Webserver "apache" }|
|{template "ApacheSite" .}|
|{ else }|
|{template "NginxSite" .}|
|{end}|

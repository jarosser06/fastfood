#
# Cookbook Name:: |{.Cookbook.Name}|
# Recipe :: |{.Name}|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#
|{ if eq .Webserver "apache" }|
include_recipe 'apache2'
|{ else }|
node.default['nginx']['default_site_enabled']
include_recipe 'nginx'
|{ end }|

application '|{ .Name }|' do
  path '|{ .Path }|'
  owner '|{ .Owner }|'
  group '|{ .Owner }|'
  repository '|{ .Repo }|'
  revision '|{ .Revision }|'
end

|{ if eq .Webserver "apache" }|
|{template "ApacheSite" .}|
|{ else }|
|{template "NginxSite" .}|
|{end}|

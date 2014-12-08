#
# Cookbook Name:: |{ .Cookbook.Name }|
# Recipe :: |{ .Options.name }|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#

|{ .Options }|
include_recipe 'mysql-multi::mysql_master'

|{ if ne .Options.openfor "" }|
search_string = "chef_environment:#{node.chef_environment} AND tags:|{ .Options.openfor }|"
search_add_iptables_rules(
  search_string,
  'INPUT',
  "-m tcp -p tcp --dport #{node['mysql']['port']} -j ACCEPT",
  70,
  'access to mysql'
)
|{ end }|

|{ if ne .Options.database "" }|
conn = {
  host:  'localhost',
  username: 'root',
  password: node['mysql']['server_root_password']
}

include_recipe 'database::mysql'
|{ if ne .Options.databag "" }|
mysql_creds = Chef::EncryptedDataBagItem.load(
  '|{ .Options.Databag }|',
  node.chef_environment
)

mysql_database |{ .QString .Options.database }| do
  connection conn
  action :create
end

mysql_database_user mysql_creds['username'] do
  connection conn
  password mysql_creds['password']
  database |{ .QString .Options.database }|
  action :create
end
|{ else }|
  mysql_database |{ .QString .Options.database }| do
  connection conn
  action :create
end

mysql_database_user |{ .Options.user }| do
  connection conn
  password |{ .Options.password}|
  database_name |{ .QString .Options.database }|
  action :create
end
|{ end }|
|{ end }|

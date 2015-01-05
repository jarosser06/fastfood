#
# Cookbook Name:: |{ .Cookbook.Name }|
# Recipe :: |{ .Options.Name }|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#

include_recipe 'mysql-multi::mysql_master'

|{ if ne .Options.Openfor "" }|
search_string = "chef_environment:#{node.chef_environment} AND tags:|{ .Options.Openfor }|"
search_add_iptables_rules(
  search_string,
  'INPUT',
  "-m tcp -p tcp --dport #{node['mysql']['port']} -j ACCEPT",
  70,
  'access to mysql'
)
|{ end }|

|{ if ne .Options.Database "" }|
conn = {
  host:  'localhost',
  username: 'root',
  password: node['mysql']['server_root_password']
}

include_recipe 'database::mysql'
|{ if ne .Options.Databag "" }|
mysql_creds = Chef::EncryptedDataBagItem.load(
  '|{ .Options.Databag }|',
  node.chef_environment
)

mysql_database |{ .QString .Options.Database }| do
  connection conn
  action :create
end

mysql_database_user mysql_creds['username'] do
  connection conn
  password mysql_creds['password']
  database |{ .QString .Options.Database }|
  action :create
end
|{ else }|
  mysql_database |{ .QString .Options.Database }| do
  connection conn
  action :create
end

mysql_database_user |{ .Options.User }| do
  connection conn
  password |{ .Options.Password}|
  database_name |{ .QString .Options.Database }|
  action :create
end
|{ end }|
|{ end }|

#
# Cookbook Name:: |{ .Cookbook.Name }|
# Recipe :: |{ .Name }|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#

include_recipe 'mysql-multi::mysql_|{ .Role }|'

|{ if ne .Openfor "" }|
search_string = "chef_environment:#{node.chef_environment} AND tags:|{ .Openfor }|"
search_add_iptables_rules(
  search_string,
  'INPUT',
  "-m tcp -p tcp --dport #{node['mysql']['port']} -j ACCEPT",
  70,
  'access to mysql'
)
|{ end }|

|{ if ne .Database "" }|
conn = {
  host:  'localhost',
  username: 'root',
  password: node['mysql']['server_root_password']
}

|{ if ne .Databag "" }|
mysql_creds = Chef::EncryptedDataBagItem.load(
  '|{ .Databag }|',
  node.chef_environment
)

mysql_database |{ .QString .Database }| do
  connection conn
  action :create
end

mysql_database_user mysql_creds['username'] do
  connection conn
  password mysql_creds['password']
  database |{ .QString .Database }|
  action :create
end
|{ else }|
mysql_database |{ .QString .Database }| do
  connection conn
  action :create
end

mysql_database_user |{ .SetorReturnDatabase .User }| do
  connection conn
  password |{ .SetOrReturnDatabase .Password}|
  database |{ .QString .Database }|
  action :create
end
|{ end }|
|{ end }|

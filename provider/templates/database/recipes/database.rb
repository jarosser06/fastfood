#
# Cookbook Name:: |{ .Cookbook.Name }|
# Recipe :: |{ .Name }|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#

include_recipe 'mysql-multi::|{ .Role }|'

|{ if ne .Database "" }|
conn = {
  host:  'localhost',
  username: 'root',
  password: node['mysql']['server_root_password']
}

mysql_database |{ .QString .Database }| do
  connection conn
  action :create
end

mysql_database_user |{ .SetOrReturnDatabase .User }| do
  connection conn
  password |{ .SetOrReturnDatabase .Password }|
  database |{ .QString .Database }|
  action :create
end
|{ end }|

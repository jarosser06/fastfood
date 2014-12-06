#
# Cookbook Name:: |{ .Cookbook.Name }|
# Recipe :: |{ .Name }|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#
|{ if ne .Databag "" }|
postgres_creds = Chef::EncryptedDataBagItem.load(
  '|{ .Databag }|',
  node.chef_environment
)
|{ end }|
|{ if ne .Openfor "" }|
search_string = "chef_environment:#{node.chef_environment} AND tags:|{ .Openfor }|"
search_add_iptables_rules(
  search_string,
  'INPUT',
  "-m tcp -p tcp --dport #{node['postgresql']['config']['port']} -j ACCEPT",
  70,
  'access to postgres'
)
|{ if ne .Database "" }|
include_recipe 'chef-sugar'

app_nodes = search(:node, search_string)
unless app_nodes.empty?
  app_nodes.each do |anode|
    node.default['postgresql']['pg_hba'] << {
      command: '# authorize app server',
      type: 'host',
      |{ if ne .Databag "" }|
      user: postgres_creds['username'],
      |{ else }|
      user: |{ .SetOrReturnDatabase .User }|,
      |{ end }|
      addr: "#{best_ip_for(anode)}/32",
      method: 'md5'
    }
  end
end
|{ end }|
|{ end }|
include_recipe 'pg-multi'
include_recipe 'database::postgres'
|{ if ne .Database "" }|
conn = {
  host: 'localhost',
  username: 'postgres',
  password: node['postgresql']['password']['postgres']
}
|{ if ne .Databag "" }|
postgresql_database |{ .Database }| do
  connection conn
  action :create
end

postgresql_database_user postgres_creds['username'] do
  connection conn
  action :create
  database_name |{ .QString .Database }|
  password postgres_creds['password']
  privileges [:all]
end
|{ else }|
postgresql_database |{ .Database }| do
  connection conn
  action :create
end

postgresql_database_user |{ .SetOrReturnDatabase .User }| do
  connection conn
  action :create
  database_name |{ .QString .Database }|
  password |{ .SetOrReturnDatabase .Password }|
  privileges [:all]
end
|{ end }|
|{ end }|

#
# Cookbook Name:: |{ .Cookbook.Name }|
# Recipe :: |{ .Options.Name }|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#
|{ if ne .Options.Databag "" }|
postgres_creds = Chef::EncryptedDataBagItem.load(
  '|{ .Options.Databag }|',
  node.chef_environment
)
|{ end }|
|{ if ne .Options.Openfor "" }|
search_string = "chef_environment:#{node.chef_environment} AND tags:|{ .Options.Openfor }|"
search_add_iptables_rules(
  search_string,
  'INPUT',
  "-m tcp -p tcp --dport #{node['postgresql']['config']['port']} -j ACCEPT",
  70,
  'access to postgres'
)
|{ if ne .Options.Database "" }|
include_recipe 'chef-sugar'

app_nodes = search(:node, search_string)
unless app_nodes.empty?
  app_nodes.each do |anode|
    node.default['postgresql']['pg_hba'] << {
      command: '# authorize app server',
      type: 'host',
      |{ if ne .Options.Databag "" }|
      user: postgres_creds['username'],
      |{ else }|
      user: |{ .Options.User }|,
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
|{ if ne .Options.Database "" }|
conn = {
  host: 'localhost',
  username: 'postgres',
  password: node['postgresql']['password']['postgres']
}
|{ if ne .Options.Databag "" }|
postgresql_database |{ .Options.Database }| do
  connection conn
  action :create
end

postgresql_database_user postgres_creds['username'] do
  connection conn
  action :create
  database_name |{ .QString .Options.Database }|
  password postgres_creds['password']
  privileges [:all]
end
|{ else }|
  postgresql_database |{ .Options.Database }| do
  connection conn
  action :create
end

postgresql_database_user |{ .Options.User }| do
  connection conn
  action :create
  database_name |{ .QString .Options.Database }|
  password |{ .Options.Password }|
  privileges [:all]
end
|{ end }|
|{ end }|

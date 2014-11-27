#
# Cookbook Name:: |{ .Cookbook.Name }|
# Recipe :: |{ .Name }|
#
# Copyright |{ .Cookbook.Year }|, Rackspace
#
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
      user: |{ .SetorReturnDatabase .User }|,
      addr: "#{best_ip_for(anode)}/32",
      method: 'md5'
    }
  end
end
|{ end }|
|{ end }|
include_recipe 'pg-multi'

|{ if ne .Database "" }|
conn = {
  host: 'localhost',
  username: 'postgres',
  password: node['postgresql']['password']['postgres']
}

postgresql_database node['employii']['postgres_database'] do
  connection conn
  action :create
end

postgresql_database_user |{ .SetorReturnDatabase .User }| do
  connection conn
  action :create
  password |{ .SetorReturnDatabase .Password }|
  privileges [:all]
end
|{ end }|

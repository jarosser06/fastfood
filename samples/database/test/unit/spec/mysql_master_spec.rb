require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Options.name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Options.name }|')
  end

  it 'includes the mysql-multi::mysql_master recipe' do
    expect(chef_run).include_recipe('mysql-multi::mysql_master')
  end
  |{ if ne .Options.database "" }|
  it 'creates the mysql_database |{ .Options.database }|' do
    expect(chef_run).create_mysql_database(|{ .Options.database }|)
  end
|{ end }|
end

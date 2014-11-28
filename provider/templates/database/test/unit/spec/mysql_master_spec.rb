require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Name }|')
  end

  it 'includes the mysql-multi::mysql_master recipe' do
    expect(chef_run).include_recipe('mysql-multi::mysql_master')
  end
|{ if ne .Database "" }|
  it 'creates the mysql_database |{ .Database }|' do
    expect(chef_run).create_mysql_database(|{ .Database }|)
  end
|{ end }|
end

require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Options.Name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Options.Name }|')
  end

  it 'includes the mysql-multi::mysql_master recipe' do
    expect(chef_run).include_recipe('mysql-multi::mysql_master')
  end
  |{ if ne .Options.Database "" }|
  it 'creates the mysql_database |{ .Options.Database }|' do
    expect(chef_run).create_mysql_database(|{ .Options.Database }|)
  end
|{ end }|
end

require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Name }|')
  end

  it 'includes the mysql-multi::mysql_slave recipe' do
    expect(chef_run).include_recipe('mysql-multi::mysql_slave')
  end
end

require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Name }|')
  end

  it 'includes the mysql-multi::|{ .Role }| recipe' do
    expect(chef_run).include_recipe('mysql-multi::|{ .Role }|')
  end
end

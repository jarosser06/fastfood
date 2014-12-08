require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Options.name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Options.name }|')
  end

  it 'includes the pg-multi::pg_master recipe' do
    expect(chef_run).include_recipe('pg-multi::pg_master')
  end

|{ if ne .Options.database "" }|
  it 'creates the database |{ .Options.database }||' do
    expect(chef_run).create_postgresql_database(|{ .Options.database }|)
  end
|{ end }|
end

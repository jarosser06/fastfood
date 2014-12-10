require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Options.Name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Options.Name }|')
  end

  it 'includes the pg-multi::pg_master recipe' do
    expect(chef_run).include_recipe('pg-multi::pg_master')
  end

|{ if ne .Options.Database "" }|
  it 'creates the database |{ .Options.Database }||' do
    expect(chef_run).create_postgresql_database(|{ .Options.Database }|)
  end
|{ end }|
end

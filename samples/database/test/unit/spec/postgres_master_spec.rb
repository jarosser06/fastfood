require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Name }|')
  end

  it 'includes the pg-multi::pg_master recipe' do
    expect(chef_run).include_recipe('pg-multi::pg_master')
  end

|{ if ne .Database "" }|
  it 'creates the database |{ .Database }||' do
    expect(chef_run).create_postgresql_database(|{ .Database }|)
  end
|{ end }|
end

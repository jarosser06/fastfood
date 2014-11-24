require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Name }|')
  end

|{ if eq .Webserver "apache" }|
  it 'includes the apache recipe' do
    expect(chef_run).include_recipe('apache')
  end
|{ else }|
  it 'includes the nginx recipe' do
    expect(chef_run).include_recipe('nginx')
  end

  |{ template "NginxSiteTest" . }|
|{ end }|
end

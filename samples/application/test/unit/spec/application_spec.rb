require_relative 'spec_helper'

describe '|{ .Cookbook.Name }|::|{ .Options.Name }|' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('|{ .Cookbook.Name }|::|{ .Options.Name }|')
  end

|{ if eq .Options.Webserver "apache" }|
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

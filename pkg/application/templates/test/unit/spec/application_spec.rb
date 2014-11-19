require_relative 'spec_helper'

descript '{{ .Cookbook.Name }}::{{ .Name }}' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('{{ .Cookbook.Name }}::{{ .Name }}')
  end
end

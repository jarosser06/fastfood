require_relative 'spec_helper'

describe 'testcookbook::default' do
  let(:chef_run) do
    ChefSpec::Runner.new.converge('testcookbook::default')
  end
end

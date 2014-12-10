|{ define "ApacheSite" }|
web_app |{ .QString .Options.Name }| do
  cookbook |{ .QString .Cookbook.Name }|
  docroot |{ .QString .Options.Path }|
  enable true
end
|{ end }|

|{ define "NginxSite" }|
site_conf = ::File.join(node['nginx']['dir'],
                        'sites-available',
                        |{ .QString .Options.Name }|)

template site_conf do
  source "|{ .Options.Name }|.nginx.erb"
  action :create
  notifies :reload, 'service[nginx]'
end

nginx_site |{ .QString .Options.Name }| do
  action enable
end
|{ end }|

|{ define "NginxSiteTest" }|
  it 'creates the site template' do
  expect(chef_run).to create_template("/etc/nginx/sites-available/|{ .Options.Name }|")
  end
|{ end }|

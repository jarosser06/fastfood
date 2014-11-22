|{ define "ApacheSite" }|
web_app |{ .Name }| do
  cookbook |{ .Cookbook.Name }|
  docroot |{ .Path }|
  enable true
end
|{ end }|

|{ define "NginxSite" }|
site_conf = ::File.join(node['nginx']['dir'],
                        'sites-available',
                        |{ .Name }|)

template site_conf do
  source "|{ .Name }|.nginx.erb"
  action :create
  notifies :reload, 'service[nginx]'
end

nginx_site |{ .Name }| do
  action enable
end
|{ end }|

|{ define "NginxSiteTest" }|
  it 'creates the site template' do
    expect(chef_run).to create_template("/etc/nginx/sites-available/|{ .Name }|")
  end
|{ end }|

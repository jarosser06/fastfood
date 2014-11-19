application '{{ .Name }}' do
  path '{{ .Path }}'
  owner 'www-data'
  group 'www-data'
  repository '{{ .Repo }}'
  revision 'master'
end

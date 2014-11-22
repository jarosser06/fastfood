application '{{ .Name }}' do
  path '{{ .Path }}'
  owner '{{ .Owner }}'
  group '{{ .Owner }}'
  repository '{{ .Repo }}'
  revision '{{ .Branch }}'
end

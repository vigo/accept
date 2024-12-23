task :default => [:test]

desc 'run test'
task :test do
  system %{ go test -v -coverprofile=coverage.out ./... }
  exit $?.exitstatus
end

desc 'show test coverage'
task :coverage do
  system %{ 
    go test -v -coverprofile=coverage.out ./... &&
    go tool cover -html=coverage.out
  }
  exit $?.exitstatus
end

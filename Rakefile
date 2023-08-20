namespace :run do
  desc "run server"
  task :server do
    system %{
      SERVER_ENV=#{ENV['SERVER_ENV'] || "local"} go run -race cmd/server/main.go
    }
    exit $?.exitstatus
  end
end

desc "default task"
task :default => ["run:server"]

desc "run golangci lint"
task :lint do
  system %{
    LOG_LEVEL=error golangci-lint run
  }
  exit $?.exitstatus
end